// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 139.

// Findlinks3 crawls the web, starting with the URLs on the command line.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
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
		responseBody, err := getPage(link)
		if err != nil {
			fmt.Printf("Can not download the page: %s\n", link)
		}
		defer responseBody.Close()
		parsedURL, err := url.Parse(link)
		if err != nil {
			fmt.Printf("Invalid link: %s\n", link)
			continue
		}
		file, err := getFile(parsedURL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(responseBody)
		for scanner.Scan() {
			_, err := file.Write(scanner.Bytes())
			if err != nil {
				fmt.Printf("Can not write to the file %s: %s\n", file.Name(), err)
				break
			}
		}
	}
	return list
}

func getPage(link string) (io.ReadCloser, error) {

	response, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("can not download the URL %s: %s", link, err)
	}
	
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receive an HTTP error %s: %s", link, response.Status)
	}
	return response.Body, nil
}

func getFile(parsedURL *url.URL) (*os.File, error) {
	dirPath := "downloaded_pages/" + parsedURL.Host
	if len(parsedURL.Path) > 1 {
		dirPath += parsedURL.Path
	}
	filePath := dirPath + ".html"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return nil, fmt.Errorf("can not create the directory: %s", dirPath)
		}
		file, err := os.Create(filePath)
		if err != nil {
			err := fmt.Errorf("can not create the file: %s", filePath)
			return nil, err
		}
		return file, err
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		err = fmt.Errorf("can not open the file: %s", filePath)
		return nil, err
	}
	return file, err
}


//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

//!-main
