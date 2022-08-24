package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

type flagsCmd struct {
	l *uint
}

type filehttp struct {
	depth, maxDepth                 int
	url, filename, mainpath, header string
	flagsCmd
	linksExist map[string]string
}

func (f *filehttp) getBody(url string) (*io.ReadCloser, string, error) {
	resp, errRespo := http.Get(url)
	if errRespo != nil {
		return nil, "", errRespo
	}
	if resp.StatusCode != 200 {
		return nil, "", errors.New("received non 200 response code")
	}
	if f.header == "" {
		f.header = resp.Header.Get("Content-Type")
	}
	return &resp.Body, path.Base(resp.Request.URL.Path), nil
}

func (f filehttp) getBodyStr(body io.ReadCloser) (string, error) {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f filehttp) utf8encode(str string) string {
	utf8, err := charset.NewReader(strings.NewReader(str), f.header)
	if err != nil {
		return getRandStr()
	}
	b, err := ioutil.ReadAll(utf8)
	if err != nil {
		return getRandStr()
	}
	return string(b)
}

func (f filehttp) CheckHtml(body string) bool {
	return strings.Contains(strings.ToLower(body), "!doctype html")
}

func (f filehttp) getLinks(body string) []string {
	reg, err := regexp.Compile(`href=".+?"|src=".+?"|href='.+?'|src='.+?'|href=.+?[> ]|src=.+?[> ]`)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return nil
	}
	return reg.FindAllString(body, -1)

}

func getRandStr() string {
	rand.Seed(time.Now().UnixNano())
	var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	bytes := make([]byte, 10)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}

func getRightLink(link, url string) (string, string, error) {
	linkWithoutHrefSrc := ""
	linkWithoutHrefSrc = strings.TrimLeft(link, "href=")
	linkWithoutHrefSrc = strings.TrimLeft(linkWithoutHrefSrc, "src=")
	linkWithoutHrefSrc = strings.Trim(linkWithoutHrefSrc, `\"`)
	linkWithoutHrefSrc = strings.Trim(linkWithoutHrefSrc, `\'`)
	linkWithoutHrefSrc = strings.Trim(linkWithoutHrefSrc, ">")
	linkWithoutHrefSrc = strings.TrimSpace(linkWithoutHrefSrc)
	if linkWithoutHrefSrc == "/" || len(linkWithoutHrefSrc) == 0 {
		return "", "", errors.New("empty link")
	}
	rightLink := linkWithoutHrefSrc
	reg, _ := regexp.Compile(`.*://.+?/`)
	domen := reg.FindString(url)
	protocol := ""
	if strings.Contains(domen, "https") {
		protocol = "https://"
	} else if strings.Contains(domen, "http") {
		protocol = "http://"
	} else {
		return "", "", errors.New("invalid protocol")
	}
	domen = strings.Trim(strings.Replace(domen, protocol, "", 1), "/")
	splitDomen := strings.Split(domen, ".")
	lenSplitDomen := len(splitDomen)
	domen = splitDomen[lenSplitDomen-2] + "." + splitDomen[lenSplitDomen-1]
	subdomens := ""
	for i := 0; i < len(splitDomen)-2; i++ {
		subdomens += splitDomen[i] + "."
	}
	if linkWithoutHrefSrc[0] == '/' {
		if linkWithoutHrefSrc[1] == '/' {
			rightLink = strings.TrimLeft(linkWithoutHrefSrc, "/")
			rightLink = protocol + subdomens + domen + "/" + rightLink
		} else {
			rightLink = protocol + subdomens + domen + linkWithoutHrefSrc
		}
	}
	if !strings.Contains(rightLink, domen) {
		return "", "", errors.New("invalid domen")
	}
	return rightLink, linkWithoutHrefSrc, nil
}

