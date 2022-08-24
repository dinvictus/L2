package main

import "testing"

func TestPipe(t *testing.T) {
	expectOut := "   how 2m\r\n  what 5b\r\nare 1g\r\nare 6k\r\ndoing 50m\r\nhello 3b\r\nyou 100g  \r\nyou 10m\r\n"
	out, err := pipe("echo file1.txt | sort")
	if err != nil {
		t.Fatal("Error pipe command")
	}
	if out != expectOut {
		t.Fatal("Error pipe command")
	}
}
