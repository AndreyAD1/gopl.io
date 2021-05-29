package main

import "fmt"

func deleteDuplicates(slice []string) []string {
	for i := 0; i < len(slice); i++ {
		if i == 0 {
			continue
		}
		if slice[i-1] == slice[i] {
			fmt.Printf("Duplicated symbols %v and %v\n", slice[i-1], slice[i])
			if i < len(slice) - 1 {
				fmt.Printf("not last symbol: %v, index:%v\n", slice[i], i)
				copy(slice[i:], slice[i + 1:])
				i -= 1
			}
			slice = slice[:len(slice) - 1]
		}
		fmt.Println(slice)
	}
	return slice
}

func main() {
	duplicateSlice := []string{"a", "b", "b", "c", "c", "c", "d", "d", "d", "e"}
	deduplicatedSlice := deleteDuplicates(duplicateSlice)
	fmt.Println(deduplicatedSlice)
}