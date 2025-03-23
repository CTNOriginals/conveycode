package tokenizer

import (
	"fmt"
	"regexp"
	"slices"
	"unicode"

	"github.com/TwiN/go-color"
)

var regStream *regexp.Regexp

// #region Handlers
type handler struct {
	test   func(cursor *Cursor) bool
	handle func(cursor *Cursor) (v []rune)

	// If the type is just a simple set of specific characters,
	// Populate them in here and leave the other fields nil
	runes []rune
}

type handlerMap = map[TokenType]handler

var handlers = handlerMap{
	0: {
		test: func(cursor *Cursor) bool {
			return unicode.IsSpace(cursor.Peek()) && cursor.Peek() != '\n'
		},
		handle: func(cursor *Cursor) (v []rune) {
			v = nil
			cursor.Read()
			return
		},
	},
	EOL: {
		test: func(cursor *Cursor) bool {
			return cursor.Peek() == '\n'
		},
		handle: func(cursor *Cursor) (v []rune) {
			v = []rune{0}
			cursor.Read()
			return
		},
	},
	Comment: {
		test: func(cursor *Cursor) bool {
			return cursor.Peek() == '/' && cursor.PeekNext() == '/'
		},
		handle: func(cursor *Cursor) (v []rune) {
			cursor.Seek(2) //? Move the cursor past the //
			return cursor.ReadUntilFunc(func(c rune) bool {
				return c == '\n'
			})
		},
	},
	String: {
		test: func(cursor *Cursor) bool {
			return slices.Contains([]rune{'"', '\'', '`'}, cursor.Peek()) && cursor.PeekPrev() != '\\'
		},
		handle: func(cursor *Cursor) (v []rune) {
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
	Number: {
		test: func(cursor *Cursor) bool {
			return unicode.IsDigit(cursor.Peek())
		},
		handle: func(cursor *Cursor) (v []rune) {
			return cursor.ReadUntilFunc(func(c rune) bool {
				return !unicode.IsDigit(c) && c != '.'
			})
		},
	},

	Operator:  {test: nil, handle: nil, runes: []rune{'+', '-', '*', '/', '%', '=', '>', '<', '!', '&', '|'}},
	Seperator: {test: nil, handle: nil, runes: []rune{','}},
	RoundL:    {test: nil, handle: nil, runes: []rune{'('}},
	RoundR:    {test: nil, handle: nil, runes: []rune{')'}},
	SquareL:   {test: nil, handle: nil, runes: []rune{'['}},
	SquareR:   {test: nil, handle: nil, runes: []rune{']'}},
	CurlyL:    {test: nil, handle: nil, runes: []rune{'{'}},
	CurlyR:    {test: nil, handle: nil, runes: []rune{'}'}},
}

// Hold the keys in the order that they are defined as in the enum
var handlerKeys []TokenType = make([]TokenType, 0, len(handlers))

//#endregion

func init() {
	var err error
	if regStream, err = regexp.Compile(`\w`); err != nil {
		panic(err)
	}

	for tt := range handlers {
		handlerKeys = append(handlerKeys, tt)
	}

	slices.Sort(handlerKeys)
}

func Tokenize(content []rune) TokenList {
	//? The tokens that are already identified in this line
	var tokens TokenList = NewTokenList()

	//? The current index in the line
	var cursor = NewCursor(content)

	for !cursor.EOF {
		var handled = false
		for _, typ := range handlerKeys {
			hand := handlers[typ]

			if hand.test == nil && hand.runes != nil {
				if slices.Contains(hand.runes, cursor.Peek()) {
					tokens.Push(typ, cursor.Read())
					handled = true
					break
				}
			} else if hand.test(cursor) {
				val := hand.handle(cursor)
				handled = true

				if val != nil {
					tokens.Push(typ, val...)
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

		tokens.Push(Other, stream...)
	}

	tokens.Push(EOF, 0)

	return tokens
}

// #region Utilities
func formatError(message string, char rune, line int, column int) string {
	return fmt.Sprintf(color.InRed("%s \"%s\" at %s:%s"), message, string(char), color.InYellow(line), color.InYellow(column))
}

// func getMatchingBracket(bracket rune) rune {
// 	switch bracket {
// 	case '(':
// 		return ')'
// 	case ')':
// 		return '('

// 	case '[':
// 		return ']'
// 	case ']':
// 		return '['

// 	case '{':
// 		return '}'
// 	case '}':
// 		return '{'

// 	default:
// 		fmt.Printf("Unexpected character \"%s\", no matching bracked", string(bracket))
// 		return 0
// 	}
// }

//#endregion
