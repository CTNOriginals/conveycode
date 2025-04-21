package main

import (
	"conveycode/compiler/utils"
	"conveycode/constents"
	"fmt"
	"slices"
	"strings"
)

var argumentSyntax = `
ConveyCode command-line arguments:
<source-file> [dest-dir] [--name <file-name>]

source-file:
	A file with the '.conv' extension to be compiled.
dest-dir:
	The destination directory for the compiled .mlog file to go in.
	If ommited, it will be created in the directory of the <source-file>
	in a directory called 'compiled'
--name
	The name of the resulting compiled file

OPTIONS:
--development, -D
	Turns development mode on
--mock
	Replaces all current arguments with the
	mock arguments in ./cmd/arguments.go
--log, -L
	Logs the tokenizer, lexer and parser results
`

func PrintArgumentSyntax() {
	fmt.Println(argumentSyntax)
}

var mockArgs = []string{"script/path/main.go",
	"tests/prototype/proto.conv",
	"--help",
	"--development",
	// "-L",
}

type Arguments struct {
	sourceFile string
	destFile   string
}

type argOption struct {
	flags  []string
	action func()
}

var options = []argOption{
	{flags: []string{"help", "--help", "-h", "-H", "?", "-?"},
		action: func() {
			PrintArgumentSyntax()
		},
	},
	{flags: []string{"--development", "-D"},
		action: func() {
			constents.DEVELOPMENT = true
		},
	},
	{flags: []string{"--log", "-L"},
		action: func() {
			constents.LOGGING = true
		},
	},
}

func NewArguments(input []string) (args Arguments) {
	input = input[1:]
	fmt.Println(strings.Join(input, "\n"))

	if len(input) == 0 {
		PrintArgumentSyntax()
		return
	}

	if containsArg(input, []string{"--mock"}) {
		return NewArguments(mockArgs)
	}

	for _, opt := range options {
		if containsArg(input, opt.flags) {
			opt.action()
		}
	}

	var sourcePath = utils.ParseFilePath(input[0])
	if !utils.FileExists(sourcePath.Full) {
		fmt.Printf("Source File does not exist '%s'\n", sourcePath.Full)
		return
	}

	args.sourceFile = sourcePath.Full

	var fileName = sourcePath.Name

	if index := slices.Index(input, "--name"); index != -1 {
		if index == len(input)-1 {
			fmt.Println("Argument flag '--name' requires the file name to follow it")
			return
		}

		fileName = input[index+1]
		if !utils.IsValidFileName(fileName) {
			fmt.Printf("Invalid file name '%s'\n", fileName)
			return
		}

		input = slices.Delete(input, index, index+2)
	}

	if len(input) == 1 || !utils.IsValidDirectoryPath(input[1]) || input[1] == "" {
		args.destFile = fmt.Sprintf("%s/compiled/%s.%s", sourcePath.Path, fileName, constents.CompiledFileExtention)
	} else {
		var dir = input[1]
		if !strings.HasSuffix(dir, constents.CompiledFileExtention) {
			if !strings.HasSuffix(dir, "/") && !strings.HasSuffix(dir, "\\") {
				dir += "/"
			}
			dir += fmt.Sprintf("%s.%s", fileName, constents.CompiledFileExtention)
		}
		args.destFile = dir
	}

	return args
}

func containsArg(input []string, search []string) (contained bool) {
	if utils.ContainsListItem(input, search) {
		for i, arg := range input {
			if slices.Contains(search, arg) {
				input = slices.Delete(input, i, i+1)
			}
		}

		return true
	}

	return false
}
