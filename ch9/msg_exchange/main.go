package main

import (
	"fmt"
	"time"
)

var (
	counter int
	channel = make(chan struct{})
)

func exchangeMsg() {
	for {
		msg := <-channel
		counter++
		channel<- msg
	}
}

func main() {
	go exchangeMsg()
	go exchangeMsg()
	channel<- struct{}{}
	time.Sleep(time.Second * 2)
	fmt.Printf("%v messages passed the channel per one second\n", counter)
}
