package lexer

import (
	"fmt"

	"github.com/TwiN/go-color"
)

type itemType int

const (
	_ itemType = iota
	Text
	Keyword
	Identifier
	Operator
	String
	Number
	EOL
	EOF
	Error
)

func (typ itemType) String() string {
	return [...]string{
		"Text",
		"Keyword",
		"Identifier",
		"Operator",
		"String",
		"Number",
		"EOL",
		"EOF",
		"Error",
	}[typ-1]
}

type item struct {
	Typ itemType
	Val []rune
}

func (it *item) String() string {
	return fmt.Sprintf("%s: %s", it.Typ, string(it.Val))
}
func (it *item) ColoredString() string {
	return fmt.Sprintf("%s: %s", color.InCyan(it.Typ), string(it.Val))
}
