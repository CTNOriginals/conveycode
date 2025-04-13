package utils

import (
	"conveycode/constents"
	"fmt"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
)

type ParsedFilePath struct {
	Raw string
	// The full path split by '/'
	Split []string
	// The full path to the directory that contains the file
	Path string
	// The directory name that contains the file
	Dir string

	// The full file name + extension
	File string
	Name string
	Ext  string

	Line   int
	Column int

	// The full path excluding any trailing text like the line number
	Full string

	// Whatever was left at the trailing end of all of it
	Trail string
}

func (this ParsedFilePath) String() string {
	var lines []string
	var values = StructValues(this)
	for i, key := range StructKeys(this) {
		lines = append(lines, fmt.Sprintf("%s: %v", color.InBlue(key), values[i]))
	}

	return strings.Join(lines, "\n")
}

// Parses the file path into a file path struct
func ParseFilePath(path string) ParsedFilePath {
	path = strings.ReplaceAll(path, "\\", "/")

	var obj = ParsedFilePath{
		Raw:   path,
		Split: strings.Split(path, "/"),
	}

	//#region Path + Dir
	obj.Path = strings.Join(obj.Split[:len(obj.Split)-1], "/")
	if len(obj.Split)-2 >= 0 {
		obj.Dir = obj.Split[len(obj.Split)-2]
	}
	//#endregion

	//#region File
	var file = obj.Split[len(obj.Split)-1]
	rStart, rEnd := GetValidStringRange(file, constents.FileNameCharacters, 0)

	if rStart+rEnd == 0 {
		panic(fmt.Sprintf("Invalid file '%s'", file))
	}

	obj.File = file[rStart:rEnd]
	obj.Trail = file[rEnd:]

	//#region Name + Ext
	var fileSplit = strings.Split(obj.File, ".")

	if len(fileSplit) != 1 {
		obj.Name = strings.Join(fileSplit[:len(fileSplit)-1], ".")

		var trail = fileSplit[len(fileSplit)-1]
		rStart, rEnd := GetValidStringRange(trail, constents.AlphaNumaric, 0)

		if rStart+rEnd > 0 {
			obj.Ext = trail[rStart:rEnd]
		}
	} else {
		obj.Name = fileSplit[0]
	}
	//#endregion
	//#endregion

	//#region Line + Colmn
	//- Make sure there is nothing between the extension and the colon before the line number
	if len(obj.Trail) > 0 && rune(obj.Trail[0]) == ':' {
		rStart, rEnd := GetValidStringRange(obj.Trail, constents.Numbers+":", 1)
		if rStart+rEnd > 0 {
			var loc = obj.Trail[rStart:rEnd]
			var locSplit = strings.Split(loc, ":")
			if len(locSplit) >= 1 {
				obj.Line, _ = strconv.Atoi(locSplit[0])
				if len(locSplit) >= 2 {
					obj.Column, _ = strconv.Atoi(locSplit[1])
				}
			}

			obj.Trail = obj.Trail[rEnd:]
		}
	}
	//#endregion

	if len(obj.Path) > 0 {
		obj.Full = fmt.Sprintf("%s/%s", strings.TrimLeft(obj.Path, " \t\n"), obj.File)
	} else {
		obj.Full = obj.File
	}

	return obj
}
