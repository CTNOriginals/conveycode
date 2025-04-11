package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type parsedFilePath struct {
	Raw string
	// The full path split by '/'
	Split []string
	// The full path to the directory that contains the file
	Path string
	// The directory name that contains the file
	Dir string

	// The full file name + extension
	File       string
	Name       string
	Ext        string
	LineNumber int

	// The full path excluding any trailing text like the line number
	Full string

	// Whatever was left at the trailing end of all of it
	Trail string
}

func (this parsedFilePath) String() string {
	var lines []string
	var values = StructValues(this)
	for i, key := range StructKeys(this) {
		lines = append(lines, fmt.Sprintf("%s: %v", key, values[i]))
	}

	return strings.Join(lines, "\n")
}

// Parses the file path into a file path struct
func ParseFilePath(path string) parsedFilePath {
	path = strings.ReplaceAll(path, "\\", "/")

	var obj = parsedFilePath{
		Raw:        path,
		Split:      strings.Split(path, "/"),
		LineNumber: -1,
	}

	obj.Path = strings.Join(obj.Split[:len(obj.Split)-1], "/")
	if len(obj.Split)-2 >= 0 {
		obj.Dir = obj.Split[len(obj.Split)-2]
	}

	obj.File = obj.Split[len(obj.Split)-1]

	var fileSplit = strings.Split(obj.File, ".")
	obj.Name = strings.Join(fileSplit[:len(fileSplit)-1], ".")

	var extSplit = strings.Split(fileSplit[len(fileSplit)-1], "")
	if len(extSplit) > 0 {
		var streamLoc = WordCharacterStream.FindStringIndex(strings.Join(extSplit, ""))

		var extBreak = streamLoc[If(len(streamLoc) > 1, 1, 0)]
		obj.Ext = strings.Join(extSplit[streamLoc[0]:streamLoc[1]], "")
		obj.Trail = strings.Join(extSplit[extBreak:], "")

		obj.File = fmt.Sprintf("%s.%s", obj.Name, obj.Ext)

		if len(extSplit) > streamLoc[1] {
			var numEnd = NonNumberCharacter.FindStringIndex(strings.Join(extSplit[extBreak+1:], ""))
			obj.LineNumber, _ = strconv.Atoi(If(
				extSplit[extBreak] == ":",
				strings.Join(extSplit[extBreak+1:(extBreak+1)+numEnd[0]], ""),
				"-1",
			))
		}
	}

	obj.Full = fmt.Sprintf("%s/%s", strings.TrimLeft(obj.Path, " \t\n"), obj.File)

	return obj
}

// Parses the file path and returns just the file name without the extension
//
// Supports file names with any number of dots (.) in it
//
//	getFileName("foo/bar/fileName.ext") // fileName
//	getFileName("foo/bar/fileName.version.data.ext") // fileName.version.data
func GetFileName(filePath string) string {
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	parts := strings.Split(filePath, "/")

	file := parts[len(parts)-1]
	split := strings.Split(file, ".")
	return strings.Join(split[:len(split)-1], ".")
}

func GetFileLines(filePath string) []string {
	file := getFile(filePath)
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var lines []string

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines
}

func GetFileRunes(filePath string) []rune {
	b, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	ret := make([]rune, len(b))

	for i, r := range b {
		ret[i] = rune(r)
	}

	return ret
}

func WriteFile(fileName string, destPath string, lines []string) {
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

func getFile(filePath string) *os.File {
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}
	return readFile
}
