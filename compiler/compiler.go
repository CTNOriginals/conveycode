package compiler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/TwiN/go-color"
)

type QuotedScope struct {
	symbol string //? The character that started with it
	state  bool
}

var regQuotes *regexp.Regexp

func init() {
	var err error
	if regQuotes, err = regexp.Compile("[\"'`]"); err != nil {
		panic(err)
	}
}

// Compile a .conv file to .mlog
//
//	compiler.CompileFile("foo/bar/file.conv", "dest/")
func CompileFile(sourceFilePath string, dest string) {
	fmt.Printf("File %s\n", color.InGreen(sourceFilePath))

	var instructions [][][]rune = Tokenize(fileLines(sourceFilePath))
	var instructionLines []string

	for _, line := range instructions {
		//? Debug Logging
		for _, seg := range line {
			fmt.Printf("%s\n", string(seg))
		}
		fmt.Println("")

		instructionLines = append(instructionLines, ParseSegments(line))
	}

	writeFile(getFileName(sourceFilePath), dest, instructionLines)
}

// Parses the file path and returns just the file name without the extension
//
// Supports file names with any number of dots (.) in it
//
//	getFileName("foo/bar/fileName.ext") // fileName
//	getFileName("foo/bar/fileName.version.data.ext") // fileName.version.data
func getFileName(filePath string) string {
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	parts := strings.Split(filePath, "/")

	file := parts[len(parts)-1]
	split := strings.Split(file, ".")
	return strings.Join(split[:len(split)-1], ".")
}

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

func writeFile(fileName string, destPath string, lines []string) {
	fileChars := strings.Split(destPath, "")
	if !slices.Contains([]string{"/", "\\"}, fileChars[len(fileChars)-1]) {
		destPath += "/"
	}

	//? Make destination dir to make sure it exists
	_ = os.MkdirAll(destPath, 0666)

	file, err := os.Create(destPath + fileName + ".mlog")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			log.Fatal(err)
		}
	}
}
