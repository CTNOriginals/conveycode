package lexer

import (
	"conveycode/compiler/tokenizer"
	"conveycode/constents"
	"fmt"

	"github.com/TwiN/go-color"
)

type ItemType int

const (
	_ ItemType = iota

	ItemText
	Keyword
	Value
	Condition
	Operator
	Identifier
	Scope
	Command
	Parameters
	Arguments
	ItemError
)

func (this ItemType) String() string {
	return [...]string{
		"ItemText",
		"Keyword",
		"Value",
		"Condition",
		"Operator",
		"Identifier",
		"Scope",
		"Command",
		"Parameters",
		"Arguments",
		"ItemError",
	}[this-1]
}

type Item struct {
	Typ    ItemType
	Tokens tokenizer.TokenList
}

func NewItem(typ ItemType, tokens ...tokenizer.Token) Item {
	return Item{
		Typ:    typ,
		Tokens: tokens,
	}
}

func (this Item) String() string {
	return fmt.Sprintf("%s\n  %s\n", color.InCyan(color.Bold+this.Typ.String()), this.Tokens.String())
}

// Returns all token values joint into a string without seperator
//
// This is mostly useful for items that usually only contain one token
func (this Item) ValueString() string {
	return this.Tokens.JoinValues("")
}

func (this Item) Compare(other Item) bool {
	return this.Typ == other.Typ && this.ValueString() == other.ValueString()
}

func (this Item) ItemFile() string {
	if this.IsError() || len(this.Tokens) == 0 {
		return constents.StringError
	}
	return this.Tokens[0].File
}
func (this Item) ItemLine() int {
	if this.IsError() || len(this.Tokens) == 0 {
		return -1
	}
	return this.Tokens[0].Line
}
func (this Item) ItemColumn() int {
	if this.IsError() || len(this.Tokens) == 0 {
		return -1
	}
	return this.Tokens[0].Column
}

func (this Item) IsError() bool {
	return this.Typ == ItemError
}
