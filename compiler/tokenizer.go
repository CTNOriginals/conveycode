package compiler

import (
	"conveycode/compiler/tools"
	"conveycode/compiler/types"
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
	// fmt.Println(cursor.Content)

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

			cursor.Read() //? This skips the closing quote charactr
			continue
		}

		//? Number
		if unicode.IsDigit(cursor.Peek()) {
			tokens.Push(types.Number, string(cursor.ReadUntilFunc(func(c rune) bool {
				return !unicode.IsDigit(c)
			})))

			continue
		}

		cursor.Read()
		continue
	}

	tokens.Push(types.EOF, "")

	return tokens
}

// func tokenizeString(cursor *tools.Cursor) (c int, token types.Token) {
// 	var quote rune = cursor.Peak(0)
// 	var value []rune

// 	cursor++

// 	for line[cursor] != quote || line[utils.Max(0, cursor-1)] == '\\' {
// 		value = append(value, line[cursor])
// 		cursor++
// 	}

// 	cursor++

// 	return cursor, types.Token{
// 		TokenType: types.String,
// 		Value:     string(value),
// 	}
// }
// func tokenizeNumber(cursor int, line []rune) (c int, token types.Token) {
// 	var value []rune

// 	for unicode.IsDigit(line[cursor]) {
// 		value = append(value, line[cursor])
// 		cursor++
// 		if cursor >= len(line) {
// 			break
// 		}
// 	}

// 	return cursor, types.Token{
// 		TokenType: types.Number,
// 		Value:     string(value),
// 	}
// }

func identifyStream(stream []rune) types.TokenType {
	if slices.Contains(types.Keywords, string(stream)) {
		return types.Keyword
	} else if slices.Contains(types.BuiltIns, string(stream)) {
		return types.BuiltIn
	} else {
		return types.Other
	}
}
