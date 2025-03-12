package compiler

import (
	"fmt"
	"slices"
	"unicode"
)

type TokenType int

const (
	_ TokenType = iota
	String
	Number
	Operator
	Bracket
	Keyword
)

func (t TokenType) String() string {
	switch t {
	case String:
		return "String"
	case Number:
		return "Number"
	case Operator:
		return "Operator"
	case Bracket:
		return "Bracket"
	case Keyword:
		return "Keyword"
	default:
		return "Unknown"
	}
}

type Token struct {
	tokenType TokenType
	value     string
}

var keywords = []string{
	"var",

	"println",
	"print",

	"func",

	"if",
	"else if",
	"else",
}

func (token Token) String() string {
	return fmt.Sprintf("%s: %s", token.tokenType.String(), token.value)
}

func Tokenize(lines []string) [][]Token {
	var tokenLines [][]Token

	for _, rawLine := range lines {
		var line []rune = []rune(rawLine)
		var tokens []Token
		var cursor int = 0

		// fmt.Println(rawLine, len(line))
		for cursor < len(line) {
			var char rune = line[cursor]

			if unicode.IsSpace(rune(char)) {
				cursor++
				continue
			}

			if slices.Contains([]rune{'"', '\'', '`'}, char) && line[cursor-1] != '\\' {
				c, token := tokenizeString(cursor, char, line)
				cursor = c
				tokens = append(tokens, token)
				continue
			}

			if unicode.IsDigit(char) {
				c, token := tokenizeNumber(cursor, char, line)
				cursor = c
				tokens = append(tokens, token)
				continue
			}

			// log.Panicf("Unknown character: %c\n", char)
			cursor++
		}

		if len(tokens) == 0 {
			continue
		}

		tokenLines = append(tokenLines, tokens)
	}

	return tokenLines
}

func tokenizeString(cursor int, char rune, line []rune) (c int, token Token) {
	var quote rune = char
	var value []rune = []rune{quote}

	cursor++
	char = line[cursor]

	for char != quote || line[cursor-1] == '\\' {
		value = append(value, char)
		cursor++
		char = line[cursor]
	}

	value = append(value, char)
	cursor++

	return cursor, Token{
		tokenType: String,
		value:     string(value),
	}
}
func tokenizeNumber(cursor int, char rune, line []rune) (c int, token Token) {
	var value []rune

	for unicode.IsDigit(char) {
		value = append(value, char)
		cursor++
		if cursor >= len(line) {
			break
		}
		char = line[cursor]
	}

	return cursor, Token{
		tokenType: Number,
		value:     string(value),
	}
}
