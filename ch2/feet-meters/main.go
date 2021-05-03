package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Foot float64
type Meter float64

func FeetToMeter(feet Foot) Meter {
	return Meter(feet * 0.3048)
}

func MeterToFeet(meter Meter) Foot {
	return Foot(meter / 0.3048)
}

func printResult(argument string) {
	number, err := strconv.ParseFloat(argument, 64)
	if err != nil {
		fmt.Printf(
			"Error. Invalid argument \"%s\". Expect number argument only\n",
			argument)
		return
	}
	meters := FeetToMeter(Foot(number))
	feet := MeterToFeet(Meter(number))
	fmt.Printf("%g feet = %g meters\n", number, meters)
	fmt.Printf("%g meters = %g feet\n", number, feet)
}

func main() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			printResult(arg)
		}
	} else {
		stdin_scanner := bufio.NewScanner(os.Stdin)
		for {
			stdin_scanner.Scan()
			input_argument := stdin_scanner.Text()
			if len(input_argument) == 0 {
				break
			}
			printResult(input_argument)
		}
	}
}
