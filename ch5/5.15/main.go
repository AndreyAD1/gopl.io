package main

import (
	"fmt"
	"math"
)

func max(values ...float64) float64 {
	maxValue := math.Inf(-1)
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func min(values ...float64) float64 {
	minValue := math.Inf(1)
	for _, value := range values {
		if value < minValue {
			minValue = value
		}
	}
	return minValue
}

func main() {
	fmt.Println(max(1.56, 9.1, -10.0))
	fmt.Println(max(1.56))
	fmt.Println(max())

	fmt.Println(min(1.56, 9.1, -10.0))
	fmt.Println(min(1.56))
	fmt.Println(min())
}