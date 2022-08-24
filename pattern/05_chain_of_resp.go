package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

	Цепочка обязанностей — это поведенческий паттерн, позволяющий передавать запрос по цепочке потенциальных обработчиков, пока один из них не обработает запрос.
	Избавляет от жёсткой привязки отправителя запроса к его получателю, позволяя выстраивать цепь из различных обработчиков динамически.
*/

// Payment Структура для хранения информации о верификации платежа
type Payment struct {
	bankSenderVerification     bool
	bankRecieverVerification   bool
	paymentGatewayVerification bool
}

// Handler интерфейс обработчика
type Handler interface {
	verification(*Payment) string
}

// BankSender структура для банка отправителя
type BankSender struct {
	next Handler
}

func (bS *BankSender) verification(p *Payment) (out string) {
	p.bankSenderVerification = true
	out += "Verification bank sender complete\n"
	if bS.next != nil {
		out += bS.next.verification(p)
	}
	return
}

// BankReciever структура для банка получателя
type BankReciever struct {
	next Handler
}

func (bR *BankReciever) verification(p *Payment) (out string) {
	if !p.bankSenderVerification {
		out += "The sender's bank rejected verification\n"
		return
	}
	p.bankRecieverVerification = true
	out += "Verification bank reciever complete\n"
	if bR.next != nil {
		out += bR.next.verification(p)
	}
	return
}

// PaymentGateway структура для платёженого шлюза
type PaymentGateway struct {
	next Handler
}

func (pG *PaymentGateway) verification(p *Payment) (out string) {
	if !p.bankSenderVerification || !p.bankRecieverVerification {
		out += "The sender's or reciever's bank rejected verification\n"
		return
	}
	p.paymentGatewayVerification = true
	out += "Verification payment gateway complete\n"
	if pG.next != nil {
		out += pG.next.verification(p)
	}
	return
}
