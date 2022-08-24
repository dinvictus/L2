package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern

	Посетитель — это поведенческий паттерн, который позволяет добавить новую операцию для иерархии структур, не изменяя код этих структур.
*/

// Интерфейс для информации
type info interface {
	getData() string
	setData(string)
	getTypeInfo() string
	accept(visitor)
}

// Интерфейс для посетителя
type visitor interface {
	visitForDelivery(*DeliveryInfo)
	visitForPayment(*PaymentInfo)
	visitForItems(*ItemsInfo)
}

// DeliveryInfo структура для хранения информации о доставке
type DeliveryInfo struct {
	data string
}

func (d *DeliveryInfo) setData(data string) {
	d.data = "delivery data" + data
}

func (d DeliveryInfo) getData() string {
	return d.data
}

func (d DeliveryInfo) getTypeInfo() string {
	return "delivery"
}

func (d *DeliveryInfo) accept(v visitor) {
	v.visitForDelivery(d)
}

// PaymentInfo структура для хранения информации о платеже
type PaymentInfo struct {
	data string
}

func (p *PaymentInfo) setData(data string) {
	p.data = "payment data" + data
}

func (p PaymentInfo) getData() string {
	return p.data
}

func (p PaymentInfo) getTypeInfo() string {
	return "payment"
}

func (p *PaymentInfo) accept(v visitor) {
	v.visitForPayment(p)
}

// ItemsInfo структура для хранения информации о предметах
type ItemsInfo struct {
	data string
}

func (it *ItemsInfo) setData(data string) {
	it.data = "items data" + data
}

func (it ItemsInfo) getData() string {
	return it.data
}

func (it ItemsInfo) getTypeInfo() string {
	return "items"
}

func (it *ItemsInfo) accept(v visitor) {
	v.visitForItems(it)
}

// ProcessingData посетитель
type ProcessingData struct {
}

func (pD *ProcessingData) visitForDelivery(d *DeliveryInfo) string {
	return d.data + "Processed delivery"
}

func (pD *ProcessingData) visitForPayment(p *PaymentInfo) string {
	return p.data + "Processed payment"
}

func (pD *ProcessingData) visitForItems(it *ItemsInfo) string {
	return it.data + "Processed items"
}
