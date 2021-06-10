package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error. The script needs a file path as an input argument")
		return
	}
}