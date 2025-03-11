package compiler

import (
	constructors "conveycode/compiler/constructor"
	"slices"
)

// import "fmt"

var variables []string

func ParseSegments(segments [][]rune) string {
	if string(segments[0]) == "var" {
		if !slices.Contains(variables, string(segments[1])) {
			variables = append(variables, string(segments[1]))
		}
		return constructors.Assignment(segments)
	} else if slices.Contains(variables, string(segments[0])) {
		return constructors.Assignment(segments)
	}

	return "unknown"

	// for _, seg := range segments {
	// }
}
