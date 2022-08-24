package main

import "testing"

func TestMonthCompare(t *testing.T) {
	month1 := "Февраль"
	month2 := "октябрь"
	comp := compareMonth(month1, month2)
	if !comp {
		t.Fatal("Error compare month")
	}
	month1 = "Декабрь"
	month2 = "март"
	comp = compareMonth(month1, month2)
	if comp {
		t.Fatal("Error compare month")
	}
}

func TestCompareLess(t *testing.T) {
	h := ""
	n, r, u, M, b, c := false, false, false, false, false, false
	var k uint = 0
	flags := flagsCmd{h: &h, n: &n, r: &r, u: &u, M: &M, b: &b, c: &c, k: &k}
	str1 := "hello"
	str2 := "world"
	comp := compareFlagsLess(str1, str2, flags)
	if !comp {
		t.Fatal("Error compare less")
	}
	n = true
	str1 = "55"
	str2 = "155"
	comp = compareFlagsLess(str1, str2, flags)
	if !comp {
		t.Fatal("Error compare less")
	}
	n = false
	comp = compareFlagsLess(str1, str2, flags)
	if comp {
		t.Fatal("Error compare less")
	}
	str1 = "hello world"
	str2 = "meow cat"
	k = 2
	comp = compareFlagsLess(str1, str2, flags)
	if comp {
		t.Fatal("Error compare less")
	}
	k = 0
	str1 = "hello"
	str2 = "     world"
	comp = compareFlagsLess(str1, str2, flags)
	if comp {
		t.Fatal("Error compare less")
	}
	b = true
	comp = compareFlagsLess(str1, str2, flags)
	if !comp {
		t.Fatal("Error compare less")
	}
}
