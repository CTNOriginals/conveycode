package utils

import (
	"conveycode/constents"
	"fmt"
	"log"
	"os"
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

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsValidFileName(name string) bool {
	return ValidateString(name, constents.FileNameCharacters)
}
func IsValidDirectoryPath(path string) bool {
	return ValidateString(path, constents.DirectoryCharacters)
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

func WriteFile(path string, lines []string) {
	var filePath = ParseFilePath(path)

	//? Make destination dir to make sure it exists
	_ = os.MkdirAll(filePath.Path, 0666)

	file, err := os.Create(fmt.Sprintf("%s/%s.mlog", filePath.Path, filePath.Name))
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
