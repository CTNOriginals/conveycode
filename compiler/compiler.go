package compiler

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/parser"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/utils"
	"conveycode/constents"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TwiN/go-color"
)

func compile(tokens tokenizer.TokenList) []string {
	//#region Tokenizer
	if constents.LOGGING {
		fmt.Printf("\n-- %s --\n", color.InBlue("Tokenizer"))
		for _, token := range tokens {
			if token.Typ == tokenizer.EOL {
				fmt.Println("")
				continue
			}

			fmt.Print(color.InUnderline(token.ColoredValue()) + " ")
			// fmt.Printf("%s: %s\n", color.InGreen(token.Typ), token.ColoredValue())
		}
	}
	//#endregion

	//#region Lexer
	if constents.LOGGING {
		fmt.Printf("\n\n-- %s --\n", color.InBlue("Lexer"))
	}

	var lx = lexer.Lex(tokens)
	var blocks []lexer.Block = lx.Construct()

	if constents.LOGGING {
		for _, block := range blocks {
			fmt.Println(block)
		}
	}
	//#endregion

	//#region Parser
	if constents.LOGGING {
		fmt.Printf("-- %s --\n", color.InBlue("Parser"))
	}

	var prs = parser.Parse(blocks, parser.GlobalScope)

	var instructions = parser.Construct(prs)
	var instructionLines []string

	for _, instruction := range instructions {
		if constents.LOGGING {
			fmt.Println(instruction)
		}

		instructionLines = append(instructionLines, instruction.String())
	}
	//#endregion

	return instructionLines
}

// TODO Allow for passing in a dir that doesnt point to a file, and then have all the files contained compiled
// Compile a .conv file to .mlog
//
//	compiler.CompileFile("foo/bar/file.conv", "dest/")
//
// TODO: fix dest functionality, currently not implemented
func CompileFile(source string, dest string) {
	var info, err = os.Stat(source)
	if err != nil {
		panic(err)
	}

	if info.IsDir() {
		CompileDir(source)
		return
	}

	dest = fmt.Sprintf("%s/compiled/%s.mlog", filepath.Dir(source), utils.GetFileName(info.Name()))

	parser.InitializeConstructor()
	parser.InitializeScope()

	fmt.Printf("\n-- %s %s --\n", color.InGreen("File"), color.InYellow(source))
	utils.WriteFile(dest, compile(tokenizer.Tokenize(source)))
}

func CompileDir(source string) {
	utils.ForEachFileInDirRecursive(source, func(file os.FileInfo, dir string) {
		var name = file.Name()

		if filepath.Ext(name) != ".conv" {
			return
		}

		var s = fmt.Sprintf("%s%s", dir, name)
		var d = fmt.Sprintf("%scompiled/%s.mlog", dir, utils.GetFileName(name))
		// fmt.Printf("Compiling: %s > %s\n", s, d)
		CompileFile(s, d)
	})
}
