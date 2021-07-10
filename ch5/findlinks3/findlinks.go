// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 139.

// Findlinks3 crawls the web, starting with the URLs on the command line.
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"gopl.io/ch5/links"
)

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

//!-breadthFirst

//!+crawl
func crawl(downloadURL string) []string {
	fmt.Println(downloadURL)
	list, err := links.Extract(downloadURL)
	if err != nil {
		log.Print(err)
	}
	for _, link := range list {
		parsedURL, err := url.Parse(link)
		if err != nil {
			fmt.Printf("Invalid link: %s\n", link)
			continue
		}
		dirPath := parsedURL.Host
		if len(parsedURL.Path) > 1 {
			dirPath += parsedURL.Path
		}
		filePath := dirPath + ".html"
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				fmt.Printf("Can not create the directory: %s\n", dirPath)
				continue
			}
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("Can not create the file: %s\n", filePath)
				continue
			}
			defer file.Close()
			_, err = file.WriteString("TEST TEST TEST")
			if err != nil {
				fmt.Printf("Can not write to the file: %s\n", filePath)
			}
		}
	}
	return list
}

//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

//!-main
