package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"
)

func TestTelnet(t *testing.T) {
	expectRes := "Hello, world!"
	CONN_HOST := "127.0.0.1"
	CONN_PORT := "3333"
	CONN_TYPE := "tcp"
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		t.Fatal("Error telnet")
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	var res string
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf, read_err := ioutil.ReadAll(conn)
		if read_err != nil {
			return
		}
		res += string(buf)
	}(wg)
	connDial, errDial := net.DialTimeout(CONN_TYPE, CONN_HOST+":"+CONN_PORT, time.Duration(10)*time.Second)
	if errDial != nil {
		t.Fatal("Error telnet")
	}
	readNWrite(connDial, "Hello, world!")
	wg.Wait()
	if res != expectRes {
		t.Fatal("Error telnet")
	}
}
