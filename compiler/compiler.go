package compiler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TwiN/go-color"

	constructors "conveycode/compiler/constructor"
)

func fileLines(filePath string) []string {
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var lines []string

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines
}

// Compile a .conv file to .mlog
//
//	compiler.CompileFile("foo/bar/file.conv", "dest/")
func CompileFile(path string, dest string) {
	fmt.Printf("File %s\n", color.InGreen(path))

	var statements []string

	for _, line := range fileLines(path) {
		var parts []string = strings.Split(line, " ")

		//? Skip empty lines
		if len(parts) == 1 && parts[0] == "" {
			continue
		}

		var outLine []string

		switch parts[0] {
		case "var", "set":
			outLine = append(outLine, constructors.Assignment(parts))
		}

		statements = append(statements, strings.Join(outLine, " "))
	}

	fmt.Println(strings.Join(statements, "\n"))
}
