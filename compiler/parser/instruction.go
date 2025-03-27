package parser

import "strings"

type Instruction struct {
	Parts []string
}

func (this Instruction) String() string {
	return strings.Join(this.Parts, "")
}
