package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func runConveyerStage(inputChannel, outputChannel chan struct{}) {
	structure := <-inputChannel
	outputChannel<- structure
}

func runConveyer(goroutineNumber int) time.Duration {
	firstInputChannel := make(chan struct{})
	var inputChannel, outputChannel chan struct{}
	for i := 0; i < goroutineNumber; i++ {
		inputChannel = outputChannel
		if i == 0 {
			inputChannel = firstInputChannel
		}
		outputChannel = make(chan struct{})
		go runConveyerStage(inputChannel, outputChannel)
	}
	start := time.Now()
	firstInputChannel<- struct{}{}
	<-outputChannel
	return time.Since(start)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println(
			"ERROR. The script has one positional argument: the goroutine number",
		)
		return
	}
	inputArgument := os.Args[1]
	goroutineNumber, err := strconv.Atoi(inputArgument)
	if err != nil {
		fmt.Printf("ERROR. The argument '%s' is not an integer\n", inputArgument)
	}
	if goroutineNumber < 1 {
		fmt.Println("ERROR. The goroutine number is less than 1.")
	}
	workDuration := runConveyer(goroutineNumber)
	fmt.Printf(
		"It took %s to pass a conveyer of %d goroutines\n",
		workDuration,
		goroutineNumber,
	)
}
