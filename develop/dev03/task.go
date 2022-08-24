package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flagsCmd struct {
	h                *string
	n, r, u, M, b, c *bool
	k                *uint
}

type filesRowsSort struct {
	RowsArr []string
	flags   flagsCmd
}

func (obj filesRowsSort) Len() int {
	return len(obj.RowsArr)
}

func fatalError(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func boolCompareStr(stri, strj string) bool {
	com := strings.Compare(stri, strj)
	return com == -1 || com == 0
}

func compareMonth(stri, strj string) bool {
	months := [12]string{"январь", "февраль", "март", "апрель", "май", "июнь", "июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь"}
	indexi := -1
	indexj := -1
	strilow := strings.ToLower(stri)
	strjlow := strings.ToLower(strj)
	for i := 0; i < 12; i++ {
		if strilow == months[i] {
			indexi = i
		}
		if strjlow == months[i] {
			indexj = i
		}
	}
	if indexi == -1 {
		return true
	}
	if indexj == -1 {
		return false
	}
	return indexi < indexj
}

func compareFlagsLess(stri, strj string, flags flagsCmd) bool {
	strResi := stri
	strResj := strj
	if *flags.b {
		strResi = strings.TrimSpace(stri)
		strResj = strings.TrimSpace(strj)
	}
	if *flags.k != 0 {
		column := *flags.k
		splitStri := strings.Split(strResi, " ")
		splitStrj := strings.Split(strResj, " ")
		splitStriRes := make([]string, 0, len(splitStri))
		splitStrjRes := make([]string, 0, len(splitStrj))
		for _, el := range splitStri {
			if !(el == "" || el == " ") {
				splitStriRes = append(splitStriRes, el)
			}
		}
		for _, el := range splitStrj {
			if !(el == "" || el == " ") {
				splitStrjRes = append(splitStrjRes, el)
			}
		}
		if int(column) > len(splitStriRes) {
			return true
		}
		if int(column) > len(splitStrjRes) {
			return false
		}
		strResi = splitStriRes[column-1]
		strResj = splitStrjRes[column-1]
	}
	if *flags.M {
		return compareMonth(strResi, strResj)
	}
	if *flags.h != "" {
		compareLetters := *flags.h
		runei := []rune(strResi)
		runej := []rune(strResj)
		lasti := runei[len(runei)-1]
		lastj := runej[len(runej)-1]
		indexi := strings.Index(compareLetters, string(lasti))
		indexj := strings.Index(compareLetters, string(lastj))
		if indexi == -1 {
			return true
		}
		if indexj == -1 {
			return false
		}
		if indexi < indexj {
			return true
		} else if indexi > indexj {
			return false
		} else {
			numi, erri := strconv.Atoi(string(runei[:len(runei)-1]))
			numj, errj := strconv.Atoi(string(runej[:len(runej)-1]))
			if erri != nil {
				return true
			}
			if errj != nil {
				return false
			}
			return numi < numj
		}
	}
	if *flags.n {
		numi, erri := strconv.ParseFloat(strResi, 64)
		numj, errj := strconv.ParseFloat(strResj, 64)
		if erri != nil {
			return true
		}
		if errj != nil {
			return false
		}
		return numi < numj
	}
	return boolCompareStr(strResi, strResj)

}

func (obj filesRowsSort) Less(i, j int) bool {
	return compareFlagsLess(obj.RowsArr[i], obj.RowsArr[j], obj.flags)
}

func (obj filesRowsSort) Swap(i, j int) {
	obj.RowsArr[i], obj.RowsArr[j] = obj.RowsArr[j], obj.RowsArr[i]
}

func (obj filesRowsSort) Sort() {
	if *obj.flags.u {
		mapU := make(map[string]struct{})
		var resSlice []string = make([]string, 0, len(obj.RowsArr))
		for _, el := range obj.RowsArr {
			if _, ok := mapU[el]; !ok {
				mapU[el] = struct{}{}
				resSlice = append(resSlice, el)
			}
		}
		obj.RowsArr = resSlice
	}
	if *obj.flags.r {
		sort.Sort(sort.Reverse(obj))
	} else {
		sort.Sort(obj)
	}
	obj.print()
	if *obj.flags.c {
		obj.check()
	}
}

func (obj filesRowsSort) print() {
	for _, el := range obj.RowsArr {
		fmt.Println(el)
	}
}

func (obj filesRowsSort) check() {
	if *obj.flags.r {
		for i := len(obj.RowsArr) - 1; i > 0; i-- {
			if !compareFlagsLess(obj.RowsArr[i], obj.RowsArr[i-1], obj.flags) {
				fatalError(errors.New("данные отсортированы неверно"))
			}
		}
	} else {
		for i := 0; i < len(obj.RowsArr)-1; i++ {
			if !compareFlagsLess(obj.RowsArr[i], obj.RowsArr[i+1], obj.flags) {
				fatalError(errors.New("данные отсортированы неверно"))
			}
		}
	}
	fmt.Println("Данные отсортированы верно")
}

func openFiles(files []string) filesRowsSort {
	var filesRows filesRowsSort
	byteBuffer := bytes.Buffer{}
	for i := 0; i < len(files); i++ {
		file, err := os.Open(files[i])
		if err != nil {
			fatalError(err)
		}
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			byteBuffer.WriteString(sc.Text())
			byteBuffer.WriteString("\n")
		}
		file.Close()
	}
	splitStr := strings.Split(byteBuffer.String(), "\n")
	filesRows.RowsArr = make([]string, len(splitStr)-1)
	copy(filesRows.RowsArr, splitStr[:len(splitStr)-1])

	return filesRows
}

func initFlags() ([]string, flagsCmd) {
	flags := flagsCmd{}
	flags.k = flag.Uint("k", 0, "Указание колонки для сортировки")
	flags.n = flag.Bool("n", false, "Сортировать по числовому значению")
	flags.r = flag.Bool("r", false, "Сортировать в обратном порядке")
	flags.u = flag.Bool("u", false, "Не выводить повторяющиеся строки")
	flags.M = flag.Bool("M", false, "Сортировать по названию месяца")
	flags.b = flag.Bool("b", false, "Игнорировать хвостовые пробелы")
	flags.c = flag.Bool("c", false, "Проверить отсортированы ли данные")
	flags.h = flag.String("h", "", "Сортировать по числовому значению с учётом суффиксов - -h `bkmg`")
	flag.Parse()
	return flag.Args(), flags
}

func main() {
	args, flags := initFlags()
	filesRows := openFiles(args)
	filesRows.flags = flags
	filesRows.Sort()
}
