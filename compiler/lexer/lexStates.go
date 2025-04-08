package lexer

import (
	"conveycode/compiler/syntax"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/utils"
	"slices"
)

type StateFn func(*lexer) StateFn

// fmt.Printf("--%d:%d --\n%s\n", lx.start, lx.pos, lx.getLocationHighlight())

func LexText(lx *lexer) (state StateFn) {
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
			case "func":
				return lexMethod
			}

			if lx.peek().Typ == tokenizer.RoundL {
				return lexCall
			}

			if string(lx.peek().Val) == "=" {
				return lexAssignment
			}
		}
	}

	return nil
}

func lexAssignment(lx *lexer) (state StateFn) {
	defer func() {
		if recover() != nil {
			state = nil
		}
	}()

	if string(lx.current().Val) == "var" {
		lx.emitItem(Keyword)
	} else {
		lx.reset()
	}

	lx.expect(tokenizer.Text)
	lx.emitItem(Identifier)

	lx.expect(tokenizer.Operator)
	lx.emitItem(Operator)

	if !lx.acceptUntilFunc(func(token tokenizer.Token) bool {
		var valueContent = append(tokenizer.ValueTokenTypes, tokenizer.Operator)

		if slices.Contains(valueContent, token.Typ) {
			return false
		}

		if token.Typ == tokenizer.RoundL {
			lx.wrapScope()
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

func lexIfStatement(lx *lexer) (state StateFn) {
	defer func() {
		if recover() != nil {
			state = nil
		}
	}()

	lx.emitItem(Keyword)

	lx.expect(tokenizer.RoundL)
	lx.wrapScope()
	lx.emitItem(Condition)

	lx.expect(tokenizer.CurlyL)
	lx.wrapScope()
	lx.emitItem(Scope)

	if lx.accept(tokenizer.EOL) {
		lx.consume()
	}
	if lx.acceptContent("else") {
		return lexElseStatement
	}

	lx.emitBlock(Statement)

	return LexText
}

func lexElseStatement(lx *lexer) (state StateFn) {
	defer func() {
		if recover() != nil {
			state = nil
		}
	}()

	if lx.acceptContent("if") {
		return lexIfStatement
	}

	lx.emitItem(Keyword)

	lx.expect(tokenizer.CurlyL)
	lx.wrapScope()
	lx.emitItem(Scope)

	lx.emitBlock(Statement)

	return LexText
}

func lexCommand(lx *lexer) (state StateFn) {
	defer func() {
		if recover() != nil {
			state = nil
		}
	}()

	if !slices.Contains(utils.Keys(syntax.Commands), string(lx.peekBack().Val)) {
		return lx.errorf("Unknown command: %s", string(lx.peekBack().Val))
	}

	lx.emitItem(Command)

	lx.expect(tokenizer.RoundL)
	lx.wrapScope()
	lx.emitItem(Arguments)

	lx.emitBlock(BuiltIn)

	return LexText
}

func lexMethod(lx *lexer) (state StateFn) {
	defer func() {
		if recover() != nil {
			state = nil
		}
	}()

	lx.emitItem(Keyword)

	lx.expect(tokenizer.Text)
	lx.emitItem(Identifier)

	lx.expect(tokenizer.RoundL)
	lx.wrapScope()
	lx.emitItem(Arguments)

	lx.expect(tokenizer.CurlyL)
	lx.wrapScope()
	lx.emitItem(Scope)

	lx.emitBlock(Method)

	return LexText
}

func lexCall(lx *lexer) (state StateFn) {
	defer func() {
		if recover() != nil {
			state = nil
		}
	}()

	lx.emitItem(Identifier)

	lx.expect(tokenizer.RoundL)
	lx.wrapScope()
	lx.emitItem(Arguments)

	lx.emitBlock(Call)

	return LexText
}
