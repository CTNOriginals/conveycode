package parser

import (
	"conveycode/compiler/lexer"
)

type parser struct {
	blocks []lexer.Block
	scope  *Scope
}

func Parse(blocks []lexer.Block, scope *Scope) (prs *parser) {
	prs = &parser{
		blocks: blocks,
		scope:  scope,
	}

	return prs
}

func (this parser) getBlock(index int) lexer.Block {
	return this.blocks[index]
}
