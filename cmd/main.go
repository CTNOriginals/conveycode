package main

import (
	"conveycode/compiler"
	"fmt"
	"time"

	"github.com/TwiN/go-color"
)

var testCases [][]string = [][]string{
	// {"tests/assignment/setAdd.conv", "tests/assignment/compiled/"},
	// {"tests/print/print.conv", "tests/print/compiled/"},
	// {"tests/print/printLine.conv", "tests/print/compiled/"},
	// {"tests/print/printInterpelate.conv", "tests/print/compiled/"},
	// {"tests/condition/ifStatement.conv", "tests/condition/compiled/"},
	{"tests/prototype/proto.conv", "tests/prototype/compiled/"},
}

func main() {
	fmt.Printf("\n---- Start %s ----\n", color.Colorize(color.Green, time.Now().Format(time.TimeOnly)))

	for _, testCase := range testCases {
		compiler.CompileFile(testCase[0], testCase[1])
	}
}
