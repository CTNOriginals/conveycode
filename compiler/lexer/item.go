package lexer

import (
	"conveycode/compiler/tokenizer"
	"fmt"
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
	Typ   itemType
	Token tokenizer.TokenList
}

func (this item) String() string {
	return fmt.Sprintf("%s: %s", this.Typ, this.Token.String())
}
