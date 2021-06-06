package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func compressSpaces(slice []byte) []byte {
	spaceBytes := []byte(" ")
	spaceByteNumber := len(spaceBytes)
	firstRune, size := utf8.DecodeRune(slice)
	for index := size; index < len(slice); {
		decodedRune, size := utf8.DecodeRune(slice[index:])
		if unicode.IsSpace(firstRune) && firstRune == decodedRune {
			copy(slice[index - size:], spaceBytes[:])
			copy(slice[index - size + spaceByteNumber:], slice[index + size:])
			slice = slice[:len(slice) - size]
		} else {
			firstRune = decodedRune
			index += size
		}
	}
	return slice
}

func main() {
	testString := "Go is a   programming\t\tlanguage\n\n"
	testSlice := []byte(testString)
	compressedSlice := compressSpaces(testSlice)
	fmt.Println(string(compressedSlice))
}