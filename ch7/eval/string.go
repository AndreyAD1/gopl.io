package eval

import (
	"fmt"
	"strconv"
)

func (v Var) String() string {
	return string(v)
}

func (l literal) String() string {
	return strconv.FormatFloat(float64(l), 'f', -1, 64)
}

func (u unary) String() string {
	return fmt.Sprint(string(u.op), u.x)
}

func (b binary) String() string {
	return fmt.Sprintf("%s %s %s", b.x, string(b.op), b.y)
}

func (c call) String() string {
	argumentString := ""
	for i, arg := range c.args {
		if i < len(c.args) - 1 {
			argumentString += fmt.Sprint(arg, ",")
			continue
		}
		argumentString += fmt.Sprint(arg)
	}
	return fmt.Sprintf("%s(%s)", c.fn, argumentString)
}