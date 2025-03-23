package compiler

import (
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

	// tools.CursorTests(utils.GetFileRunes(sourceFilePath))

	/*instructions*/
	var tokens tokenizer.TokenList = tokenizer.Tokenize(utils.GetFileRunes(sourceFilePath))

	//? Debug logging
	for _, token := range tokens {
		if token.Typ == tokenizer.EOL {
			fmt.Println("")
			continue
		}

		var col string = color.Gray
		var val any = string(token.Val)

		switch token.Typ {
		case tokenizer.String:
			col = color.Red
		case tokenizer.Number:
			col = color.Green
		case tokenizer.Operator:
			col = color.Blue
		case tokenizer.Seperator:
			col = color.Cyan
		case tokenizer.RoundL, tokenizer.RoundR, tokenizer.SquareL, tokenizer.SquareR, tokenizer.CurlyL, tokenizer.CurlyR:
			col = color.Yellow
		}
		fmt.Print(color.InUnderline(color.Colorize(col, val)) + " ")
	}

	// utils.WriteFile(utils.GetFileName(sourceFilePath), dest, instructionLines)
}
