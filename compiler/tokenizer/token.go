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
	Boolean
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
	TokenError
)

func (this TokenType) String() string {
	return [...]string{
		"EOL",
		"Comment",
		"Boolean",
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
		"TokenError",
	}[this-1]
}

var ValueTokenTypes = []TokenType{String, Number, Boolean, Text}

// #region Token
type Token struct {
	Typ    TokenType
	Val    []rune
	File   string
	Line   int
	Column int
}

func NewToken(file string, typ TokenType, val []rune, line int, col int) Token {
	return Token{
		Typ:    typ,
		Val:    val,
		File:   file,
		Line:   line,
		Column: col,
	}
}

func (this Token) String() string {
	return fmt.Sprintf("%s:%s %s: %s", color.InYellow(this.Line), color.InYellow(this.Column), color.InGreen(this.Typ), string(this.Val))
}

func (this Token) GetTypeColor() string {
	switch this.Typ {
	case String:
		return color.Red
	case Number:
		return color.Green
	case Boolean:
		return color.Blue
	case Operator:
		return color.Cyan
	case Command:
		return color.Purple
	case Seperator:
		return color.Gray
	case RoundL, RoundR, SquareL, SquareR, CurlyL, CurlyR:
		return color.Yellow
	}

	return ""
}

func (this Token) ColoredValue() string {
	return this.GetTypeColor() + string(this.Val)
}
