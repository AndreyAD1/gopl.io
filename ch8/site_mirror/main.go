package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"gopl.io/ch5/links"
)

var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := links.Extract(url)
	<-tokens

	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	flag.Parse()
	inputURLs := flag.Args()
	var initialDomains []string
	for _, inputURL := range inputURLs {
		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			log.Printf("Invalid input URL %s: %v", inputURL, err)
			continue
		}
		inputHost := parsedURL.Hostname()
		if inputHost == "" {
			log.Fatalf("Can not determine the host for %s", inputURL)
		}
		initialDomains = append(initialDomains, inputHost)
	}
	worklists := make(chan []string)
	var n int
	n++
	go func() { worklists <- flag.Args() }()
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		worklist := <-worklists
		for _, link := range worklist {
			if !seen[link] && hostIsInitial(initialDomains, link) {
				seen[link] = true
				n++
				go func(link string) {
					worklists <- crawl(link)
				}(link)
			}
		}
	}
}

func hostIsInitial(initialHosts []string, link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		log.Printf("Get invalid URL from the worklist: %s", link)
		return false
	}
	linkHost := parsedURL.Hostname()
	linkHostWithoutWWW := strings.TrimPrefix(linkHost, "www.")
	for _, initialHost := range initialHosts {
		initialHostWithoutWWW := strings.TrimPrefix(initialHost, "www.")
		if linkHostWithoutWWW == initialHostWithoutWWW {
			return true
		}
	}
	return false
}