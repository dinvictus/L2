Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Вывод: error. Случай аналогичен с вопросом 3. error это интерфейс, а из устройства интерфейсов известно, что они хранят два указателя:
один на тип объекта и другой на область в памяти с данными. Интерфейс будет равняться nil, когда эти оба указателя будут nil.
В данном же случае указатель на тип будет *customError, поэтому условие err != nil будет всегда true.

```
