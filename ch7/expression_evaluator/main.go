package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopl.io/ch7/eval"
)

func main() {
	if len(os.Args) != 2 {
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
	reader := bufio.NewReader(os.Stdin)
	expressionEnvironment := make(map[eval.Var]float64)
	variableString, err := reader.ReadString('\n')
	if err == io.EOF {
		fmt.Println("No variable values")
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	variableString = strings.Replace(variableString, "\n", "", -1)
	variableList := strings.Split(variableString, ",")
	for _, variableString := range variableList {
		variableStringItems := strings.Split(variableString, "=")
		if len(variableStringItems) != 2 {
			fmt.Println("Invalid input: ", variableString)
			return
		}
		variableName := variableStringItems[0]
		if _, ok := vars[eval.Var(variableName)]; !ok {
			fmt.Println("Unexpected variable name: ", variableString)
			return
		}
		variableValue, err := strconv.ParseFloat(variableStringItems[1], 64)
		if err != nil {
			fmt.Println("Unexpected variable value: ", variableStringItems[1])
			return
		}
		expressionEnvironment[eval.Var(variableName)] = variableValue
	}
	// variableValues := map[eval.Var]float64 {"a": 2, "b": 1, "c": 3}
	calculatedResult := expr.Eval(eval.Env(expressionEnvironment))
	fmt.Println("The expression result is: ", calculatedResult)
}