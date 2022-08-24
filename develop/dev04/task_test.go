package main

import "testing"

func TestAnaram(t *testing.T) {
	str1 := "волос"
	str2 := "слово"
	anaram := isAnagram(str1, str2)
	if !anaram {
		t.Fatal("Error isanagram")
	}
	str2 = "слововолос"
	anaram = isAnagram(str1, str2)
	if !anaram {
		t.Fatal("Error isanagram")
	}
	str1 = "воло"
	anaram = isAnagram(str1, str2)
	if anaram {
		t.Fatal("Error isanagram")
	}
}
