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
	fmt.Printf("File %s\n", color.InYellow(sourceFilePath))

	// tools.CursorTests(utils.GetFileRunes(sourceFilePath))

	var instructions types.TokenList = Tokenize(utils.GetFileRunes(sourceFilePath))
	var instructionLines []string

	for _, content := range instructions.Tokens {
		//? Debug Logging
		fmt.Printf("\n%s", content)
	}

	utils.WriteFile(utils.GetFileName(sourceFilePath), dest, instructionLines)
}
