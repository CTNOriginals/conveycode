package lexer

import "conveycode/compiler/tokenizer"

type StateFn func(*lexer) StateFn

var valueTokenTypes = []tokenizer.TokenType{tokenizer.String, tokenizer.Number, tokenizer.Text}

func LexText(lx *lexer) StateFn {
	var token = lx.next()

	if token.Typ == tokenizer.EOL {
		lx.consume()
		return LexText
	}

	switch string(token.Val) {
	case "var":
		return lexAssignment
	case "if":
		return lexIfStatement
	case "else":
		return lexElseStatement
	}

	return nil
}

func lexAssignment(lx *lexer) StateFn {
	lx.emitItem(Keyword)

	if valid, err := lx.expect(tokenizer.Text); !valid {
		return err
	}

	lx.emitItem(Identifier)

	if valid, err := lx.expect(tokenizer.Operator); !valid {
		return err
	}

	lx.emitItem(Operator)

	if valid, err := lx.expect(valueTokenTypes...); !valid {
		return err
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

	var depth = 0
	if !lx.acceptUntilFunc(func(token tokenizer.Token) bool {
		if token.Typ == tokenizer.RoundR {
			if depth == 0 {
				return true
			}

			depth--
		} else if token.Typ == tokenizer.RoundL {
			depth++
		}

		return false
	}) {
		return lx.errorf("Unmatched bracket for conditional statement")
	}

	lx.emitItem(Condition)

	if valid, err := lx.expect(tokenizer.CurlyL); !valid {
		return err
	}

	depth = 0
	if !lx.acceptUntilFunc(func(token tokenizer.Token) bool {
		if token.Typ == tokenizer.CurlyR {
			if depth == 0 {
				return true
			}

			depth--
		} else if token.Typ == tokenizer.CurlyL {
			depth++
		}

		return false
	}) {
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

	var depth = 0
	if !lx.acceptUntilFunc(func(token tokenizer.Token) bool {
		if token.Typ == tokenizer.CurlyR {
			if depth == 0 {
				return true
			}

			depth--
		} else if token.Typ == tokenizer.CurlyL {
			depth++
		}

		return false
	}) {
		return lx.errorf("Unmatched bracket for statement scope")
	}

	lx.emitItem(Scope)

	lx.emitBlock(Statement)

	return LexText
}
