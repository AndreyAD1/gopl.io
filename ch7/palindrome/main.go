package main

import (
	"fmt"
	"sort"
)

type PalindromeChecker string

func (p PalindromeChecker) Len() int {
	return len(p)
}

func (p PalindromeChecker) Swap(i, j int) {}

func (p PalindromeChecker) Less(i, j int) bool {
	return p[i] < p[j]
}

func IsPalindrome(s sort.Interface) bool {
	length := s.Len()
	for i := 0; i < length - 1 - i; i++{
		if !s.Less(i, length - 1 - i) && !s.Less(length - 1 - i, i) {
			continue
		}
		return false
	}
	return true
}

func main() {
	test := PalindromeChecker("not palindrome")
	fmt.Println(IsPalindrome(test))
	test2 := PalindromeChecker("popop")
	fmt.Println(IsPalindrome(test2))
	test3 := PalindromeChecker("rats live on no evil star")
	fmt.Println(IsPalindrome(test3))
}