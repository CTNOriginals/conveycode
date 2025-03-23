package compiler

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/utils"
	"fmt"
	"strconv"

	"github.com/TwiN/go-color"
)

// Compile a .conv file to .mlog
//
//	compiler.CompileFile("foo/bar/file.conv", "dest/")
func CompileFile(sourceFilePath string, dest string) {
	fmt.Printf("File %s\n", color.InYellow(sourceFilePath))

	// tools.CursorTests(utils.GetFileRunes(sourceFilePath))

	var content = utils.GetFileRunes(sourceFilePath)
	var tokens tokenizer.TokenList = tokenizer.Tokenize(content)

	//? Debug logging
	fmt.Printf("-- %s --\n", color.InBlue("Tokenizer"))
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
			val, _ = strconv.ParseInt(string(token.Val), 0, 64)
		case tokenizer.Operator:
			col = color.Blue
		case tokenizer.Seperator:
			col = color.Cyan
		case tokenizer.RoundL, tokenizer.RoundR, tokenizer.SquareL, tokenizer.SquareR, tokenizer.CurlyL, tokenizer.CurlyR:
			col = color.Yellow

		}
		fmt.Print(color.InUnderline(color.Colorize(col, val)) + " ")
	}

	fmt.Printf("\n\n-- %s --\n", color.InBlue("Lexer"))
	var lx = lexer.Lex(content)
	var item, ok = <-lx.Items

	for ok {
		fmt.Println(item.ColoredString())
		item, ok = <-lx.Items
	}

	// for item.Typ != lexer.LexEOF {
	// 	// fmt.Println(string(item.Val))
	// }

	// utils.WriteFile(utils.GetFileName(sourceFilePath), dest, instructionLines)
}
