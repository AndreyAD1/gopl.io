// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 73.

// Comma prints its argument numbers with a comma at each power of 1000.
//
// Example:
// 	$ go build gopl.io/ch3/comma
//	$ ./comma 1 12 123 1234 1234567890
// 	1
// 	12
// 	123
// 	1,234
// 	1,234,567,890
//
package main

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

//!+
// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	integerPart := s
	fractionalPart := ""
	for runeIndex, rune := range s {
		if !unicode.IsDigit(rune) {
			if runeIndex == 0 && rune != '-' {
				return "Error. String should contain a decimal number."
			}
			if runeIndex != 0 && rune != '.' {
				return "Error. String should contain a decimal number."
			}
			if rune == '.' {
				integerPart = s[:runeIndex]
				fractionalPart = s[runeIndex:]
				break
			}
		}
	}

	if len(integerPart) <= 3 {
		return s
	}

	firstSeparatorIndex := len(integerPart) % 3
	if firstSeparatorIndex == 0 {
		firstSeparatorIndex = 3
	}
	buffer := bytes.NewBufferString(integerPart[:firstSeparatorIndex])
	for index := firstSeparatorIndex; len(integerPart) - index > 1; {
		buffer.WriteString(",")
		buffer.WriteString(integerPart[index:index + 3])
		index += 3
	}
	buffer.WriteString(fractionalPart)
	return buffer.String()
}

//!-
