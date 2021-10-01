// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 133.

// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"unicode/utf8"

	"golang.org/x/net/html"
)

const exampleHTML = "<head>\n    <title>Example Domain</title>\n\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"Content-type\" content=\"text/html; charset=utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <style type=\"text/css\">\n    body {\n        background-color: #f0f0f2;\n        margin: 0;\n        padding: 0;\n        font-family: -apple-system, system-ui, BlinkMacSystemFont, \"Segoe UI\", \"Open Sans\", \"Helvetica Neue\", Helvetica, Arial, sans-serif;\n        \n    }\n    div {\n        width: 600px;\n        margin: 5em auto;\n        padding: 2em;\n        background-color: #fdfdff;\n        border-radius: 0.5em;\n        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);\n    }\n    a:link, a:visited {\n        color: #38488f;\n        text-decoration: none;\n    }\n    @media (max-width: 700px) {\n        div {\n            margin: 0 auto;\n            width: auto;\n        }\n    }\n    </style>    \n<style>@media print {#ghostery-tracker-tally {display:none !important}}</style></head>\n\n<body>\n<div>\n    <h1>Example Domain</h1>\n    <p>This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.</p>\n    <p><a href=\"https://www.iana.org/domains/example\">More information...<a href=\"https://www.iana.org/domains/example2\"></a></p>\n</div>\n\n\n</body>"

type StringReader string

func (s StringReader) Read(p []byte) (int, error) {
	if len(s) == 0 {
		return 0, io.EOF
	}
	
	readByteNumber := 0
	for {
		run, runeSize := utf8.DecodeRuneInString(string(s))
		if run == utf8.RuneError {
			return readByteNumber, io.EOF
		}
		readByteNumber += runeSize
		if len(p) < readByteNumber {
			return readByteNumber - runeSize, nil
		}
		copy(p[readByteNumber - runeSize:], s[:runeSize])
		s = s[runeSize:]
	}
}

func main() {
	commandError := "The script requires either 'url' or 'string' command"
	if len(os.Args) < 2 {
		fmt.Println(commandError)
	}
	switch os.Args[1] {
		case "url":
			for _, url := range os.Args[2:] {
				resp, err := http.Get(url)	
				if err != nil {
					fmt.Println(err)
					break
				}
				err = outline(resp.Body)
				if err != nil {
					fmt.Println(err)
				}
				resp.Body.Close()
			}
		case "string":
			stringReader := StringReader(exampleHTML)
			err := outline(stringReader)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println(commandError)
	}

}

func outline(responseBody io.Reader) error {
	doc, err := html.Parse(responseBody)
	if err != nil {
		return err
	}

	//!+call
	forEachNode(doc, startElement, endElement)
	//!-call

	return nil
}

//!+forEachNode
// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

//!-forEachNode

//!+startend
var depth int

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		attributeString := ""
		for _, attr := range n.Attr {
			attributeString += fmt.Sprintf(" %s=\"%s\"", attr.Key, attr.Val)
		}
		if n.FirstChild != nil {
			fmt.Printf("%*s<%s%s>\n", depth*2, "", n.Data, attributeString)
			depth++
		} else {
			fmt.Printf("%*s<%s%s/>\n", depth*2, "", n.Data, attributeString)
		}
	}
	if n.Type == html.TextNode {
		if !regexp.MustCompile(`^[\s]+$`).MatchString(n.Data) {
			fmt.Printf("%*s%s\n", depth*2, "", n.Data)
		}
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode && n.FirstChild != nil {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}

//!-startend
