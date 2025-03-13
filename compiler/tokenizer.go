package compiler

import (
	"conveycode/compiler/utils"
	"fmt"
	"log"
	"regexp"
	"slices"
	"unicode"
)

// #region Class Token
type TokenType int

const (
	_ TokenType = iota
	String
	Number

	Operator
	Comparator
	Seperator

	Bracket

	Keyword
	Variable
	BuiltIn

	Method
	Parameter
	Call
	Argument

	Comment

	Other
)

func (t TokenType) String() string {
	switch t {
	case String:
		return "String"
	case Number:
		return "Number"
	case Comparator:
		return "Comparator"
	case Operator:
		return "Operator"
	case Seperator:
		return "Seperator"
	case Bracket:
		return "Bracket"
	case Keyword:
		return "Keyword"
	case Variable:
		return "Variable"
	case BuiltIn:
		return "BuiltIn"
	case Method:
		return "Method"
	case Parameter:
		return "Parameter"
	case Call:
		return "Call"
	case Argument:
		return "Argument"
	case Comment:
		return "Comment"
	case Other:
		return "Other"
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

//#endregion

// #region Constents
var keywords = []string{
	"var",
	"func",

	"if",
	"else",

	"continue",
	"break",
	"return",
}
var builtIn = []string{
	"print",
	"println",
	"flush",
}

//#endregion

var regStream *regexp.Regexp

func init() {
	var err error
	if regStream, err = regexp.Compile("\\w"); err != nil {
		panic(err)
	}
}

func Tokenize(lines []string) [][]Token {
	var tokenLines [][]Token
	// var variables []string

	for _, rawLine := range lines {
		var line []rune = []rune(rawLine)

		//? The tokens that are already identified in this line
		var tokens []Token

		//? The current index in the line
		var cursor int = 0

		for cursor < len(line) {
			var char rune = line[cursor]

			if unicode.IsSpace(char) {
				cursor++
				continue
			}

			if char == '/' && line[cursor+1] == '/' {
				tokens = append(tokens, Token{
					tokenType: Comment,
					value:     string(line[cursor+2:]),
				})
				break
			}

			if slices.Contains([]rune{'"', '\'', '`'}, char) && line[utils.Max(0, cursor-1)] != '\\' {
				c, token := tokenizeString(cursor, line)
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

			var singleTokenType TokenType = 0

			switch char {
			case '(', ')', '[', ']', '{', '}':
				singleTokenType = Bracket
			case '=':
				if utils.ContainsListItem([]rune{'=', '>', '<', '!'}, []rune{line[cursor-1], line[cursor+1]}) {
					singleTokenType = Comparator
				} else {
					singleTokenType = Operator
				}
			case '+', '-', '/', '*':
				singleTokenType = Operator
			case '>', '<', '!', '&', '|':
				singleTokenType = Comparator
			case ',':
				singleTokenType = Seperator
			}

			if singleTokenType != 0 {
				tokens = append(tokens, Token{
					tokenType: singleTokenType,
					value:     string(char),
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
				if cursor >= len(line) {
					break
				}
				char = line[cursor]
			}

			if stream == nil {
				log.Panicf("Unknown character: %s", string(char))
				cursor++
				continue
			}

			tokens = append(tokens, Token{
				tokenType: identifyStream(stream),
				value:     string(stream),
			})
		}

		if len(tokens) == 0 {
			continue
		}

		tokenLines = append(tokenLines, tokens)
	}

	return tokenLines
}

func tokenizeString(cursor int, line []rune) (c int, token Token) {
	var quote rune = line[cursor]
	var value []rune

	cursor++

	for line[cursor] != quote || line[utils.Max(0, cursor-1)] == '\\' {
		value = append(value, line[cursor])
		cursor++
	}

	// value = append(value, char)
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

func identifyStream(stream []rune) TokenType {
	if slices.Contains(keywords, string(stream)) {
		return Keyword
	} else if slices.Contains(builtIn, string(stream)) {
		return BuiltIn
	} else {
		return Other
	}
}
