package parser

import (
	"conveycode/compiler/lexer"
)

type parser struct {
	blocks []lexer.Block
}

func Parse(blocks []lexer.Block) (prs *parser) {
	prs = &parser{
		blocks: blocks,
	}

	return prs
}

func (this *parser) Construct() (instructions []Instruction) {
	for _, block := range this.blocks {
		var constructor, ok = constructors[block.Typ]

		//- Does the block type have a constructor defined
		if !ok {
			continue
		}

		instructions = append(instructions, constructor(block)...)
	}

	return instructions
}
