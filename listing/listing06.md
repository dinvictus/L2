Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
Вывод: [3, 2, 3]. Слайс хранит в себе указатели на значения, поэтому при передаче его в функцию указатели на значения копируются,
но все метаданные среза просто копируются, поэтому i[0] = "3" изменит первый элемент переданного слайса, а при использовании функции
append выделится новое место в памяти и это будет уже другой срез

```
