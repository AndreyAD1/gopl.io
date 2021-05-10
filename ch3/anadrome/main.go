package main

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Expect two strings as script arguments")
		return
	}

	if isAnadrome(os.Args[1], os.Args[2]) {
		fmt.Println("These strings are anadromes")
	} else {
		fmt.Println("These strings are not anadromes")
	}
}

func isAnadrome(firstString string, secondString string) bool {
	var buffer bytes.Buffer
	runeBytes := []byte(secondString)
	for len(runeBytes) > 0 {
		run, size := utf8.DecodeLastRune(runeBytes)
		buffer.WriteRune(run)
		runeBytes = runeBytes[:len(runeBytes) - size]
	}

	return firstString == buffer.String()
}