package tokenizer

import (
	"fmt"

	"github.com/TwiN/go-color"
)

type TokenType int

const (
	_ TokenType = iota

	EOL
	Comment
	String
	Number
	Operator
	Seperator

	RoundL
	RoundR
	SquareL
	SquareR
	CurlyL
	CurlyR

	Text
	EOF
)

func (this TokenType) String() string {
	return [...]string{
		"EOL",
		"Comment",
		"String",
		"Number",
		"Operator",
		"Seperator",

		"RoundL",
		"RoundR",
		"SquareL",
		"SquareR",
		"CurlyL",
		"CurlyR",

		"Other",
		"EOF",
	}[this-1]
}

// #region Token
type Token struct {
	Typ TokenType
	Val []rune
}

func NewToken(t TokenType, v []rune) Token {
	return Token{
		Typ: t,
		Val: v,
	}
}

func (token Token) String() string {
	return fmt.Sprintf("%s: %s", token.Typ, string(token.Val))
}

func (this Token) ColoredValue() string {
	switch this.Typ {
	case String:
		return color.Colorize(color.Red, string(this.Val))
	case Number:
		return color.Colorize(color.Green, string(this.Val))
	case Operator:
		return color.Colorize(color.Blue, string(this.Val))
	case Seperator:
		return color.Colorize(color.Cyan, string(this.Val))
	case RoundL, RoundR, SquareL, SquareR, CurlyL, CurlyR:
		return color.Colorize(color.Yellow, string(this.Val))
	}

	return string(this.Val)
}
