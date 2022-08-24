package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern

	Паттерн Строитель также используется, когда нужный продукт сложный и требует нескольких шагов для построения.
	В таких случаях несколько конструкторных методов подойдут лучше, чем один большой конструктор.
	При использовании пошагового построения объектов потенциальной проблемой является выдача клиенту частично построенного нестабильного продукта.
	Паттерн "Строитель" скрывает объект до тех пор, пока он не построен до конца.

*/

// Order структура для информации о заказе
type Order struct {
	paymentsInfo string
	deliveryInfo string
	itemsInfo    string
}

// BuilderInterface Интерфейс строителя
type BuilderInterface interface {
	getPaymentsInfo()
	getDeliveryInfo()
	getItemsInfo()
	getOrder() Order
}

// BuildOrder Конкретный строитель
type BuildOrder struct {
	paymentsInfo string
	deliveryInfo string
	itemsInfo    string
}

// Метод для получения информации о доставке
func (b *BuildOrder) getDeliveryInfo() {
	b.deliveryInfo = "delivery info"
}

// Метод для получения информации о предметах в заказе
func (b *BuildOrder) getItemsInfo() {
	b.itemsInfo = "items info"
}

// Метод для получения информации о платежах
func (b *BuildOrder) getPaymentsInfo() {
	b.paymentsInfo = "payment info"
}

// Метод для создания заказа с помощью строителя
func (b *BuildOrder) getOrder() Order {
	return Order{
		paymentsInfo: b.paymentsInfo,
		deliveryInfo: b.deliveryInfo,
		itemsInfo:    b.itemsInfo,
	}
}

// Director для отправки команд строителю
type Director struct {
	builder BuilderInterface
}

// CreateOrder метод для создания заказа с помощью директора
func (d *Director) CreateOrder() Order {
	d.builder.getDeliveryInfo()
	d.builder.getPaymentsInfo()
	d.builder.getItemsInfo()
	return d.builder.getOrder()
}
