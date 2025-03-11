package compiler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"

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

	var debugStatement [][][]rune

	var statement []string

	for _, line := range fileLines(sourceFilePath) {
		//#region Sudo Code
		//? The first objective of this scanner is to subdivide the line into relevant segments
		//? Each segment represents a portion that has meaning by its self, for example:
		//* the first word/character of the line
		//> var
		//> print
		//> if
		//> etc
		//+ this can be detected in a few different ways
		//- 1. stop at the first whitespace character and push the segment.
		//! this would not work with functions like print()
		//- 2. Have a predifined const list of keywords to look for, and once the current segment of characters matches any of the,
		//< my preference goes to this
		//! This would require constent comparrison and might need some special edge cases
		//> check if there is a longer variation of the keyword,
		//> like print also has println, so if print is matched, check if ln directly follows to match that instead
		//- 3. Dont do anything special
		//< This would be the most hands off solution for this loop
		//! might dump too much on the next loop
		//> just create segments on whitespaces and other points
		//> let the next loop handle the keyword checking instead
		//* an enclosure
		//? In this case, the enclosure would be captured with the symbols on the outer end
		//? and everything inside it without breaking it into more segments
		//> a quoted string
		//+ "", '', ``
		//> any open and close braces on the same line
		//+ (), [], {}
		//- if the enclosure is not close on the same line,
		//- it needs to be handled later on by a more advanced function that can parse scoped blocks
		//? once the whole line is scanned, clean up some unresolved groups, like unmet enclosures
		//? after that, it should be appended to the statements array
		//? when this for loop is done, it would have to be given to another for loop after this one to make sense of what all of it is
		//#endregion

		var segments [][]rune

		var current []rune

		quoteScope := &QuotedScope{
			symbol: "",
			state:  false,
		}

		line = strings.TrimSpace(line)

		for i, char := range line {
			//? If we encounter a space and we're not inside a quote, finalize the current segment.
			if !quoteScope.state && unicode.IsSpace(char) {
				if current != nil {
					segments = append(segments, current)
					current = nil
				}
				continue
			}

			//? Handle quotes (either starting or ending a quoted string).
			if regQuotes.MatchString(string(char)) && (i == 0 || line[i-1] != '\\') {
				//? If we are inside a quote and we encounter the same quote symbol, finalize the current string.
				if quoteScope.state && string(char) == quoteScope.symbol {
					current = append(current, char)
					segments = append(segments, current)
					current = nil
					quoteScope.state = false
					continue
				}

				//? If we are outside a quote and have accumulated characters, finalize the current segment.
				if !quoteScope.state && current != nil {
					segments = append(segments, current)
					current = nil
				}

				//? Toggle quote state and set the symbol to the current quote type.
				quoteScope.symbol = string(char)
				quoteScope.state = !quoteScope.state
			}

			//? Add the current character to the current segment.
			current = append(current, char)

			//? If it's the last character of the line, finalize the current segment.
			if i == len(line)-1 && current != nil {
				segments = append(segments, current)
				current = nil
			}
		}

		if len(segments) == 0 {
			continue
		}

		statement = append(statement, ParseSegments(segments))
		debugStatement = append(debugStatement, segments)
	}

	//? Debug Logging
	// for _, line := range debugStatement {
	// 	for j, seg := range line {
	// 		fmt.Printf("%d: %q\n", j, seg)
	// 	}
	// 	fmt.Println("")
	// }

	writeFile(getFileName(sourceFilePath), dest, statement)
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
