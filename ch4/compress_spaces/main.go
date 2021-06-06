package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func compressSpaces(slice []byte) []byte {
	firstRune, size := utf8.DecodeRune(slice)
	for index := size; index < len(slice); {
		decodedRune, size := utf8.DecodeRune(slice[index:])
		fmt.Printf("current rune '%c'\t previous rune '%c'", decodedRune, firstRune)
		if unicode.IsSpace(firstRune) && firstRune == decodedRune {
			fmt.Printf("delete symbol '%c'\n", decodedRune)
			copy(slice[index:], slice[index + size:])
			slice = slice[:len(slice) - size]
			// index -= size
			fmt.Println(index, string(slice))
			continue
		} else {
			firstRune = decodedRune
			index += size
		}
		fmt.Println(index, string(slice))
	}
	return slice
}

func main() {
	// testString := "Go is a   programming\t\tlanguage\n\n"
	testString := "Test\t\ttest"
	testSlice := []byte(testString)
	compressedSlice := compressSpaces(testSlice)
	fmt.Println(string(compressedSlice))
}