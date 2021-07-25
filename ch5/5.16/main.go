package main

import "fmt"

func Join(separator string, strs ...string) string {
	var outputStr string
	for index, str := range strs {
		outputStr += str
		if index != len(strs) - 1 {
			outputStr += separator
		}
	}
	return outputStr
}

func main() {
	fmt.Println(Join("-", "first", "middle", "last"))
	fmt.Println(Join("-", )) // ""
	fmt.Println(Join("-", "last")) // "last"
}