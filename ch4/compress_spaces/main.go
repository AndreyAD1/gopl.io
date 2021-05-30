package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func compressSpaces(slice []byte) []byte {
	for index := 0; index < len(slice); {
		spaceSlice := []byte(" ")
		_, spaceSize := utf8.DecodeRune(spaceSlice)
		decodedRune, size := utf8.DecodeRune(slice)
		if unicode.IsSpace(decodedRune) {
			copy(slice[index:], spaceSlice)
			copy(slice[index + spaceSize:], slice[index + size:])
			slice = slice[:len(slice) - size]
			index -= size + spaceSize
		} else {
			index += size
		}
	}

	firstRune, size := utf8.DecodeRune(slice)
	for index := size; index < len(slice); {
		decodedRune, size := utf8.DecodeRune(slice[index:])
		fmt.Printf("current rune '%c'\t previous rune '%c'", decodedRune, firstRune)
		if unicode.IsSpace(firstRune) && firstRune == decodedRune {
			fmt.Printf("delete symbol '%c'\n", decodedRune)
			copy(slice[index:], slice[index + size:])
			slice = slice[:len(slice) - size]
			index -= size
			continue
		} else {
			firstRune = decodedRune
			index += size
		}
		fmt.Println(string(slice))
	}
	return slice
}

func main() {
	testString := "Go is a   programming \t\t   language\n\n"
	testSlice := []byte(testString)
	compressedSlice := compressSpaces(testSlice)
	fmt.Println(string(compressedSlice))
}