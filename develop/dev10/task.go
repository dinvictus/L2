package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/
func readNWrite(conn net.Conn, message string) {
	_, errWrite := conn.Write([]byte(message))
	if errWrite != nil {
		os.Stderr.WriteString(errWrite.Error() + "\n")
		os.Exit(1)
	}
	conn.(*net.TCPConn).CloseWrite()

	buf, errRead := ioutil.ReadAll(conn)
	if errRead != nil {
		os.Stderr.WriteString(errRead.Error() + "\n")
		os.Exit(1)
	}
	os.Stdout.WriteString(string(buf) + "\n")
}

func connect(address, port string, timeout time.Duration, channelShutdown chan os.Signal) {
	var connDial net.Conn
	var errDial error
	timer := time.NewTimer(timeout)
	timer.Stop()
	go func() {
		select {
		case <-channelShutdown:
			connDial.Close()
			os.Exit(0)
		case <-timer.C:
			os.Stdout.WriteString("Timeout\n")
			connDial.Close()
			os.Exit(1)
		}
	}()
	for loop := true; loop; {
		connDial, errDial = net.DialTimeout("tcp", address+":"+port, timeout)
		if errDial != nil {
			os.Stderr.WriteString(errDial.Error() + "\n")
			os.Exit(1)
		}
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		timer.Reset(timeout)
		readNWrite(connDial, sc.Text()+"\n")
		timer.Stop()
		connDial.Close()
	}
}

func main() {
	channelShutdown := make(chan os.Signal, 1)
	signal.Notify(channelShutdown, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	timeout := flag.Uint("timeout", 10, "Terminating the program after reaching the connection timeout")
	flag.Parse()
	args := flag.Args()
	if len(args) == 2 {
		connect(args[0], args[1], time.Duration(*timeout)*time.Second, channelShutdown)
	} else {
		os.Stderr.WriteString("Error args\n")
		os.Exit(1)
	}
}
