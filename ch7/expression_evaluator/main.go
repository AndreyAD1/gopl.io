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

func getVariableList() ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	variableString, err := reader.ReadString('\n')
	if err == io.EOF {
		return nil, fmt.Errorf("no variable values")
	}
	if err != nil {
		return nil, err
	}
	variableString = strings.Replace(variableString, "\n", "", -1)
	variableList := strings.Split(variableString, ",")
	return variableList, nil
}

func getExpressionEnvironment(
	variableList []string,
	vars map[eval.Var]bool,
	varDescription string,
) (map[eval.Var]float64, error) {
	expressionEnvironment := make(map[eval.Var]float64)
	for _, variableString := range variableList {
		variableStringItems := strings.Split(variableString, "=")
		if len(variableStringItems) != 2 {
			return nil, fmt.Errorf("invalid input: %s", variableString)
		}
		variableName := variableStringItems[0]
		if _, ok := vars[eval.Var(variableName)]; !ok {
			err := fmt.Errorf("unexpected variable name: %s", variableString)
			return nil, err
		}
		if _, ok := expressionEnvironment[eval.Var(variableName)]; ok {
			return nil, fmt.Errorf("repeated variable name: %v", variableName)
		}
		variableValue, err := strconv.ParseFloat(variableStringItems[1], 64)
		if err != nil {
			err := fmt.Errorf(
				"unexpected variable value: %v",
				variableStringItems[1],
			)
			return nil, err
		}
		expressionEnvironment[eval.Var(variableName)] = variableValue
	}
	if len(expressionEnvironment) != len(vars) {
		err := fmt.Errorf(
			"not enough variable values: expect: %s",
			varDescription,
		)
		return nil, err
	}
	return expressionEnvironment, nil
}

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
		if i == len(vars)-1 {
			varDescription += string(v)
			break
		}
		varDescription += string(v) + ","
		i++
	}
	fmt.Println("Enter the variable values. Variable names: ", varDescription)
	fmt.Println("Expected input format: 'x=1;y=2;...'")
	variableList, err := getVariableList()
	if err != nil {
		fmt.Println(err)
		return
	}
	expressionEnvironment, err := getExpressionEnvironment(
		variableList,
		vars,
		varDescription,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	calculatedResult := expr.Eval(eval.Env(expressionEnvironment))
	fmt.Println("The expression result is: ", calculatedResult)
}
