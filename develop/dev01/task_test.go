package main

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestTime(t *testing.T) {
	strTime := GetCurTime()
	strTime = strTime[:len(strTime)-1]
	strTimeSplit := strings.Split(strTime, ":")
	if len(strTimeSplit) != 3 {
		t.Fatal("Error time")
	}
	for _, el := range strTimeSplit {
		if utf8.RuneCountInString(el) != 2 {
			t.Fatal("Error time: ", utf8.RuneCountInString(el))
		} else if !unicode.IsDigit(rune(el[0])) && !unicode.IsDigit(rune(el[1])) {
			t.Fatal("Error time: " + el)
		}
	}
}
