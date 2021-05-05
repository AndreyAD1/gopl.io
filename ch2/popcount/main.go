package main

import "fmt"

var byteArray [256]byte

func main() {
	fmt.Print(byteArray)
	
	for index, item := range byteArray {
		fmt.Printf("%v,%v\n", item)
		fmt.Printf("%v\n", byteArray[index/2])
		byteArray[index] = byteArray[index/2] + byte(index&1)
	}
	fmt.Print(byteArray)
}
