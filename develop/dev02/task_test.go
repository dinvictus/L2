package main

import "testing"

func TestStr1(t *testing.T) {
	str := "a4bc2d5e"
	expected := "aaaabccddddde"
	recieved, err := UnpackingString(str)
	if err != nil {
		t.Fatal(err)
	}
	if recieved != expected {
		t.Fatal("Incorrect func work")
	}
}

func TestStr2(t *testing.T) {
	str := "aa4bc2d5ee7"
	expected := "aaaaabccdddddeeeeeeee"
	recieved, err := UnpackingString(str)
	if err != nil {
		t.Fatal(err)
	}
	if recieved != expected {
		t.Fatal("Incorrect func work")
	}
}

func TestStr3(t *testing.T) {
	str := "54534534"
	_, err := UnpackingString(str)
	if err == nil {
		t.Fatal("Incorrect func work")
	}
}

func TestStr4(t *testing.T) {
	str := `qwe4\\5\710`
	expected := `qweeee\\\\\7777777777`
	recieved, err := UnpackingString(str)
	if err != nil {
		t.Fatal(err)
	}
	if recieved != expected {
		t.Fatal("Incorrect func work")
	}
}

func TestStr5(t *testing.T) {
	str := `qwe\\\\\\5\10`
	expected := `qwe\\\\\\\`
	recieved, err := UnpackingString(str)
	if err != nil {
		t.Fatal(err)
	}
	if recieved != expected {
		t.Fatal("Incorrect func work")
	}
}

func TestStr6(t *testing.T) {
	str := `5q\55q5\`
	expected := "q55555qqqqq"
	recieved, err := UnpackingString(str)
	if err != nil {
		t.Fatal(err)
	}
	if recieved != expected {
		t.Fatal("Incorrect func work")
	}
}
