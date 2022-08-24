package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

	Стратегия — это поведенческий паттерн, выносит набор алгоритмов в собственные структуры и делает их взаимозаменимыми.
	Другие объекты содержат ссылку на объект-стратегию и делегируют ей работу.
	Программа может подменить этот объект другим, если требуется иной способ решения задачи.
*/

// DeleteInterface интерфейс для удаления
type DeleteInterface interface {
	delete(*Data, int)
}

// Data структура для хранения какой-то информации
type Data struct {
	sliceData []int
	name      string
	DeleteInterface
}

// Add метод добавления новой информации
func (sD *Data) Add(element int) {
	sD.sliceData = append(sD.sliceData, element)
}

// Init метод инициализации структуры
func (sD *Data) Init(name string) {
	sD.name = name
	sD.sliceData = make([]int, 0, 10)
}

// Delete метод для удаления информации по индексу
func (sD *Data) Delete(index int) {
	sD.delete(sD, index)
}

func (sD *Data) setDeleteMethod(delObj DeleteInterface) {
	sD.DeleteInterface = delObj
}

// DeleteWithoutOrder удаления без сохранения порядка
type DeleteWithoutOrder struct {
}

func (dWtO DeleteWithoutOrder) delete(data *Data, index int) {
	data.sliceData[index] = data.sliceData[len(data.sliceData)-1]
	data.sliceData[len(data.sliceData)-1] = 0
	data.sliceData = data.sliceData[:len(data.sliceData)-1]
}

// DeleteWithOrder удаления c сохранением порядка
type DeleteWithOrder struct {
}

func (dWhO DeleteWithOrder) delete(data *Data, index int) {
	copy(data.sliceData[index:], data.sliceData[index+1:])
	data.sliceData[len(data.sliceData)-1] = 0
	data.sliceData = data.sliceData[:len(data.sliceData)-1]
}
