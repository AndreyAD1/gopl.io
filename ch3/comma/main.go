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
	"strings"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

//!+
// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	var buffer bytes.Buffer

	if strings.HasPrefix(s, "-") {
		if strings.Count(s, "-") != 1 {
			return "Error. String should not contain a minus in a middle."
		}
		buffer.WriteRune('-')
		s = s[1:]
	}

	splittedString := strings.Split(s, ".")
	if len(splittedString) > 2 {
		return "Error. String should contain only one point."
	}

	integerPart := splittedString[0]
	fractionalPart := ""
	if len(splittedString) == 2 {
		fractionalPart = splittedString[1]
	}

	firstSeparatorIndex := len(integerPart) % 3
	if firstSeparatorIndex == 0 {
		firstSeparatorIndex = 3
	}
	buffer.WriteString(integerPart[:firstSeparatorIndex])
	for index := firstSeparatorIndex; len(integerPart) - index > 1; {
		buffer.WriteString(",")
		buffer.WriteString(integerPart[index:index + 3])
		index += 3
	}
	if len(fractionalPart) > 0 {
		buffer.WriteString(".")
		buffer.WriteString(fractionalPart)
	}
	return buffer.String()
}

//!-
