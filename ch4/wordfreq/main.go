package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error. The script needs a file path as an input argument")
		return
	}
	fileBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Can not read the file %s", os.Args[1])
		return
	}
	fileReader := bytes.NewReader(fileBytes)
	fileScanner := bufio.NewScanner(fileReader)
	fileScanner.Split(bufio.ScanWords)
	var maxWordLength int
	wordInfo := make(map[string][2]int)
	for fileScanner.Scan() {
		word := fileScanner.Text()
		lowerWord := strings.ToLower(word)
		wordLenght := utf8.RuneCountInString(lowerWord)
		if wordLenght > maxWordLength {
			maxWordLength = wordLenght
		}
		wordNumber := wordInfo[lowerWord][0]
		wordInfo[lowerWord] = [2]int{wordNumber + 1, wordLenght}
	}
	if fileScanner.Err() != nil {
		fmt.Println("Error occured while scanner was reading the file")
		return
	}
	for word, wordInfo := range wordInfo {
		wordNumber, wordLength := wordInfo[0], wordInfo[1]
		additionalRuneNumber := maxWordLength - wordLength
		for r := 0; r < additionalRuneNumber; r++ {
			word += " "
		}
		fmt.Printf("%s\t%d\n", word, wordNumber)
	}
}