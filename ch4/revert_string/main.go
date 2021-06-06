package main

import (
	"fmt"
	"unicode/utf8"
)

func reverse(slice []byte) {
	for i, j := 0, len(slice); i < j; {
		shortenSlice := slice[i: j]
		_, firstRuneSize := utf8.DecodeRune(shortenSlice)
		_, lastRuneSize := utf8.DecodeLastRune(shortenSlice)
		sizeDiff := firstRuneSize - lastRuneSize
		if sizeDiff >= 0 {
			firstRune := make([]byte, firstRuneSize)
			copy(firstRune, shortenSlice[:firstRuneSize])
			copy(shortenSlice, shortenSlice[len(shortenSlice) - lastRuneSize:])
			copy(shortenSlice[lastRuneSize:], shortenSlice[firstRuneSize:])
			copy(shortenSlice[len(shortenSlice) - firstRuneSize:], firstRune)
		} else {
			lastRune := make([]byte, lastRuneSize)
			copy(lastRune, shortenSlice[len(shortenSlice) - lastRuneSize:])
			copy(shortenSlice[len(shortenSlice) - firstRuneSize:], shortenSlice[:firstRuneSize])
			copy(shortenSlice[-sizeDiff:len(shortenSlice) - lastRuneSize], shortenSlice[:lastRuneSize])
			copy(shortenSlice[:lastRuneSize], lastRune)
		}
		i += lastRuneSize
		j -= firstRuneSize
	}
}

func main() {
	testSlice := []byte("a")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("ab")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("abc")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("abcd")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("Ю")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("Юл")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("Юля")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("Тест тест   тест")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)

	testSlice = []byte("Сxмесь QЯ")
	reverse(testSlice)
	fmt.Printf("%s\n", testSlice)
}