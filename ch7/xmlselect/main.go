// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 214.
//!+

// Xmlselect prints the text of selected elements of an XML document.
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type ArgumentList []string

var inputPath ArgumentList
var inputIDs ArgumentList

func (p *ArgumentList) String() string {
	return fmt.Sprint(*p)
}

func (p *ArgumentList) Set(value string) error {
	if len(*p) > 0 {
		return fmt.Errorf("flag %s has been already set", *p)
	}
	*p = append(*p, strings.Split(value, " ")...)
	return nil
}

func verifyInput() error {
	if inputPath == nil && inputIDs == nil {
		errMsg := "the script awaits at least one argument: element path or ID"
		return fmt.Errorf("error: no input arguments, %s", errMsg)
	}
	return nil
}

func main() {
	flag.Var(&inputPath, "path", "A target path. Example: -path='div div h2'")
	flag.Var(&inputIDs, "ids", "A targets` ids. Example: -id='id1 id2'")
	flag.Parse()
	if err := verifyInput(); err != nil {
		fmt.Println(err)
		return
	}
	dec := xml.NewDecoder(os.Stdin)
	var nameStack []string // stack of element names
	var startElementStack []xml.StartElement
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			nameStack = append(nameStack, tok.Name.Local) // push
			startElementStack = append(startElementStack, tok)
		case xml.EndElement:
			nameStack = nameStack[:len(nameStack)-1] // pop
			startElementStack = startElementStack[:len(startElementStack)-1]
		case xml.CharData:
			if len(startElementStack) == 0 {
				continue
			}
			currentElement :=  startElementStack[len(startElementStack)-1]
			if !ElementIDIsCorrect(inputIDs, currentElement) {
				continue
			}
			if containsAll(nameStack, inputPath) {
				fmt.Printf("%s: %s\n", strings.Join(nameStack, " "), tok)
			}
		}
	}
}

func ElementIDIsCorrect(IDSlice []string, element xml.StartElement) bool {
	if inputIDs != nil {
		for _, attr := range element.Attr {
			if attr.Name.Local == "id" && containsAll(inputIDs, []string{attr.Value}) {
				return true
			}
		}
	}
	return false
}

// containsAll reports whether x contains the elements of y, in order.
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

//!-
