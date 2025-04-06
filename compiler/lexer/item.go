package lexer

import (
	"conveycode/compiler/tokenizer"
	"fmt"

	"github.com/TwiN/go-color"
)

type itemType int

const (
	_ itemType = iota

	ItemText
	Keyword
	Value
	Condition
	Operator
	Identifier
	Scope
	Command
	Arguments
	ItemError
)

func (this itemType) String() string {
	return [...]string{
		"ItemText",
		"Keyword",
		"Value",
		"Condition",
		"Operator",
		"Identifier",
		"Scope",
		"Command",
		"Arguments",
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

// Returns all token values joint into a string without seperator
//
// This is mostly useful for items that usually only contain one token
func (this item) ValueString() string {
	return this.Tokens.JoinValues("")
}

func (this item) Compare(other item) bool {
	return this.Typ == other.Typ && this.ValueString() == other.ValueString()
}

func (this item) ItemFile() string {
	return this.Tokens[0].File
}
func (this item) ItemLine() int {
	return this.Tokens[0].Line
}
func (this item) ItemColumn() int {
	return this.Tokens[0].Column
}
