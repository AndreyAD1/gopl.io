package main

import (
	"fmt"
	"strings"
)


func getUpperCase(input string) string {
	return strings.ToUpper(input)
}

func expand(input string, convertationFunc func(string) string) string {
	stringParts := strings.Split(input, "$")
	var convertedStrings []string
	convertedStrings = append(convertedStrings, stringParts[0])
	for _, stringPart := range stringParts[1:] {
		convertedString := convertationFunc(stringPart)
		convertedStrings = append(convertedStrings, convertedString)
	}
	return strings.Join(convertedStrings, "")
}

func main() {
	input := "Some $string inpu$t"
	convertedString := expand(input, getUpperCase)
	fmt.Println(convertedString)
}
