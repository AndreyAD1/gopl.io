// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	error_message := "The script has two positional arguments: " +
		"output file name and at least one URL"
	if len(os.Args) < 3 {
		fmt.Println(error_message)
		return
	}

	for _, url := range os.Args[2:] {
		fmt.Printf("start a goroutine for %s", url)
		go fetch(url, ch) // start a goroutine
	}
	var result []string
	for range os.Args[2:] {
		result = append(result, <-ch) // receive from channel ch
	}
	elapsed_time := fmt.Sprintf("%.2fs elapsed\n", time.Since(start).Seconds())
	result = append(result, elapsed_time)
	output_filename := os.Args[1]
	result_string := strings.Join(result, "\n")
	err := ioutil.WriteFile(output_filename, []byte(result_string), 0644)
	if err != nil {
		fmt.Printf("Can not write the result to the file %s", output_filename)
		return
	}
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

//!-
