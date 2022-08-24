package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// multiplyString - функция для нахождения количества повторения символа и его повторения
func multiplyString(strRune []rune, start int) (string, int, error) {
	strBuilder := strings.Builder{}
	for i := start + 1; i < len(strRune); i++ {
		if unicode.IsDigit(strRune[i]) {
			strBuilder.WriteRune(strRune[i])
		} else {
			break
		}
	}
	num, err := strconv.Atoi(strBuilder.String())
	if err != nil {
		return "", 0, err
	}
	step := utf8.RuneCountInString(strBuilder.String())
	strBuilder.Reset()
	for i := 0; i < num; i++ {
		strBuilder.WriteRune(strRune[start])
	}
	return strBuilder.String(), step, nil
}

// UnpackingString - функция для распаковки строки ("a4bc2d5e" => "aaaabccddddde", qwe\45 => qwe44444)
func UnpackingString(str string) (string, error) {
	strBuilder := strings.Builder{}
	strRune := []rune(str)
	lenStrRune := len(strRune)
	start := 0
	for i := 0; i < lenStrRune; i++ {
		if unicode.IsDigit(strRune[i]) {
			start++
		} else {
			break
		}
	}
	if start == lenStrRune && start != 0 {
		return "", errors.New("некорректная строка")
	}
	for i := start; i < lenStrRune; i++ {
		if string(strRune[i]) == `\` {
			i++
			if !(i < lenStrRune) {
				break
			}
		}
		if i+1 < lenStrRune && unicode.IsDigit(strRune[i+1]) {
			mulStr, step, err := multiplyString(strRune, i)
			if err != nil {
				return "", err
			}
			strBuilder.WriteString(mulStr)
			i += step
			continue
		}
		strBuilder.WriteRune(strRune[i])
	}
	return strBuilder.String(), nil
}

func main() {
	str := `abc40de5b6`
	unpackStr, err := UnpackingString(str)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	os.Stdout.WriteString(unpackStr + "\n")
}
