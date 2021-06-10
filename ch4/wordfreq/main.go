package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error. The script needs a file path as an input argument")
		return
	}
	fileBytes, err := ioutil.ReadFile(os.Args[1])
	fileReader := bytes.NewReader(fileBytes)
	fileScanner := bufio.NewScanner(fileReader)
	fileScanner.Split(bufio.ScanWords)
	for fileScanner.Scan() {
		word := fileScanner.Text()
	}
}