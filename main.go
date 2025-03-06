package main

import (
	"fmt"
	"time"

	"github.com/TwiN/go-color"

	"conveycode/compiler"
)

var filePath string = "tests/code/setAdd.conv"
var destDir string = "tests/compiled/"

func main() {
	fmt.Printf("\n---- Start %s ----\n", color.Colorize(color.Green, time.Now().Format(time.TimeOnly)))

	compiler.CompileFile(filePath, destDir)
}
