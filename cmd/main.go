package main

import (
	"conveycode/compiler"
	"conveycode/compiler/parser"
	"fmt"
	"time"

	"github.com/TwiN/go-color"
)

var testCases [][]string = [][]string{
	// {"tests/assignment/setAdd.conv", "tests/assignment/compiled/"},
	// {"tests/print/print.conv", "tests/print/compiled/"},
	// {"tests/print/printInterpelate.conv", "tests/print/compiled/"},
	// {"tests/condition/if.conv", "tests/condition/compiled/"},
	// {"tests/condition/ifElse.conv", "tests/condition/compiled/"},
	// {"tests/condition/elseIf.conv", "tests/condition/compiled/"},
	// {"tests/condition/conditions.conv", "tests/condition/compiled/"},
	{"tests/prototype/proto.conv", "tests/prototype/compiled/"},
}

func main() {
	fmt.Printf("\n\n---- Start %s ----\n", color.Colorize(color.Green, time.Now().Format(time.TimeOnly)))

	// debug.SetMaxStack(64)
	// debug.SetMaxThreads(1)

	for _, testCase := range testCases {
		parser.InitializeScope()
		compiler.CompileFile(testCase[0], testCase[1])
	}

	// fmt.Printf("\n\n---- End %s ----\n", color.Colorize(color.Green, time.Now().Format(time.TimeOnly)))
}
