package utils

import (
	"bytes"
	"conveycode/internal"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/TwiN/go-color"
)

var showPointers = false

func PrintStackTrace() {
	var seperator = []byte("\n")
	var stack = debug.Stack()
	//? Remove the fers few lines that contain the traces that lead here that
	var stackLines = bytes.Split(stack, seperator)[9:]
	var rootPath = strings.ReplaceAll(internal.RootPath, "\\", "/")

	var lines = make([]string, len(stackLines))

	for i, byteLine := range stackLines {
		if len(byteLine) == 0 {
			continue
		}
		var line = strings.ReplaceAll(string(byteLine), rootPath, internal.ProjectName)
		var filePath = ParseFilePath(line)

		if strings.HasPrefix(line, "\t") {
			line = fmt.Sprintf(
				"%s/%s.%s:%s",
				color.InBlue(filePath.Path),
				color.InGreen(filePath.Name),
				color.InCyan(filePath.Ext),
				color.InYellow(filePath.LineNumber),
			)
		} else {
			line = fmt.Sprintf(
				"\x1b[2m%s\x1b[2m/%s\x1b[2m%s\x1b[0m",
				color.InBlue(filePath.Path),
				color.InCyan(filePath.File),
				If(showPointers, filePath.Trail, ""),
			)
		}

		lines[i] = line
	}

	fmt.Println(strings.ReplaceAll(strings.Join(lines, string(seperator)), rootPath, internal.ProjectName))
}
