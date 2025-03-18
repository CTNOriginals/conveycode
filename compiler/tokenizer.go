package compiler

import (
	"conveycode/compiler/tools"
	"conveycode/compiler/types"
	"fmt"
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

func Tokenize(content []rune) types.TokenList {
	//? The tokens that are already identified in this line
	var tokens types.TokenList = types.NewTokenList()

	//? The current index in the line
	var cursor = tools.NewCursor(content)

	for !cursor.EOF {
		//? EOL
		if cursor.Peek() == 10 {
			tokens.Push(types.EOL, "")
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
			tokens.Push(types.Comment, string(cursor.ReadUntilFunc(func(c rune) bool {
				return c == '\n'
			})))

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

			tokens.Push(types.String, string(stream))
			continue
		}

		//? Number
		if unicode.IsDigit(cursor.Peek()) {
			tokens.Push(types.Number, string(cursor.ReadUntilFunc(func(c rune) bool {
				return !unicode.IsDigit(c)
			})))

			cursor.Seek(-1) //? the stream had to jump over the next character, so we bring it back here to not skip it
			continue
		}

		//? Scope
		if slices.Contains([]rune{'(', '[', '{'}, cursor.Peek()) {
			var openBracket = cursor.Peek()
			var closeBracket = getMatchingBracket(openBracket)
			var depth = 0

			cursor.Seek(1) //? Go past the opening bracket

			var body = cursor.ReadUntilFunc(func(c rune) bool {
				if c == closeBracket {
					if depth == 0 {
						return true
					} else {
						depth--
					}
				} else if c == openBracket {
					depth++
				}

				return false
			})

			tokens.Push(types.Scope, string(body))
			cursor.Seek(1) //? Go past the closing bracket
			continue
		}

		var singleTokenType types.TokenType = 0

		switch cursor.Peek() {
		case '(', ')', '[', ']', '{', '}': //TODO wrap all of the body after the open brace into a Scope token type
			singleTokenType = types.Bracket
		case '+', '-', '/', '*', '=', '>', '<', '!', '&', '|':
			singleTokenType = types.Operator
		case ',':
			singleTokenType = types.Seperator
		}

		if singleTokenType != 0 {
			tokens.Push(singleTokenType, string(cursor.Peek()))

			cursor.Read()
			continue
		}

		var stream = cursor.ReadUntilFunc(func(c rune) bool {
			return !regStream.MatchString(string(cursor.PeekPrev()))
		})

		if len(stream) == 0 {
			fmt.Printf("Unknown character: %v %d:%d\n", cursor.Peek(), cursor.Line, cursor.Column)

			cursor.Read()
			continue
		}

		tokens.Push(types.Other, string(stream))
		cursor.Seek(-1) //? the stream had to jump over the next character, so we bring it back here to not skip it

		continue
	}

	tokens.Push(types.EOF, "")

	return tokens
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
