package types

import "fmt"

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
)

func (t TokenType) String() string {
	switch t {
	case String:
		return "String"
	case Number:
		return "Number"
	case Assigner:
		return "Assigner"
	case Comparator:
		return "Comparator"
	case Operator:
		return "Operator"
	case Seperator:
		return "Seperator"
	case Bracket:
		return "Bracket"
	case Keyword:
		return "Keyword"
	case Variable:
		return "Variable"
	case BuiltIn:
		return "BuiltIn"
	case Comment:
		return "Comment"
	case Other:
		return "Other"
	default:
		return "Unknown"
	}
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
