package main

import (
	"fmt"
	"os"

	"gopl.io/ch7/eval"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("An expression is required script argument")
		fmt.Println("Expression example: 'a * (b + c)'")
		return
	}
	inputExpression := os.Args[1]
	expr, err := eval.Parse(inputExpression)
	if err != nil {
		fmt.Println(err)
		return
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		fmt.Println(err)
		return
	}
	varDescription := "" 
	i := 0
	for v := range vars {
		if i == len(vars) - 1 {
			varDescription += string(v)
			break
		}
		varDescription += string(v) + ","
		i++
	}
	fmt.Println("Enter the variable values. Variable names: ", varDescription)
	fmt.Println("Expected input format: 'x=1;y=2;...'")
	variableValues := map[eval.Var]float64 {"a": 2, "b": 1, "c": 3}
	calculatedResult := expr.Eval(eval.Env(variableValues))
	fmt.Println("The expression result is: ", calculatedResult)
}