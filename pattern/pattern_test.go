package pattern

import (
	"testing"
	"time"
)

func TestFacade(t *testing.T) {
	purchaseFacade := createFacade("1", "2")
	ok := purchaseFacade.purchase("2000", "user1")
	if !ok {
		t.Fatal("Error pattern facade")
	}
}

func TestBuilder(t *testing.T) {
	expectDeliveryInfo := "delivery info"
	expectPaymentInfo := "payment info"
	expectItemsInfo := "items info"
	director := Director{&BuildOrder{}}
	order := director.CreateOrder()
	if order.deliveryInfo != expectDeliveryInfo || order.itemsInfo != expectItemsInfo || order.paymentsInfo != expectPaymentInfo {
		t.Fatal("Error pattern builder")
	}
}

func TestVisitor(t *testing.T) {
	expectDeliveryInfo := " city, street Processed delivery"
	expectPaymentInfo := " sale, count Processed payment"
	expectItemsInfo := " name, price Processed items"
	delInfo := DeliveryInfo{data: " city, street "}
	payInfo := PaymentInfo{data: " sale, count "}
	itInfo := ItemsInfo{data: " name, price "}

	processData := ProcessingData{}

	processDataDelivery := processData.visitForDelivery(&delInfo)
	processDataPayment := processData.visitForPayment(&payInfo)
	processDataItems := processData.visitForItems(&itInfo)
	if processDataDelivery != expectDeliveryInfo || processDataPayment != expectPaymentInfo || processDataItems != expectItemsInfo {
		t.Fatal("Error pattern visitor")
	}

}

func TestCommand(t *testing.T) {
	expectOut := "Server starting...\nServer restarted\nServer stopped\n"
	httpServer := ServerHTTP{}
	serverStartCommand := StartHTTPServerCommand{httpServer: &httpServer}
	serverStopCommand := StopHTTPServerCommand{httpServer: &httpServer}
	serverRestartCommand := RestartHTTPServerCommand{httpServer: &httpServer}
	commandStorage := CommandStorage{}
	commandStorage.SaveCommand(&serverStartCommand)
	commandStorage.SaveCommand(&serverRestartCommand)
	commandStorage.SaveCommand(&serverStopCommand)
	out := commandStorage.ExecuteCommands()
	if out != expectOut {
		t.Fatal("Error pattern command")
	}
}

func TestChainOfResp(t *testing.T) {
	expectVerification := "Verification bank sender complete\nVerification bank reciever complete\nVerification payment gateway complete\n"
	payment := Payment{false, false, false}
	handlers := &BankSender{next: &BankReciever{next: &PaymentGateway{}}}

	verification := handlers.verification(&payment)

	if verification != expectVerification {
		t.Fatal("Error pattern chain of resp")
	}
}

func TestFactoryMethod(t *testing.T) {
	expectExecute := "Executed wget utilityData for wgetwgetExecuted sort utilityData for sortsortExecuted grep utilityData for grepgrep"
	commands := []CommandsInterface{
		CreateCommand(WGET, "Data for wget"),
		CreateCommand(SORT, "Data for sort"),
		CreateCommand(GREP, "Data for grep"),
	}
	var resultExecute string
	for _, cmd := range commands {
		resultExecute += cmd.execute()
	}
	if resultExecute != expectExecute {
		t.Fatal("Error pattern factory method")
	}
}

func TestStrategy(t *testing.T) {
	expectDeleteWithoutOrder := []int{0, 1, 2, 3, 4, 9, 6, 7, 8}
	expectDeleteWithOrder := []int{0, 1, 2, 3, 4, 6, 7, 8}
	deleteWithoutOrder := DeleteWithoutOrder{}
	deleteWithOrder := DeleteWithOrder{}
	data := Data{}
	data.Init("Data storage")
	for i := 0; i < 10; i++ {
		data.Add(i)
	}
	data.setDeleteMethod(deleteWithoutOrder)
	data.Delete(5)
	for i := 0; i < len(expectDeleteWithoutOrder); i++ {
		if expectDeleteWithoutOrder[i] != data.sliceData[i] {
			t.Fatal("Error pattern strategy")
		}
	}
	data.setDeleteMethod(deleteWithOrder)
	data.Delete(5)
	for i := 0; i < len(expectDeleteWithOrder); i++ {
		if expectDeleteWithOrder[i] != data.sliceData[i] {
			t.Fatal("Error pattern strategy")
		}
	}
}

func TestState(t *testing.T) {
	expectResult := "Server startedServer restartedserver restartingServer stopped"
	var result string
	server := newServer()
	resStart, errStart := server.start()
	if errStart != nil {
		t.Fatal("Error pattern state")
	}
	result += resStart
	var errStop error
	go func() {
		time.Sleep(1 * time.Second)
		_, errStop = server.stop()
	}()
	resRestart, errRestart := server.restart()
	if errStop == nil {
		t.Fatal("Error pattern state")
	}
	if errRestart != nil {
		t.Fatal("Error pattern state")
	}
	result += resRestart
	result += errStop.Error()
	time.Sleep(time.Second * 6)
	resStop, errStop := server.stop()
	if errStop != nil {
		t.Fatal("Error pattern state")
	}
	result += resStop
	if result != expectResult {
		t.Fatal("Error pattern state")
	}
}
