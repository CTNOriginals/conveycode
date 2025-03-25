package lexer

import (
	"conveycode/compiler/tokenizer"
	"slices"
)

type StateFn func(*lexer) StateFn

var valueTokenTypes = []tokenizer.TokenType{tokenizer.String, tokenizer.Number, tokenizer.Text}
var bracketTokenTypes = []tokenizer.TokenType{
	tokenizer.RoundL,
	tokenizer.SquareL,
	tokenizer.CurlyL,
	tokenizer.RoundR,
	tokenizer.SquareR,
	tokenizer.CurlyR,
}
var openBracketTokenTypes = bracketTokenTypes[:3]
var closeBracketTokenTypes = bracketTokenTypes[3:]

func LexText(lx *lexer) StateFn {
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

func lexAssignment(lx *lexer) (state StateFn) {
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
		return lx.errorf("Invalid value assignment")
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