func (f filehttp) getFilename(body, url, respFilename string) string {
	if f.filename != "" {
		return f.filename
	}
	if respFilename != "" && strings.Contains(respFilename, ".") {
		return respFilename
	}
	name := ""
	if !f.CheckHtml(body) {
		splitUrl := strings.Split(url, "/")
		name = splitUrl[len(splitUrl)-1]
		if name == "" {
			name = getRandStr()
		}
	} else {
		regName, _ := regexp.Compile(`<title.*?>.+?</title>`)
		name = regName.FindString(body)
		if name == "" {
			name = getRandStr()
		}
	}
	regNameCorrect, _ := regexp.Compile(`[\?\*/\\\|<>:\"\#\&]|title| itemprop=name`)
	name = regNameCorrect.ReplaceAllString(name, "")
	if _, errFileExist := os.Stat(name); !os.IsNotExist(errFileExist) {
		name = getRandStr() + name
	}
	name = strings.ReplaceAll(name, " ", "")
	return name
}

func (f *filehttp) validHtmlAndGetStrBody(url string) (string, string, error) {
	if f.depth >= f.maxDepth {
		return "", "", errors.New("maximum depth exceeded")
	}
	bodyReader, filename, errGet := f.getBody(url)
	if errGet != nil {
		return "", "", errGet
	}
	body, errGetStr := f.getBodyStr(*bodyReader)
	if errGetStr != nil {
		return "", "", errGetStr
	}
	if !f.CheckHtml(body) {
		return body, filename, errors.New("not valid html file")
	}
	return body, "", nil
}

func (f *filehttp) createDir(name string) error {
	if f.depth == 1 {
		utf8name := f.utf8encode(name)
		err := os.Mkdir(utf8name, os.ModeDir)
		if err != nil {
			return err
		}
		errChdir := os.Chdir(utf8name)
		if errChdir != nil {
			return errChdir
		}
		f.mainpath = name
	}
	return nil
}

func (f filehttp) writeMainHtmlFile(filename, body string) error {
	if f.depth == 1 {
		os.Chdir("..")
	}
	errWrite := f.writeFile(filename+".html", io.NopCloser(strings.NewReader(body)))
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func (f filehttp) replaceL(body, link, linkWithoutHrefSrc, filename, filenameLink string) string {
	replaceLink := ""
	if f.depth == 1 {
		replaceLink = strings.Replace(link, linkWithoutHrefSrc, filename+"/"+filenameLink, 1)
	} else {
		replaceLink = strings.Replace(link, linkWithoutHrefSrc, filenameLink, 1)
	}

	body = strings.Replace(body, link, replaceLink, 1)

	return body
}

func (f filehttp) replaceEmptylinks(body string) string {
	if f.depth == 1 {
		body = strings.ReplaceAll(body, `href="/"`, `href="`+f.mainpath+`.html"`)
		body = strings.ReplaceAll(body, `href='/'`, `href="`+f.mainpath+`.html'`)
	} else {
		body = strings.ReplaceAll(body, `href="/"`, `href="../`+f.mainpath+`.html"`)
		body = strings.ReplaceAll(body, `href='/'`, `href='../`+f.mainpath+`.html'`)
	}
	return body
}

func (f *filehttp) openHtml(url string) (string, error) {
	body, _, errBody := f.validHtmlAndGetStrBody(url)
	if errBody != nil {
		return "", errBody
	}
	links := f.getLinks(body)
	filename := f.getFilename(body, url, "")
	errMkdir := f.createDir(filename)
	if errMkdir != nil {
		return "", errMkdir
	}
	body = f.replaceEmptylinks(body)
	for i := 0; i < len(links); i++ {
		if f.depth == 1 {
			os.Stdout.WriteString(fmt.Sprint(int(float32(i)/float32(len(links)-1)*100), "%\n"))
		}
		link, linkWithoutHrefSrc, err := getRightLink(links[i], url)
		if err != nil {
			continue
		}
		fNew := filehttp{depth: f.depth + 1, maxDepth: f.maxDepth, linksExist: f.linksExist, url: link, mainpath: f.mainpath}
		filepath, errOpen := fNew.openHtml(link)
		if errOpen == nil {
			body = f.replaceL(body, links[i], linkWithoutHrefSrc, filename, filepath)
			continue
		}
		if l, ok := f.linksExist[links[i]]; ok {
			body = f.replaceL(body, links[i], linkWithoutHrefSrc, filename, l)
			continue
		}
		bodyLink, respFilename, errBodyLink := f.validHtmlAndGetStrBody(link)
		filenameLink := ""
		if errBodyLink == nil {
			filenameLink = ".html"
		}
		filenameLink = f.getFilename(bodyLink, link, respFilename) + filenameLink
		body = f.replaceL(body, links[i], linkWithoutHrefSrc, filename, filenameLink)
		errWrite := f.writeFile(filenameLink, io.NopCloser(strings.NewReader(bodyLink)))
		if errWrite != nil {
			continue
		}
		f.linksExist[links[i]] = filenameLink
	}
	errMain := f.writeMainHtmlFile(filename, body)
	if errMain != nil {
		return "", errMain
	}
	if f.depth == 1 {
		f.editFiles()
	}
	return filename + ".html", nil
}

func (f *filehttp) open() {
	if *f.l == 1 {
		f.download()
	} else {
		f.depth = 1
		f.maxDepth = int(*f.l)
		f.linksExist = make(map[string]string)
		_, err := f.openHtml(f.url)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
		}
	}
}

