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

func (token Token) String() string {
	return fmt.Sprintf("%s: %s", token.tokenType.String(), token.value)
}

func Tokenize(lines []string) [][]Token {
	var tokenLines [][]Token

	for _, rawLine := range lines {
		var line []rune = []rune(rawLine)
		var tokens []Token
		var cursor int = 0

		fmt.Println(rawLine, len(line))
		for cursor < len(line) {
			var char rune = line[cursor]

			if unicode.IsSpace(rune(char)) {
				cursor++
				continue
			}

			if slices.Contains([]rune{'"', '\'', '`'}, char) && line[cursor-1] != '\\' {
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
				char = line[cursor]

				tokens = append(tokens, Token{
					tokenType: String,
					value:     string(value),
				})

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
