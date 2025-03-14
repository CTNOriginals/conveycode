package compiler

import (
	"conveycode/src/compiler/constructor"
	"slices"
	"strings"
)

// import "fmt"

var variables []string

func ParseSegments(tokens []string) string {
	var parts []string

	for _, seg := range tokens {
		parts = append(parts, string(seg))
	}

	if string(tokens[0]) == "var" {
		if !slices.Contains(variables, string(tokens[1])) {
			variables = append(variables, string(tokens[1]))
		}
		return constructor.Assignment(parts)
	} else if slices.Contains(variables, parts[0]) {
		return constructor.Assignment(parts)
	}

	return ("# ERROR: " + strings.Join(parts, ""))

	// for _, seg := range segments {
	// }
}
