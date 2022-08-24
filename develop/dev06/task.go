package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flagsCmd struct {
	f, d *string
	s    *bool
}

type fileRows struct {
	flagsCmd
	rows    []string
	cutRows []string
}

func fatalError(err string) {
	os.Stderr.WriteString(err)
	os.Exit(1)
}

func (f *fileRows) openFile(args []string) {
	for _, file := range args {
		fil, err := os.Open(file)
		if err != nil {
			fatalError(err.Error())
		}
		sc := bufio.NewScanner(fil)
		for sc.Scan() {
			f.rows = append(f.rows, sc.Text())
		}
		fil.Close()
	}
}

func (f *fileRows) cut() {
	columns := strings.Split(*f.f, ",")
	var columnsInt []int
	for _, el := range columns {
		number, err := strconv.Atoi(el)
		if err != nil {
			fatalError("Некорректно заданы колонки")
		}
		columnsInt = append(columnsInt, number-1)
	}
	strBuilder := strings.Builder{}
	lenColumns := len(columnsInt) - 1
	for _, el := range f.rows {
		splitStr := strings.Split(el, *f.d)
		if len(splitStr) < 2 && *f.s {
			continue
		}
		for i := 0; i < lenColumns; i++ {
			if columnsInt[i] >= len(splitStr) || columnsInt[i] < 0 {
				continue
			}
			strBuilder.WriteString(splitStr[columnsInt[i]])
			strBuilder.WriteString(*f.d)
		}
		if !(columnsInt[lenColumns] >= len(splitStr) || columnsInt[lenColumns] < 0) {
			strBuilder.WriteString(splitStr[columnsInt[lenColumns]])
		}
		f.cutRows = append(f.cutRows, strBuilder.String())
		strBuilder.Reset()
	}
	f.print()
}

func (f fileRows) print() {
	for _, el := range f.cutRows {
		fmt.Println(el)
	}
}

func initFlags() ([]string, flagsCmd) {
	flags := flagsCmd{}
	flags.d = flag.String("d", "\t", "Использовать другой разделитель")
	flags.s = flag.Bool("s", false, "Только строки с разделителем")
	flags.f = flag.String("f", "", "Выбрать поля(колонки)")
	flag.Parse()
	return flag.Args(), flags
}

func main() {
	args, flags := initFlags()
	rows := fileRows{flagsCmd: flags}
	rows.openFile(args)
	rows.cut()
}
