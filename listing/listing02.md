Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
Вывод: 2 1. Вызов defer функции происходит после присваивания возратного значения x. 
При вызове функции test defer функция возьмёт x из области видимости, в которой находится возвратное значение x и перезапишет его.
При вызове функции anotherTest defer функция изменит только локально объявленный x, но не изменит возвратное значение.

```
