package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func runConveyer(goroutineNumber int) {

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
		fmt.Printf("The argument '%s' is not an integer\n", inputArgument)
	}
	start := time.Now()
	runConveyer(goroutineNumber)
	workDuration := time.Since(start)
	fmt.Printf(
		"It took %s to pass a conveyer of %d goroutines\n",
		workDuration,
		goroutineNumber,
	)
}
