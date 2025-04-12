package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"
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

func ParseContent(file string, content string, scope *Scope) (prs *parser) {
	var tokens = tokenizer.TokenizeContent(file, content)
	var blocks = lexer.Lex(tokens).Construct()
	return Parse(blocks, scope)
}

func ParseItem(item lexer.Item, scope *Scope) (prs *parser) {
	return Parse(lexer.Lex(item.Tokens).Construct(), scope)
}
