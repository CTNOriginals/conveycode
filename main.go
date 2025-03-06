package main

import (
	"fmt"
	"time"

	"github.com/TwiN/go-color"

	"conveycode/compiler"
)

var testCases [][]string = [][]string{
	{"tests/assignment/setAdd.conv", "tests/assignment/compiled"},
}

func main() {
	fmt.Printf("\n---- Start %s ----\n", color.Colorize(color.Green, time.Now().Format(time.TimeOnly)))

	for _, testCase := range testCases {
		compiler.CompileFile(testCase[0], testCase[1])
	}
}
