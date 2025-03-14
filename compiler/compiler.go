package compiler

import (
	"conveycode/compiler/types"
	"conveycode/compiler/utils"
	"fmt"
	"regexp"

	"github.com/TwiN/go-color"
)

type QuotedScope struct {
	symbol string //? The character that started with it
	state  bool
}

var regQuotes *regexp.Regexp

func init() {
	var err error
	if regQuotes, err = regexp.Compile("[\"'`]"); err != nil {
		panic(err)
	}
}

// Compile a .conv file to .mlog
//
//	compiler.CompileFile("foo/bar/file.conv", "dest/")
func CompileFile(sourceFilePath string, dest string) {
	fmt.Printf("File %s\n", color.InGreen(sourceFilePath))

	var instructions [][]types.Token = Tokenize(utils.GetFileLines(sourceFilePath))
	var instructionLines []string

	for _, line := range instructions {
		//? Debug Logging
		fmt.Println("")
		for _, seg := range line {
			fmt.Printf("%s\n", seg)
		}
		// instructionLines = append(instructionLines, ParseSegments(line))
	}

	utils.WriteFile(utils.GetFileName(sourceFilePath), dest, instructionLines)
}