func (f *filehttp) download() {
	body, _, errGet := f.getBody(f.url)
	if errGet != nil {
		os.Stderr.WriteString(errGet.Error() + "\n")
		os.Exit(1)
	}
	errWrite := f.writeFile(f.filename, *body)
	if errWrite != nil {
		os.Stderr.WriteString(errWrite.Error() + "\n")
		os.Exit(1)
	}
}

func (f filehttp) editFileBody(body string) string {
	f.depth = 2
	body = f.replaceEmptylinks(body)
	links := f.getLinks(body)
	for i := 0; i < len(links); i++ {
		if file, ok := f.linksExist[links[i]]; ok {
			linkWithoutHrefSrc := strings.Trim(links[i], "href=")
			linkWithoutHrefSrc = strings.Trim(linkWithoutHrefSrc, "src=")
			linkWithoutHrefSrc = strings.Trim(linkWithoutHrefSrc, `\"`)
			linkWithoutHrefSrc = strings.Trim(linkWithoutHrefSrc, `\'`)
			linkWithoutHrefSrc = strings.Trim(linkWithoutHrefSrc, ">")
			linkWithoutHrefSrc = strings.TrimSpace(linkWithoutHrefSrc)
			body = f.replaceL(body, links[i], linkWithoutHrefSrc, "", file)
		}
	}
	return body
}

func (f filehttp) editFiles() {
	os.Chdir(f.mainpath)
	files, err := os.ReadDir(".")
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}
	for _, el := range files {
		if strings.Contains(el.Name(), ".html") {
			file, err := os.Open(el.Name())
			if err != nil {
				continue
			}
			bytes, errRead := ioutil.ReadAll(file)
			if errRead != nil {
				file.Close()
				continue
			}
			file.Close()
			os.Remove(el.Name())
			body := f.editFileBody(string(bytes))
			fileNew, errCreate := os.Create(el.Name())
			if errCreate != nil {
				continue
			}
			fileNew.WriteString(body)
			fileNew.Close()
		}
	}
}

func (f filehttp) writeFile(filename string, data io.ReadCloser) error {
	defer data.Close()
	filename = f.utf8encode(filename)
	file, errCreate := os.Create(filename)
	if errCreate != nil {
		return errCreate
	}
	defer file.Close()
	_, errCopy := io.Copy(file, data)
	if errCopy != nil {
		return errCopy
	}
	return nil
}

func initFlags() ([]string, flagsCmd) {
	flags := flagsCmd{}
	flags.l = flag.Uint("l", 1, "maximum depth level")
	flag.Parse()
	return flag.Args(), flags
}

func main() {
	args, flags := initFlags()
	var f filehttp
	switch len(args) {
	case 2:
		f = filehttp{url: args[1], filename: args[0], flagsCmd: flags}
	case 1:
		if *flags.l == 1 {
			os.Stderr.WriteString("Missing filename\n")
			os.Exit(1)
		}
		f = filehttp{url: args[0], flagsCmd: flags}
	default:
		os.Stderr.WriteString("Error count arguments")
		os.Exit(1)
	}
	f.open()
}
