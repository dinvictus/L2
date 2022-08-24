package pattern

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern

	Фасад — это структурный паттерн, который предоставляет простой (но урезанный) интерфейс к сложной системе объектов, библиотеке или фреймворку.
	Кроме того, что Фасад позволяет снизить общую сложность программы, он также помогает вынести код, зависимый от внешней системы в единственное место.

	Пример реального использования приведён в реализации ниже
*/

// PurchaseFacade Фасад для совершения платежа
type PurchaseFacade struct {
	account      *Account
	server       *HTTPServer
	bankSender   *Bank
	BankReciever *Bank
}

// Функция для создания васада
func createFacade(bankIDSender, bankIDReciever string) *PurchaseFacade {
	purchaseFacade := PurchaseFacade{
		account:      new(Account),
		server:       new(HTTPServer),
		bankSender:   createBank(bankIDSender),
		BankReciever: createBank(bankIDReciever),
	}
	return &purchaseFacade
}

// Account структура для пользовательского аккаунта
type Account struct {
}

// Метод для верификации аккаунта пользователя
func (aC Account) verificationAccount(userid string) bool {
	return userid != ""
}

// Метод совершения платежа, обёрнутый фасадом
func (pF PurchaseFacade) purchase(purchaseData, userid string) bool {
	ver := pF.account.verificationAccount(userid)
	if !ver {
		return false
	}
	statusCodeSender := pF.server.sendRequest(*pF.bankSender, purchaseData)
	statusCodeReciever := pF.server.sendRequest(*pF.BankReciever, purchaseData)
	if statusCodeReciever == 200 && statusCodeSender == 200 {
		return true
	}
	return false
}

// HTTPServer Сервер для посылания запросов к банкам
type HTTPServer struct {
}

func (h HTTPServer) sendRequest(bank Bank, purchaseInfo string) int {
	if bank.bankID == "" || purchaseInfo == "" {
		return 404
	}
	return 200
}

// Bank структура для банков
type Bank struct {
	bankID string
}

// Функция для создания нового банка
func createBank(bankID string) *Bank {
	return &Bank{bankID: bankID}
}
