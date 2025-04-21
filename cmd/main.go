package main

import (
	"conveycode/compiler"
	"conveycode/compiler/parser"
	"conveycode/constents"
	"fmt"
	"os"
	"time"

	"github.com/TwiN/go-color"
)

var VERSION = "v0.0.0"

// var testCases [][]string = [][]string{
// 	// {"tests/assignment/setAdd.conv", "tests/assignment/compiled/"},
// 	// {"tests/print/print.conv", "tests/print/compiled/"},
// 	// {"tests/print/printInterpelate.conv", "tests/print/compiled/"},
// 	// {"tests/condition/if.conv", "tests/condition/compiled/"},
// 	// {"tests/condition/ifElse.conv", "tests/condition/compiled/"},
// 	// {"tests/condition/elseIf.conv", "tests/condition/compiled/"},
// 	// {"tests/condition/conditions.conv", "tests/condition/compiled/"},
// 	{"tests/prototype/proto.conv", "tests/prototype/compiled/"},
// }

func main() {
	constents.VERSION = VERSION

	if constents.DEVELOPMENT {
		executeDevelopmentTests()
		return
	}

	Args = NewArguments(os.Args)
	parser.InitializeScope()
	compiler.CompileFile(Args.sourceFile, Args.destFile)
}

func executeDevelopmentTests() {
	fmt.Printf("\n\n---- Start %s ----\n", color.Colorize(color.Green, time.Now().Format(time.TimeOnly)))
	Args = NewArguments(mockArgs)
	fmt.Println(Args)

	parser.InitializeScope()
	compiler.CompileFile(Args.sourceFile, Args.destFile)

	// for _, testCase := range testCases {
	// 	parser.InitializeScope()
	// 	compiler.CompileFile(testCase[0], testCase[1])
	// }
}
