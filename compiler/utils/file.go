package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

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
	file := getFile(filePath)
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanRunes)

	var runes []rune

	for fileScanner.Scan() {
		tokens := fileScanner.Bytes()
		runes = append(runes, rune(tokens[0]))
	}

	return runes
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

// func getFileContent[T bufio.SplitFunc, R []string | []rune | []byte](filePath string, scanType T) R {
// 	file := getFile(filePath)
// 	defer file.Close()

// 	fileScanner := bufio.NewScanner(file)
// 	fileScanner.Split(scanType)

// 	var content []R

// 	for fileScanner.Scan() {
// 		content = append(content, fileScanner.Text())
// 	}

// 	return content
// }
