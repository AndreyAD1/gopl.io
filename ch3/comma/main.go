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
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

//!+
// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	s1 := s[len(s)-3:]
	fmt.Println("initial string", s1)
	for index := 3; len(s) - index >= 1; {
		if len(s) - index < 3 {
			s1 = s[:len(s) - index] + "," + s1
			break
		}
		s1 = s[len(s) - index - 3:len(s) - index] + "," + s1
		fmt.Printf("string per index %d: %s\n", index, s1)
		index += 3
	}
	return s1
}

//!-
