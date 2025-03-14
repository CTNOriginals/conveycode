package compiler

import (
	"conveycode/compiler/types"
	"conveycode/compiler/utils"
	"log"
	"regexp"
	"slices"
	"unicode"
)

var regStream *regexp.Regexp

func init() {
	var err error
	if regStream, err = regexp.Compile("\\w"); err != nil {
		panic(err)
	}
}

func Tokenize(content []rune) []types.Token {
	//? The tokens that are already identified in this line
	var tokens []types.Token

	//? The current index in the line
	var cursor int = 0

	for cursor < len(content) {
		var char rune = content[cursor]

		//? Whitespace skip
		if unicode.IsSpace(char) {
			cursor++
			continue
		}

		//? Comments
		if char == '/' && content[cursor+1] == '/' {
			tokens = append(tokens, types.Token{
				TokenType: types.Comment,
				Value:     string(content[cursor+2:]),
			})
			break
		}

		//? Strings
		if slices.Contains([]rune{'"', '\'', '`'}, char) && content[utils.Max(0, cursor-1)] != '\\' {
			c, token := tokenizeString(cursor, content)
			cursor = c
			tokens = append(tokens, token)
			continue
		}

		//? Numbers
		if unicode.IsDigit(char) {
			c, token := tokenizeNumber(cursor, content)
			cursor = c
			tokens = append(tokens, token)
			continue
		}

		var singleTokenType types.TokenType = 0

		switch char {
		case '(', ')', '[', ']', '{', '}':
			singleTokenType = types.Bracket
		case '=':
			if utils.ContainsListItem([]rune{'=', '>', '<', '!'}, []rune{content[cursor-1], content[cursor+1]}) {
				singleTokenType = types.Comparator
			} else {
				singleTokenType = types.Assigner
			}
		case '+', '-', '/', '*':
			singleTokenType = types.Operator
		case '>', '<', '!', '&', '|':
			singleTokenType = types.Comparator
		case ',':
			singleTokenType = types.Seperator
		}

		if singleTokenType != 0 {
			tokens = append(tokens, types.Token{
				TokenType: singleTokenType,
				Value:     string(char),
			})
			cursor++
			continue
		}

		//? Nothing was found to match at this point
		//? All the following characters will be put in a stream until the next character is anything that can make its own token again

		var stream []rune

		for regStream.Match([]byte{byte(char)}) {
			stream = append(stream, char)
			cursor++
			if cursor >= len(content) {
				break
			}
			char = content[cursor]
		}

		if stream == nil {
			log.Panicf("Unknown character: %s", string(char))
			cursor++
			continue
		}

		tokens = append(tokens, types.Token{
			TokenType: identifyStream(stream),
			Value:     string(stream),
		})
	}

	return tokens
}

func tokenizeString(cursor int, line []rune) (c int, token types.Token) {
	var quote rune = line[cursor]
	var value []rune

	cursor++

	for line[cursor] != quote || line[utils.Max(0, cursor-1)] == '\\' {
		value = append(value, line[cursor])
		cursor++
	}

	cursor++

	return cursor, types.Token{
		TokenType: types.String,
		Value:     string(value),
	}
}
func tokenizeNumber(cursor int, line []rune) (c int, token types.Token) {
	var value []rune

	for unicode.IsDigit(line[cursor]) {
		value = append(value, line[cursor])
		cursor++
		if cursor >= len(line) {
			break
		}
	}

	return cursor, types.Token{
		TokenType: types.Number,
		Value:     string(value),
	}
}

func identifyStream(stream []rune) types.TokenType {
	if slices.Contains(types.Keywords, string(stream)) {
		return types.Keyword
	} else if slices.Contains(types.BuiltIns, string(stream)) {
		return types.BuiltIn
	} else {
		return types.Other
	}
}
