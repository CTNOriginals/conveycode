package types

import "fmt"

const (
	// End Of Transmission.
	// Used as an End Of File (EOF) indicator character
	EOT = 0x00
)

// #region Class Token
type TokenType int

const (
	_ TokenType = iota
	String
	Number

	Assigner
	Operator

	Bracket
	Seperator

	// Keyword
	// BuiltIn

	Scope

	Comment

	Other

	EOL
	EOF
)

func (this TokenType) String() string {
	return [...]string{
		"String",
		"Number",

		"Assigner",
		"Operator",

		"Bracket",
		"Seperator",

		// "Keyword",
		// "BuiltIn",

		"Scope",

		"Comment",

		"Other",

		"EOL",
		"EOF",
	}[this-1]
}

//#region Token
type Token struct {
	TokenType TokenType
	Value     []rune
}

func NewToken(t TokenType, v []rune) Token {
	return Token{
		TokenType: t,
		Value:     v,
	}
}

func (token Token) String() string {
	return fmt.Sprintf("%s: %s", token.TokenType, string(token.Value))
}

//#endregion

//#region Token List
type TokenList struct {
	Tokens []Token
}

func NewTokenList() TokenList {
	return TokenList{}
}

func (this *TokenList) String() (str string) {
	for _, token := range this.Tokens {
		str += token.String() + "\n"
	}

	return
}

func (this *TokenList) Push(t TokenType, v ...rune) {
	this.Tokens = append(this.Tokens, NewToken(t, v))
	fmt.Println(NewToken(t, v).String())
}

//#endregion

//#endregion

// #region Constents
var Keywords = []string{
	"var",
	"func",

	"if",
	"else",

	"continue",
	"break",
	"return",
}
var BuiltIns = []string{
	"print",
	"println",
	"flush",
}

//#endregion
