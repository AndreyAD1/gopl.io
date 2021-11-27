// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup

func echo(c net.Conn, shout string, delay time.Duration) {
	defer waitGroup.Done()
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

//!+
func handleConn(c net.Conn) {
	TCPConn := c.(*net.TCPConn)
	input := bufio.NewScanner(TCPConn)
	for input.Scan() {
		waitGroup.Add(1)
		// fmt.Println("work group counter ", waitGroup)
		go echo(TCPConn, input.Text(), 1*time.Second)
	}
	go func() {
		fmt.Println("Run waiting goroutine")
		waitGroup.Wait()
		TCPConn.CloseWrite()
	}()
	// NOTE: ignoring potential errors from input.Err()
	TCPConn.CloseRead()
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
