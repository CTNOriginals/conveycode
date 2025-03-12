package compiler

import (
	"strings"
	"unicode"
)

func Tokenize(lines []string) [][][]rune {
	var instructions [][][]rune

	for _, line := range lines {
		var tokens [][]rune
		var current []rune

		quoteScope := &QuotedScope{
			symbol: "",
			state:  false,
		}

		line = strings.TrimSpace(line)

		for i, char := range line {
			//? If we encounter a space and we're not inside a quote, finalize the current segment.
			if !quoteScope.state && unicode.IsSpace(char) {
				if current != nil {
					tokens = append(tokens, current)
					current = nil
				}
				continue
			}

			//? Handle quotes (either starting or ending a quoted string).
			if regQuotes.MatchString(string(char)) && (i == 0 || line[i-1] != '\\') {
				//? If we are inside a quote and we encounter the same quote symbol, finalize the current string.
				if quoteScope.state && string(char) == quoteScope.symbol {
					current = append(current, char)
					tokens = append(tokens, current)
					current = nil
					quoteScope.state = false
					continue
				}

				//? If we are outside a quote and have accumulated characters, finalize the current segment.
				if !quoteScope.state && current != nil {
					tokens = append(tokens, current)
					current = nil
				}

				//? Toggle quote state and set the symbol to the current quote type.
				quoteScope.symbol = string(char)
				quoteScope.state = !quoteScope.state
			}

			//? Add the current character to the current segment.
			current = append(current, char)

			//? If it's the last character of the line, finalize the current segment.
			if i == len(line)-1 && current != nil {
				tokens = append(tokens, current)
				current = nil
			}
		}

		if len(tokens) == 0 {
			continue
		}

		instructions = append(instructions, tokens)
	}

	return instructions
}
