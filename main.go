package main

import (
	"fmt"
	"time"

	"github.com/TwiN/go-color"
)

func main() {
	fmt.Printf("\n---- Start %s ----\n", color.Colorize(color.Green, time.Now().Format(time.TimeOnly)))
}
