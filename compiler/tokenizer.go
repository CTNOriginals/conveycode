package compiler

import (
	"conveycode/compiler/tools"
	"conveycode/compiler/types"
	"fmt"
	"regexp"
	"slices"
	"unicode"

	"github.com/TwiN/go-color"
)

var regStream *regexp.Regexp

func init() {
	var err error
	if regStream, err = regexp.Compile("\\w"); err != nil {
		panic(err)
	}
}

func Tokenize(content []rune) types.TokenList {
	//? The tokens that are already identified in this line
	var tokens types.TokenList = types.NewTokenList()

	//? The current index in the line
	var cursor = tools.NewCursor(content)

	for !cursor.EOF {
		//? EOL
		if cursor.Peek() == '\n' {
			tokens.Push(types.EOL, 0)
			cursor.Read()
			continue
		}

		//? Whitespace skip
		if unicode.IsSpace(cursor.Peek()) {
			cursor.Read()
			continue
		}

		//? Comment
		if cursor.Peek() == '/' && cursor.PeekNext() == '/' {
			cursor.Seek(2) //? Move the cursor past the //
			tokens.Push(types.Comment, cursor.ReadUntilFunc(func(c rune) bool {
				return c == '\n'
			})...)

			continue
		}

		//? String
		if slices.Contains([]rune{'"', '\'', '`'}, cursor.Peek()) && cursor.PeekPrev() != '\\' {
			var quote = cursor.Read()
			var stream []rune

			for true {
				s := cursor.ReadUntilFunc(func(c rune) bool {
					return c == quote
				})

				stream = append(stream, s...)
				if stream[len(stream)-1] == '\\' {
					stream = append(stream, quote)
				} else {
					break
				}
			}

			tokens.Push(types.String, stream...)
			continue
		}

		//? Number
		if unicode.IsDigit(cursor.Peek()) {
			tokens.Push(types.Number, cursor.ReadUntilFunc(func(c rune) bool {
				return !unicode.IsDigit(c)
			})...)

			continue
		}

		//? Scope
		if slices.Contains([]rune{'(', '[', '{'}, cursor.Peek()) {
			var openBracket = cursor.Peek()
			var startLocation = []int{cursor.Line, cursor.Column}
			var closeBracket = getMatchingBracket(openBracket)
			var depth = 0

			cursor.Seek(1) //? Go past the opening bracket

			var body = cursor.ReadUntilFunc(func(c rune) bool {
				switch c {
				case closeBracket:
					if depth == 0 {
						return true
					} else {
						depth--
					}
				case openBracket:
					depth++
				default:
					break
				}

				if cursor.EOF {
					fmt.Println(formatError("Unmatched bracket", openBracket, startLocation[0], startLocation[1]))
					return true
				}

				return false
			})

			tokens.Push(types.Scope, body...)
			cursor.Seek(1) //? Go past the closing bracket
			continue
		}

		var singleTokenType types.TokenType = 0

		switch cursor.Peek() {
		case '+', '-', '/', '*', '=', '>', '<', '!', '&', '|':
			singleTokenType = types.Operator
		case ',':
			singleTokenType = types.Seperator
		}

		if singleTokenType != 0 {
			tokens.Push(singleTokenType, cursor.Peek())

			cursor.Read()
			continue
		}

		var stream = cursor.ReadUntilFunc(func(c rune) bool {
			return !regStream.MatchString(string(c))
		})

		if len(stream) == 0 {
			fmt.Println(formatError("Unknown character", cursor.Peek(), cursor.Line, cursor.Column))

			cursor.Read()
			continue
		}

		tokens.Push(types.Other, stream...)

		continue
	}

	tokens.Push(types.EOF, 0)

	return tokens
}

func formatError(message string, char rune, line int, column int) string {
	return fmt.Sprintf(color.InRed("%s \"%s\" at %s:%s"), message, string(char), color.InYellow(line), color.InYellow(column))
}

func getMatchingBracket(bracket rune) rune {
	switch bracket {
	case '(':
		return ')'
	case ')':
		return '('

	case '[':
		return ']'
	case ']':
		return '['

	case '{':
		return '}'
	case '}':
		return '{'

	default:
		fmt.Printf("Unexpected character \"%s\", no matching bracked", string(bracket))
		return 0
	}
}
