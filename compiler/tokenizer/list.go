package tokenizer

import "strings"

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

func (this *TokenList) Push(t TokenType, v ...rune) {
	*this = append(*this, NewToken(t, v))
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
