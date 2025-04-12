package tokenizer

import (
	"fmt"
	"slices"

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
	Link

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
		"Link",

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

var ValueTokenTypes = []TokenType{String, Number, Boolean, Link, Text}
var OpenBracketTypes = []TokenType{RoundL, SquareL, CurlyL}
var CloseBracketTypes = []TokenType{RoundR, SquareR, CurlyR}

func (this TokenType) IsValue() bool {
	return slices.Contains(ValueTokenTypes, this)
}
func (this TokenType) IsOpenBracket() bool {
	return slices.Contains(OpenBracketTypes, this)
}
func (this TokenType) IsCloseBracket() bool {
	return slices.Contains(CloseBracketTypes, this)
}

func (this TokenType) GetMatchingBracket() TokenType {
	switch this {
	case RoundL:
		return RoundR
	case SquareL:
		return SquareR
	case CurlyL:
		return CurlyR

	case RoundR:
		return RoundL
	case SquareR:
		return SquareL
	case CurlyR:
		return CurlyL
	}

	return 0
}

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

func (this Token) Compare(other Token) bool {
	return this.Typ == other.Typ &&
		string(this.Val) == string(other.Val) &&
		this.File == other.File &&
		this.Line == other.Line &&
		this.Column == other.Column
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
