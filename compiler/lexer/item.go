package lexer

import (
	"conveycode/compiler/tokenizer"
	"fmt"

	"github.com/TwiN/go-color"
)

type itemType int

const (
	_ itemType = iota

	Keyword
	Value
	Condition
	Operator
	Identifier
	Scope
	ItemError
)

func (this itemType) String() string {
	return [...]string{
		"Keyword",
		"Value",
		"Condition",
		"Operator",
		"Identifier",
		"Scope",
		"ItemError",
	}[this-1]
}

type item struct {
	Typ    itemType
	Tokens tokenizer.TokenList
}

func NewItem(typ itemType, tokens ...tokenizer.Token) item {
	return item{
		Typ:    typ,
		Tokens: tokens,
	}
}

func (this item) String() string {
	return fmt.Sprintf("%s\n  %s\n", color.InCyan(color.Bold+this.Typ.String()), this.Tokens.String())
}

func (this *item) push(token ...tokenizer.Token) {
	this.Tokens = append(this.Tokens, token...)
}
