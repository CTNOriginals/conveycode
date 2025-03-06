package compiler

import (
	"bufio"
	"conveycode/compiler/utils"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/TwiN/go-color"
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
		var operatorSymbols []string = []string{"+", "-", "*", "/"}

		if len(parts) == 0 {
			continue
		}

		var outLine []string

		if idx := slices.Index(parts, "="); idx != -1 {
			if utils.ContainsListItem(parts, operatorSymbols) {
				outLine = append(outLine, "op")
			} else {
				outLine = append(outLine, "set", parts[idx-1])
				outLine = append(outLine, parts[idx:]...)
			}
		}

		statements = append(statements, strings.Join(outLine, " "))
	}

	fmt.Println(strings.Join(statements, "\n"))
}
