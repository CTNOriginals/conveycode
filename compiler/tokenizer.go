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

// #region Handlers
type handler struct {
	test   func(cursor *tools.Cursor) bool
	handle func(cursor *tools.Cursor) (v []rune)
}

type handlerMap = map[types.TokenType]handler

var handlers = handlerMap{
	0: {
		test: func(cursor *tools.Cursor) bool {
			return unicode.IsSpace(cursor.Peek()) && cursor.Peek() != '\n'
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
			v = nil
			cursor.Read()
			return
		},
	},
	types.EOL: {
		test: func(cursor *tools.Cursor) bool {
			return cursor.Peek() == '\n'
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
			v = []rune{0}
			cursor.Read()
			return
		},
	},
	types.Comment: {
		test: func(cursor *tools.Cursor) bool {
			return cursor.Peek() == '/' && cursor.PeekNext() == '/'
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
			cursor.Seek(2) //? Move the cursor past the //
			return cursor.ReadUntilFunc(func(c rune) bool {
				return c == '\n'
			})
		},
	},
	types.String: {
		test: func(cursor *tools.Cursor) bool {
			return slices.Contains([]rune{'"', '\'', '`'}, cursor.Peek()) && cursor.PeekPrev() != '\\'
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
			var quote = cursor.Read()
			var stream []rune

			for {
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

			cursor.Seek(1) //? Skip the closing quote
			return stream
		},
	},
	types.Number: {
		test: func(cursor *tools.Cursor) bool {
			return unicode.IsDigit(cursor.Peek())
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
			return cursor.ReadUntilFunc(func(c rune) bool {
				return !unicode.IsDigit(c)
			})
		},
	},
	types.Scope: {
		test: func(cursor *tools.Cursor) bool {
			return slices.Contains([]rune{'(', '[', '{'}, cursor.Peek())
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
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

			if !cursor.Seek(1) { //? Go past the closing bracket
				cursor.Read() //? EOF is past the character we are trying to skip, so make sure EOF is set to true
			}

			return body
		},
	},
	types.Operator: {
		test: func(cursor *tools.Cursor) bool {
			return slices.Contains([]rune{'+', '-', '*', '/', '%', '=', '>', '<', '!', '&', '|'}, cursor.Peek())
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
			v = []rune{cursor.Peek()}
			cursor.Read()
			return
		},
	},
	types.Seperator: {
		test: func(cursor *tools.Cursor) bool {
			return cursor.Peek() == ','
		},
		handle: func(cursor *tools.Cursor) (v []rune) {
			v = []rune{cursor.Peek()}
			cursor.Read()
			return
		},
	},
}

// Hold the keys in the order that they are defined as in the enum
var handlerKeys []types.TokenType = make([]types.TokenType, 0, len(handlers))

//#endregion

func init() {
	var err error
	if regStream, err = regexp.Compile("\\w"); err != nil {
		panic(err)
	}

	for tt := range handlers {
		handlerKeys = append(handlerKeys, tt)
	}

	slices.Sort(handlerKeys)
}

func Tokenize(content []rune) types.TokenList {
	//? The tokens that are already identified in this line
	var tokens types.TokenList = types.NewTokenList()

	//? The current index in the line
	var cursor = tools.NewCursor(content)

	for !cursor.EOF {
		var handled = false
		for _, tt := range handlerKeys {
			h := handlers[tt]
			if h.test(cursor) {
				v := h.handle(cursor)
				handled = true

				if v != nil {
					tokens.Push(tt, v...)
				}
				break
			}
		}

		if handled {
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
	}

	tokens.Push(types.EOF, 0)

	return tokens
}

// #region Utilities
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

//#endregion
