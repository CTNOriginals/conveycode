package parser

import "strings"

type Instruction struct {
	Parts []string
}

func NewInstruction(parts []string) Instruction {
	return Instruction{
		Parts: parts,
	}
}

func (this Instruction) String() string {
	return strings.Join(this.Parts, " ")
}
