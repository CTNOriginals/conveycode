package compiler

import (
	"conveycode/compiler/types"
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

	/*instructions*/
	var tokens types.TokenList = Tokenize(utils.GetFileRunes(sourceFilePath))

	//? Debug logging
	for _, token := range tokens {
		if token.Typ == types.EOL {
			fmt.Println("")
			continue
		}

		var col string = color.Gray
		var val any = string(token.Val)

		switch token.Typ {
		case types.String:
			col = color.Red
		case types.Number:
			col = color.Green
			val, _ = strconv.ParseInt(string(token.Val), 0, 64)
		case types.Operator:
			col = color.Blue
		case types.Seperator:
			col = color.Cyan
		case types.RoundL, types.RoundR, types.SquareL, types.SquareR, types.CurlyL, types.CurlyR:
			col = color.Yellow

		}
		fmt.Print(color.InUnderline(color.Colorize(col, val)) + " ")
	}

	// utils.WriteFile(utils.GetFileName(sourceFilePath), dest, instructionLines)
}
