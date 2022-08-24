package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func fatalError(err string) {
	os.Stderr.WriteString(err)
	os.Exit(1)

}

type flagsCmd struct {
	A, B, C       *uint
	c, i, v, n, F *bool
}

type fileSearch struct {
	rows      []string
	flags     flagsCmd
	searchArg string
}

func (f fileSearch) Compare(str string) bool {
	strRes := str
	if *f.flags.i {
		strRes = strings.ToLower(strRes)
	}
	if *f.flags.F {
		return strRes == f.searchArg
	}
	return strings.Contains(strRes, f.searchArg)
}

func (f fileSearch) Search() {
	rowsMatchIndex := make([]int, 0)
	for i := 0; i < len(f.rows); i++ {
		if f.Compare(f.rows[i]) {
			rowsMatchIndex = append(rowsMatchIndex, i)
		}
	}
	if len(rowsMatchIndex) == 0 {
		return
	}
	f.print(rowsMatchIndex)
}

func (f fileSearch) print(rowsIndex []int) {
	if *f.flags.c {
		if *f.flags.v {
			os.Stdout.WriteString(fmt.Sprint(len(f.rows)-len(rowsIndex), "\n"))
		} else {
			os.Stdout.WriteString(fmt.Sprint(len(rowsIndex), "\n"))
		}
		return
	}
	var after, before uint = 0, 0
	if *f.flags.C != 0 {
		after, before = *f.flags.C, *f.flags.C
	} else {
		after, before = *f.flags.A, *f.flags.B
	}
	if *f.flags.v {
		start := 0
		for _, el := range rowsIndex {
			f.rowsWriter(start, el-int(after)-1)
			start = el + int(before) + 1
		}
		f.rowsWriter(start, len(f.rows))
	} else {
		for _, el := range rowsIndex {
			f.rowsWriter(el-int(after), el+int(before))
		}
	}
}

func (f fileSearch) rowsWriter(start, end int) {
	if start < 0 {
		start = 0
	}
	if end >= len(f.rows) {
		end = len(f.rows) - 1
	}
	for i := start; i <= end; i++ {
		printStr := f.rows[i] + "\n"
		if *f.flags.n {
			printStr = fmt.Sprint(i, ": ") + printStr
		}
		os.Stdout.WriteString(printStr)
	}
}

func (f *fileSearch) Open(args []string) {
	if len(args) < 2 {
		fatalError("недостаточно аргументов")
	}
	f.searchArg = args[0]
	file, err := os.Open(args[1])
	if err != nil {
		fatalError(err.Error())
	}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		f.rows = append(f.rows, sc.Text())
	}
	file.Close()
}

func initFlags() ([]string, flagsCmd) {
	flags := flagsCmd{}
	flags.c = flag.Bool("c", false, "Количество строк")
	flags.i = flag.Bool("i", false, "Игнорировать регистр")
	flags.v = flag.Bool("v", false, "Вместо совпадения, исключать")
	flags.n = flag.Bool("n", false, "Выводить номер строки")
	flags.A = flag.Uint("A", 0, "Печатать +N строк после совпадения")
	flags.B = flag.Uint("B", 0, "Печатать +N строк до совпадения")
	flags.C = flag.Uint("C", 0, "Печатать +-N строк вокруг совпадения")
	flags.F = flag.Bool("F", false, "Точное совпадение со строкой")
	flag.Parse()
	return flag.Args(), flags
}

func main() {
	args, flags := initFlags()
	file := fileSearch{flags: flags}
	file.Open(args)
	file.Search()
}
