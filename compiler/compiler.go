package compiler

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/parser"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/utils"
	"fmt"

	"github.com/TwiN/go-color"
)

func compile(tokens tokenizer.TokenList) []string {
	//#region Tokenizer
	fmt.Printf("\n-- %s --\n", color.InBlue("Tokenizer"))
	for _, token := range tokens {
		if token.Typ == tokenizer.EOL {
			fmt.Println("")
			continue
		}

		fmt.Print(color.InUnderline(token.ColoredValue()) + " ")
		// fmt.Printf("%s: %s\n", color.InGreen(token.Typ), token.ColoredValue())
	}
	//#endregion

	//#region Lexer
	fmt.Printf("\n\n-- %s --\n", color.InBlue("Lexer"))
	var lx = lexer.Lex(tokens)
	var blocks []lexer.Block = lx.Construct()

	for _, block := range blocks {
		fmt.Println(block)
	}
	//#endregion

	//#region Parser
	fmt.Printf("-- %s --\n", color.InBlue("Parser"))
	var prs = parser.Parse(blocks, parser.GlobalScope)

	var instructions = parser.Construct(prs)
	var instructionLines []string

	for _, instruction := range instructions {
		fmt.Println(instruction)
		instructionLines = append(instructionLines, instruction.String())
	}
	//#endregion

	return instructionLines
}

// TODO Allow for passing in a dir that doesnt point to a file, and then have all the files contained compiled
// Compile a .conv file to .mlog
//
//	compiler.CompileFile("foo/bar/file.conv", "dest/")
func CompileFile(sourceFilePath string, dest string) {
	fmt.Printf("\n-- %s %s --\n", color.InGreen("File"), color.InYellow(sourceFilePath))
	utils.WriteFile(dest, compile(tokenizer.Tokenize(sourceFilePath)))
}
