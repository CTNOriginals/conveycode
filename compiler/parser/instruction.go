package parser

import "strings"

type Instruction struct {
	Parts []string
}

func NewInstruction(parts ...string) Instruction {
	return Instruction{
		Parts: parts,
	}
}

func (this Instruction) String() string {
	return strings.Join(this.Parts, " ")
}

func (this *Instruction) Push(parts ...string) {
	this.Parts = append(this.Parts, parts...)
}

func (this Instruction) isLabel() bool {
	if len(this.Parts) != 1 {
		return false
	}

	var part = this.Parts[0]
	return len(strings.Split(part, "_")) > 1 && rune(part[len(part)-1]) == ':'
}
