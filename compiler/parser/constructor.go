package parser

import (
	"conveycode/compiler/lexer"
)

type Constructor = func(block lexer.Block) []Instruction

type ConstructorMap = map[lexer.BlockType]Constructor

var constructors = ConstructorMap{
	lexer.Assignment: func(block lexer.Block) []Instruction {
		return []Instruction{NewInstruction([]string{"some", "parts", "here"})}
	},
}
