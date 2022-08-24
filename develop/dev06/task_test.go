package main

import "testing"

func TestCut(t *testing.T) {
	d, f := "\t", ""
	s := false
	flags := flagsCmd{d: &d, f: &f, s: &s}
	fileRw := fileRows{flagsCmd: flags, rows: []string{"hello	world", "my	name", "is	dmitry"}}
	f = "1"
	fileRw.cut()
	if fileRw.cutRows[0] != "hello" && fileRw.cutRows[1] != "my" && fileRw.cutRows[2] != "is" {
		t.Fatal("Error cut")
	}
	fileRw.rows = []string{"hello,world", "my,name", "is,dmitry"}
	f = "2"
	d = ","
	fileRw.cutRows = []string{}
	fileRw.cut()
	if fileRw.cutRows[0] != "world" && fileRw.cutRows[1] != "name" && fileRw.cutRows[2] != "dmitry" {
		t.Fatal("Error cut")
	}
}
