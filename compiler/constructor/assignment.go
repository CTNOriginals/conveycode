package constructors

import (
	"conveycode/compiler/utils"
	"slices"
	"strings"
)

var operatorSymbols []string = []string{"+", "-", "*", "/"}

// Get the correct operator syntax from the operator symbol that was used
//
//	"+" = "add"
//	"-" = "sub"
func getOperator(operator string) string {
	switch operator {
	case "+":
		return "add"
	case "-":
		return "sub"
	case "*":
		return "mul"
	case "/":
		return "div"
	default:
		return operator
	}
}

func constructVariable(parts []string) []string {
	var outLine []string

	idx := slices.Index(parts, "=")
	outLine = append(outLine, "set", parts[idx-1])
	outLine = append(outLine, parts[idx+1:]...)

	return outLine
}
func constructOperation(parts []string) []string {
	var outLine []string
	outLine = append(outLine, "op")

	idx := slices.Index(parts, "=")
	variable := parts[:idx]
	operand := parts[idx+1:]

	//TODO Add support for multi operational assignments
	//* example: var z = x + y * (x - (y / x)) + y

	//! Currently only suppoerts single assignment
	outLine = append(outLine, getOperator(operand[1]))

	//? Append the variable name to the line
	outLine = append(outLine, variable[len(variable)-1])

	//? Append the 2 operational arguments
	outLine = append(outLine, operand[0], operand[2])

	return outLine
}

// Construct a variable assignment line
//
// Pass in the parts of the line which has been split by spaces
func Assignment(parts [][]rune) string {
	var stringParts []string

	for _, seg := range parts {
		stringParts = append(stringParts, string(seg))
	}

	if utils.ContainsListItem(stringParts, operatorSymbols) {
		return strings.Join(constructOperation(stringParts), " ")
	} else {
		return strings.Join(constructVariable(stringParts), " ")
	}
}
