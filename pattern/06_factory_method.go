package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern

	Фабричный метод — это порождающий паттерн проектирования, который решает проблему создания различных продуктов, без указания конкретных структур продуктов.
	Фабричный метод задаёт метод, который следует использовать вместо создания объектов-продуктов.

	Паттерн Factory Method полезен, когда система должна оставаться легко расширяемой путем добавления объектов новых типов.
	Этот паттерн является основой для всех порождающих паттернов и может легко трансформироваться под нужды системы.
	По этому, если перед разработчиком стоят не четкие требования для продукта или не ясен способ организации взаимодействия между продуктами,
	то для начала можно воспользоваться паттерном Factory Method, пока полностью не сформируются все требования.
*/

type objType string

// Константные значения для типов объектов
const (
	SORT objType = "sort"
	GREP objType = "grep"
	WGET objType = "wget"
)

// CommandsInterface интерфейс для комманд
type CommandsInterface interface {
	setData(string)
	getData() string
	execute() string
	getType() string
	setType(string)
}

type utility struct {
	typeCommand string
	data        string
}

func (c *utility) setData(data string) {
	c.data = data
}

func (c utility) getType() string {
	return c.typeCommand
}

func (c utility) getData() string {
	return c.data
}

func (c *utility) setType(typeCmd string) {
	c.typeCommand = typeCmd
}

type grepUtility struct {
	utility
}

func (gU grepUtility) execute() string {
	return "Executed grep utility" + gU.data + gU.getType()
}

type sortUtility struct {
	utility
}

func (sU sortUtility) execute() string {
	return "Executed sort utility" + sU.data + sU.getType()
}

type wgetUtility struct {
	utility
}

func (wU wgetUtility) execute() string {
	return "Executed wget utility" + wU.data + wU.getType()
}

// CreateCommand фабричный метод для создания объекта по типу
func CreateCommand(cmdType objType, data string) CommandsInterface {
	var cmd CommandsInterface
	switch cmdType {
	case GREP:
		cmd = &grepUtility{}
	case SORT:
		cmd = &sortUtility{}
	case WGET:
		cmd = &wgetUtility{}
	default:
		return nil
	}
	cmd.setData(data)
	cmd.setType(string(cmdType))
	return cmd
}
