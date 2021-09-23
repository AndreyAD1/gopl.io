package main

import (
	"bufio"
	"fmt"
)

type WordCounter int

func (wc *WordCounter) Write(data []byte) (n int, err error) {
	for start := 0;; {
		wordLength, word, err := bufio.ScanWords(data[start:], true)
		if err != nil {
			return 0, err
		}
		*wc++
		if len(data[start:]) == len(word) {
			break
		}
		start += wordLength
	}
	return len(data), nil
}

func main() {
	var wordCounter WordCounter
	word := "first"
	word2 := "second"
	fmt.Fprintf(&wordCounter, "Some string with words: %s, %s", word, word2)
	fmt.Println(wordCounter)
}