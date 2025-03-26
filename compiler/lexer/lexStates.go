package lexer

import (
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/types"
	"slices"
)

type StateFn func(*lexer) StateFn

var valueTokenTypes = []tokenizer.TokenType{tokenizer.String, tokenizer.Number, tokenizer.Text}

// var bracketTokenTypes = []tokenizer.TokenType{
// 	tokenizer.RoundL,
// 	tokenizer.SquareL,
// 	tokenizer.CurlyL,
// 	tokenizer.RoundR,
// 	tokenizer.SquareR,
// 	tokenizer.CurlyR,
// }
// var openBracketTokenTypes = bracketTokenTypes[:3]
// var closeBracketTokenTypes = bracketTokenTypes[3:]

// fmt.Printf("--%d:%d --\n%s\n", lx.start, lx.pos, lx.getLocationHighlight())

func LexText(lx *lexer) StateFn {

	for {
		var token = lx.read()

		if token.Typ == tokenizer.EOF {
			break
		}

		switch token.Typ {
		case tokenizer.EOL:
			lx.consume()
		case tokenizer.Command:
			return lexCommand
		case tokenizer.Text:
			switch string(token.Val) {
			case "var":
				return lexAssignment
			case "if":
				return lexIfStatement
			case "else":
				return lexElseStatement
			}

			if string(lx.peek().Val) == "=" {
				return lexAssignment
			}
		}
	}

	return nil
}

func lexAssignment(lx *lexer) (state StateFn) {
	if string(lx.current().Val) == "var" {
		lx.emitItem(Keyword)
	} else {
		lx.reset()
	}

	if valid, err := lx.expect(tokenizer.Text); !valid {
		return err
	}
	lx.emitItem(Identifier)

	if valid, err := lx.expect(tokenizer.Operator); !valid {
		return err
	}

	lx.emitItem(Operator)

	if !lx.acceptUntilFunc(func(token tokenizer.Token) bool {
		var valueContent = append(valueTokenTypes, tokenizer.Operator)

		if slices.Contains(valueContent, token.Typ) {
			return false
		}

		if token.Typ == tokenizer.RoundL {
			if !lx.wrapScope() {
				lx.errorf("Value assignment contained unmatched bracket")
				lx.consume()
				return true
			}

			return false
		}

		if token.Typ == tokenizer.EOL {
			lx.backup()
			return true
		}

		return true
	}) {
		lx.backup() //? Dont include the EOF
	}

	lx.emitItem(Value)

	lx.emitBlock(Assignment)

	return LexText
}

func lexIfStatement(lx *lexer) StateFn {
	lx.emitItem(Keyword)

	if valid, err := lx.expect(tokenizer.RoundL); !valid {
		return err
	}

	if !lx.wrapScope() {
		return lx.errorf("Unmatched bracket for conditional statement")
	}

	lx.emitItem(Condition)

	if valid, err := lx.expect(tokenizer.CurlyL); !valid {
		return err
	}

	if !lx.wrapScope() {
		return lx.errorf("Unmatched bracket for statement scope")
	}

	lx.emitItem(Scope)

	lx.emitBlock(Statement)

	return LexText
}

func lexElseStatement(lx *lexer) StateFn {
	if lx.acceptContent("if") {
		return lexIfStatement
	}

	lx.emitItem(Keyword)

	if valid, err := lx.expect(tokenizer.CurlyL); !valid {
		return err
	}

	if !lx.wrapScope() {
		return lx.errorf("Unmatched bracket for statement scope")
	}

	lx.emitItem(Scope)

	lx.emitBlock(Statement)

	return LexText
}

func lexCommand(lx *lexer) StateFn {
	if !slices.Contains(types.Commands, string(lx.peekBack().Val)) {
		return lx.errorf("Unknown command: %s", string(lx.peekBack().Val))
	}

	lx.emitItem(Command)

	if valid, err := lx.expect(tokenizer.RoundL); !valid {
		return err
	}

	if !lx.wrapScope() {
		return lx.errorf("Unmatched bracket")
	}

	lx.emitItem(Arguments)

	lx.emitBlock(Instruction)

	return LexText
}
