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
	Comparator

	Bracket
	Seperator

	Keyword
	Variable
	BuiltIn

	Comment

	Other

	EOL
	EOF
)

func (e TokenType) String() string {
	return [...]string{
		"String",
		"Number",

		"Assigner",
		"Operator",
		"Comparator",

		"Bracket",
		"Seperator",

		"Keyword",
		"Variable",
		"BuiltIn",

		"Comment",

		"Other",

		"EOL",
		"EOF",
	}[e-1]
}

type Token struct {
	TokenType TokenType
	Value     string
}

func (token Token) String() string {
	return fmt.Sprintf("%s: %s", token.TokenType.String(), token.Value)
}

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
