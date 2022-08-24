package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func isAnagram(str1, str2 string) bool {
	mapSymb1 := make(map[rune]struct{})
	mapSymb2 := make(map[rune]struct{})
	if str1 == str2 {
		return false
	}
	for _, s := range str1 {
		mapSymb1[s] = struct{}{}
	}
	for _, s := range str2 {
		mapSymb2[s] = struct{}{}
	}
	if len(mapSymb1) != len(mapSymb2) {
		return false
	}
	for _, s := range str2 {
		if _, ok := mapSymb1[s]; !ok {
			return false
		}
	}
	for _, s := range str1 {
		if _, ok := mapSymb2[s]; !ok {
			return false
		}
	}
	return true
}

func searchAnagram(dictionary []string) map[string][]string {
	mapAnagram := make(map[string][]string)
	mapAdded := make(map[string]struct{})
	for i := range dictionary {
		dictionary[i] = strings.ToLower(dictionary[i])
	}
	for i := range dictionary {
		for j := range dictionary {
			if _, ok := mapAdded[dictionary[i]]; !ok && isAnagram(dictionary[i], dictionary[j]) {
				mapAnagram[dictionary[i]] = append(mapAnagram[dictionary[i]], dictionary[j])
				mapAdded[dictionary[j]] = struct{}{}
			}
		}
		mapAdded[dictionary[i]] = struct{}{}
		sort.Strings(mapAnagram[dictionary[i]])
	}
	return mapAnagram
}

func main() {
	words := []string{"пятак", "пятка", "тяПка", "листок", "слиток", "столик", "мама", "амам", "маам", "листок", "столикс", "test", "tset"}
	fmt.Println(searchAnagram(words))
}
