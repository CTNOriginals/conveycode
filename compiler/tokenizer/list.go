package tokenizer

import (
	"slices"
	"strings"
)

type TokenList []Token

func NewTokenList() TokenList {
	return make(TokenList, 0)
}

func (this TokenList) String() (str string) {
	var list = make([]string, len(this))
	for i, token := range this {
		list[i] = token.String()
	}

	return strings.Join(list, "\n  ")
}

func (this TokenList) Contains(typ TokenType) bool {
	for _, token := range this {
		if token.Typ == typ {
			return true
		}
	}

	return false
}

func (this TokenList) FindTokenByType(typ TokenType) Token {
	for _, token := range this {
		if token.Typ == typ {
			return token
		}
	}

	return Token{
		Typ: TokenError,
	}
}

// Returns the stream of values contained in the list
func (this TokenList) Stream() (str string) {
	for _, token := range this {
		if token.Typ == EOL {
			str += "\n"
			continue
		}

		str += token.String()
	}

	return str
}
func (this TokenList) ColoredStream() (str string) {
	for _, token := range this {
		if token.Typ == EOL {
			str += "\n"
			continue
		}

		str += token.ColoredValue()
	}

	return str
}

func (this *TokenList) Push(file string, cursor Cursor, typ TokenType, val ...rune) {
	var col = cursor.Column - len(string(val))

	if cursor.EOF {
		col++
	}

	*this = append(*this, NewToken(file, typ, val, cursor.Line, col))
	// fmt.Println(NewToken(t, v))
}

func (this TokenList) Values() (ret [][]rune) {
	ret = make([][]rune, len(this))
	for i, token := range this {
		ret[i] = token.Val
	}

	return ret
}

func (this TokenList) ValuesAsString() (ret []string) {
	ret = make([]string, len(this))
	for i, token := range this {
		ret[i] = string(token.Val)
	}

	return ret
}

func (this TokenList) JoinValues(seperator string) string {
	return strings.Join(this.ValuesAsString(), seperator)
}

// type TokenList []Token
func (this *TokenList) Remove(start int, count int) {
	*this = slices.Delete(*this, start, start+count)
}
