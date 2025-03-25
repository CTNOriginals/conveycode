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
	Command

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
		"Command",

		"RoundL",
		"RoundR",
		"SquareL",
		"SquareR",
		"CurlyL",
		"CurlyR",

		"Text",
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

func (this Token) String() string {
	return fmt.Sprintf("%s: %s", color.InGreen(this.Typ), string(this.Val))
}

func (this Token) GetTypeColor() string {
	switch this.Typ {
	case String:
		return color.Red
	case Number:
		return color.Green
	case Operator:
		return color.Blue
	case Seperator:
		return color.Cyan
	case RoundL, RoundR, SquareL, SquareR, CurlyL, CurlyR:
		return color.Yellow
	}

	return ""
}

func (this Token) ColoredValue() string {
	return this.GetTypeColor() + string(this.Val)
}
