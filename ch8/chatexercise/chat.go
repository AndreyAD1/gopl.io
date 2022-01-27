// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

type ClientInfo struct {
	Channel client
	Name    string
}

var (
	entering = make(chan ClientInfo)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]string) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case clientInfo := <-entering:
			var connectedClients []string
			for _, name := range clients {
				connectedClients = append(connectedClients, name)
			}
			firstMsg := "Connected clients: " + strings.Join(connectedClients, ", ")
			clients[clientInfo.Channel] = clientInfo.Name
			clientInfo.Channel <- firstMsg

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientInput(
	scanner *bufio.Scanner,
	inputChannel chan<- string,
	exitChannel chan<- struct{},
) {
	for scanner.Scan() {
		inputChannel <- scanner.Text()
	}
	exitChannel <- struct{}{}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ClientInfo{ch, who}

	clientInputChannel := make(chan string)
	clientExitChannel := make(chan struct{})
	input := bufio.NewScanner(conn)
	go clientInput(input, clientInputChannel, clientExitChannel)
	for {
		connectionIsAlive := true
		select {
			case clientMessage := <-clientInputChannel:
				messages <- who + ": " + clientMessage
			case <-clientExitChannel:
				leaving <- ch
				messages <- who + " has left"
				connectionIsAlive = false
			case <-time.After(time.Minute * 5):
				leaving <- ch
				messages <- who + " kicked out for keeping silence"
				connectionIsAlive = false
		}
		if !connectionIsAlive {
			break
		}
	}
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
