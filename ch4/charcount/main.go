// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 97.
//!+

// Charcount computes counts of Unicode characters.
package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)    // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters
	stdinScanner := bufio.NewScanner(os.Stdin)
	for {
		stdinScanner.Scan()
		inputText := stdinScanner.Text()
		if len(inputText) == 0 {
			break
		}
		for i := 0; i < len(inputText); {
			run, byte_number := utf8.DecodeRuneInString(inputText[i:])
			if run == utf8.RuneError {
				break
			}
			if run == unicode.ReplacementChar && byte_number == 1 {
				invalid++
				continue
			}
			counts[run]++
			utflen[byte_number]++
			i += byte_number
		}
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

//!-
