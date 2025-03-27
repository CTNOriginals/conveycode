package compiler

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/parser"
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
	var blocks []lexer.Block

	for lx.State != nil {
		var block = lx.NextBlock()
		blocks = append(blocks, block)
		fmt.Println(block)
	}

	var prs = parser.Parse(blocks)
	var instruction, ok = <-prs.Channel

	for ok {
		fmt.Println(instruction)
		instruction, ok = <-prs.Channel
	}

	// utils.WriteFile(utils.GetFileName(sourceFilePath), dest, instructionLines)
}
