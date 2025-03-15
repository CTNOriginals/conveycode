package types

import (
	"fmt"
)

type Instruction struct {
	Raw    string
	OpCode string
	Values []string
}

func NewInstruction(raw string, opCode string, values ...string) Instruction {
	return Instruction{Raw: raw, OpCode: opCode, Values: values}
}

func (i Instruction) String() string {
	return fmt.Sprintf(`
		Raw: %s
		OpCode: %s
		Values: %v
	`, i.Raw, i.OpCode, i.Values)
}

type Instructions = []Instruction
