package tokenizer

import "fmt"

// #region Class Token
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

	Other
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

//#region Token
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

//#endregion

//#region Token List
type TokenList []Token

func NewTokenList() TokenList {
	return make(TokenList, 0)
}

func (tl TokenList) String() (str string) {
	for _, token := range tl {
		str += token.String() + "\n"
	}

	return
}

func (tl *TokenList) Push(t TokenType, v ...rune) {
	*tl = append(*tl, NewToken(t, v))
	// fmt.Println(NewToken(t, v))
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
