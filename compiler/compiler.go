package compiler

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/utils"
	"fmt"

	"github.com/TwiN/go-color"
)

// Compile a .conv file to .mlog
//
//	compiler.CompileFile("foo/bar/file.conv", "dest/")
func CompileFile(sourceFilePath string, dest string) {
	fmt.Printf("File %s\n", color.InYellow(sourceFilePath))

	var tokens tokenizer.TokenList = tokenizer.Tokenize(utils.GetFileRunes(sourceFilePath))

	//? Debug logging
	fmt.Printf("\n\n-- %s --\n", color.InBlue("Tokenizer"))
	for _, token := range tokens {
		if token.Typ == tokenizer.EOL {
			fmt.Println("")
			continue
		}

		fmt.Print(color.InUnderline(token.ColoredValue()) + " ")
	}

	fmt.Printf("\n\n-- %s --\n", color.InBlue("Lexer"))
	var lx = lexer.Lex(tokens)

	for lx.State != nil {
		fmt.Print(lx.NextBlock())
	}

	// utils.WriteFile(utils.GetFileName(sourceFilePath), dest, instructionLines)
}
