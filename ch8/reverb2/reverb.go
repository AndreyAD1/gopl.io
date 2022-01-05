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

func scanInput(scanner *bufio.Scanner, msgChannel, finalChannel chan struct{}) {
	for scanner.Scan() {
		fmt.Println("Receive a message")
		msgChannel <- struct{}{}
	}
	finalChannel <- struct{}{}
}

//!+
func handleConn(c net.Conn) {
	TCPConn := c.(*net.TCPConn)
	input := bufio.NewScanner(TCPConn)
	msgChannel := make(chan struct{})
	clientAbortChannel := make(chan struct{})
	go scanInput(input, msgChannel, clientAbortChannel)
loop:
	for {
		select {
		case <-msgChannel:
			waitGroup.Add(1)
			go echo(TCPConn, input.Text(), 1*time.Second)
		case <-clientAbortChannel:
			fmt.Println("The client have finished the conversation")
			break loop
		case <-time.After(time.Second * 10):
			fmt.Println("Receive no message in 10 seconds")
			break loop
		}
	}
	go func() {
		fmt.Println("Run waiting goroutine")
		waitGroup.Wait()
		TCPConn.CloseWrite()
	}()
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
