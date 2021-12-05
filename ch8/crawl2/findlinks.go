// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

//!-sema

type WorkList struct {
	Urls []string
	UrlLevel int
}

//!+
func main() {
	worklists := make(chan WorkList)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() { worklists <- WorkList{os.Args[1:], 0} }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		worklist := <-worklists
		if worklist.UrlLevel > 2 {
			continue
		}
		for _, link := range worklist.Urls {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklists <- WorkList{crawl(link), worklist.UrlLevel + 1}
				}(link)
			}
		}
	}
}

//!-
