package parser

import (
	"conveycode/compiler/lexer"
)

type parser struct {
	blocks  []lexer.Block
	Channel chan Instruction
}

func Parse(blocks []lexer.Block) (prs *parser) {
	prs = &parser{
		blocks:  blocks,
		Channel: make(chan Instruction),
	}

	return prs
}
