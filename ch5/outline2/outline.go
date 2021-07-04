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

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
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
