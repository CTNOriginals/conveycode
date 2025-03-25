package lexer

import (
	"conveycode/compiler/tokenizer"
)

type StateFn func(*lexer) StateFn

var valueTokenTypes = []tokenizer.TokenType{tokenizer.String, tokenizer.Number, tokenizer.Text}

func LexText(lx *lexer) StateFn {
	// var redirect = func(state StateFn) StateFn {
	// 	if lx.length() > 1 {
	// 		lx.emitItem(ItemText)
	// 		lx.emitBlock(BlockText)
	// 	}

	// 	return state
	// }

	for {
		var token = lx.next()

		if lx.isEOF() {
			break
		}

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
	}

	return nil
}

func lexAssignment(lx *lexer) StateFn {
	defer func() {

	}()

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
