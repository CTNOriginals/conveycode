package utils

import (
	"conveycode/constents"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Parses the file path and returns just the file name without the extension
//
// Supports file names with any number of dots (.) in it
//
//	getFileName("foo/bar/fileName.ext") // fileName
//	getFileName("foo/bar/fileName.version.data.ext") // fileName.version.data
func GetFileName(path string) string {
	// filePath = strings.ReplaceAll(filePath, "\\", "/")
	// parts := strings.Split(filePath, "/")
	//
	// file := parts[len(parts)-1]
	// split := strings.Split(file, ".")
	// return strings.Join(split[:len(split)-1], ".")

	var _, name = filepath.Split(path)
	var split = strings.Split(name, ".")
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

// Calls fn for each file in dir.
// Does not recurse into child directories.
// Also returns directories.
func ForEachFileInDir(dir string, fn func(file os.FileInfo)) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
			continue
		}

		fn(info)
	}
}

// Calls fn for each file in dir recursivly
func ForEachFileInDirRecursive(dir string, fn func(file os.FileInfo, dir string)) {
	ForEachFileInDir(dir, func(file os.FileInfo) {
		if file.IsDir() {
			childDir := fmt.Sprintf("%s/%s", dir, file.Name())
			ForEachFileInDirRecursive(childDir, fn)
			return
		}

		fn(file, dir)
	})
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

// TODO: fix clutter and consistency by
// assuming path is the correct resulting file path.
func WriteFile(path string, lines []string) {
	// var filePath = ParseFilePath(path)
	var dir, name = filepath.Split(path)
	name = GetFileName(name)

	//? Make destination dir to make sure it exists
	_ = os.MkdirAll(dir, os.ModePerm)

	file, err := os.Create(fmt.Sprintf("%s/%s.mlog", dir, name))
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
