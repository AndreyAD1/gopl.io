package main

import (
	"fmt"
	"os"
	"strconv"
)

func rotate(slice []int, offset int) []int {
	rotatedSlice := make([]int, len(slice))
	for index, item := range slice {
		rotatedIndex := index + offset
		if rotatedIndex >= len(slice) {
			rotatedIndex = rotatedIndex - len(slice)
		}
		rotatedSlice[rotatedIndex] = item
	}
	return rotatedSlice
}

func main () {
	if len(os.Args) != 2 {
		fmt.Println("Offset number is a required script argument")
		return
	}
	offset, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("An offset argument should be a number")
		return
	}
	slice := []int{1, 2, 3, 4, 5}
	rotatedSlice := rotate(slice, offset)
	fmt.Println(slice, rotatedSlice)
}