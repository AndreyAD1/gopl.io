package main

import "fmt"

import "gopl.io/ch7/eval"

func main() {
	inputExpression := "a * (b + c)"
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
		if i < len(vars) - 1 {
			varDescription += string(v) + ","
			i++
		}
		varDescription += string(v)
	}
	fmt.Println("Enter the variable values. Variable names: ", varDescription)
	fmt.Println("Expected input format: 'x=1;y=2...'")
	variableValues := map[eval.Var]float64 {"a": 2, "b": 1, "c": 3}
	calculatedResult := expr.Eval(eval.Env(variableValues))
	fmt.Println("The expression result is: ", calculatedResult)
}