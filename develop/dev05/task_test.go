package main

import "testing"

func TestCompare(t *testing.T) {
	var A, B, C uint = 0, 0, 0
	c, i, v, n, F := false, false, false, false, false
	flags := flagsCmd{A: &A, B: &B, C: &C, c: &c, i: &i, v: &v, n: &n, F: &F}
	fileSrch := fileSearch{flags: flags, rows: []string{"hello", "world", "my", "name", "hell"}, searchArg: "el"}
	comp := fileSrch.Compare("hello")
	if !comp {
		t.Fatal("Error compare")
	}
	F = true
	comp = fileSrch.Compare("hello")
	if comp {
		t.Fatal("Error compare")
	}
	F = false
	comp = fileSrch.Compare("hEllo")
	if comp {
		t.Fatal("Error compare")
	}
	i = true
	comp = fileSrch.Compare("hEllo")
	if !comp {
		t.Fatal("Error compare")
	}
}
