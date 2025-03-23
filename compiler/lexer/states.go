package lexer

import (
	"unicode"
)

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var digits = "0123456789"
var hexadecimal = digits + "abcdefABCDEF"
var wordBoundry = alphabet + digits + "_"
var quotes = "\"'`"

func LexText(lx *lexer) StateFn {
	for {
		if lx.isEOF() {
			break
		}

		if lx.char() == '\n' {
			return LexEOL
		}

		if lx.containsChar(quotes) {
			if lx.length() > 0 {
				lx.backup()
				lx.emit(Text)
				lx.next()
			}
			return LexString
		}

		if lx.containsChar(digits) || (lx.containsChar("+-") && unicode.IsDigit(lx.peek())) {
			if lx.length() > 0 {
				lx.backup()
				lx.emit(Text)
				lx.next()
			}
			return LexNumber
		}

		if unicode.IsSpace(lx.peek()) {
			lx.next()
			switch string(lx.stream()) {
			case "var":
				return LexVariable
			case "if", "else":
				return LexStatement
			}

			if lx.length() > 0 {
				lx.emit(Text)
			}

			lx.consume()
		}

		lx.next()
	}

	if lx.length() > 0 {
		lx.emit(Text)
	}
	lx.emit(EOF)

	return nil
}

func LexVariable(lx *lexer) StateFn {
	lx.emit(Keyword) //? Emit "var"
	lx.skipRun(" ")

	if !lx.accept(alphabet) { //? Accept the first letter to be an alpha character, not anything else
		return lx.errorf("Unexpected starting character of identifier: '%s'", string(lx.char()))
	}
	lx.acceptRun(wordBoundry) //? accept all of the following word defining characters that makes up the variable identifier

	lx.emit(Identifier)

	lx.skipRun(" ")

	if lx.accept("=") {
		lx.emit(Operator)
	}

	lx.skipRun(" ") //? Remove any trailing spaces

	return LexText
}

func LexString(lx *lexer) StateFn {
	lx.next() //? Pass over the opening quote

	if lx.acceptUntil(string(lx.startChar())) {
		if lx.peekBack() == '\\' {
			return LexString //? Start from the current position again as the quote is escaped
		}
	} else {
		return lx.errorf("Unmatched quote (%s)", string(lx.startChar()))
	}

	lx.next() //? Include the closing quote
	lx.emit(String)

	return LexText
}
func LexNumber(lx *lexer) StateFn {
	var valid = digits
	var isHex = false

	lx.skipRun(" ") //? TrimStart
	lx.accept("+-") //? Optional Sign

	//? Hexedecimal
	if lx.accept("0") && lx.accept("xX") {
		valid = hexadecimal
		isHex = true
	}

	lx.acceptRun(valid)

	if !isHex && lx.accept(".") { //? Floating point
		lx.acceptRun(valid)
	}

	if !isHex && lx.accept("eE") {
		lx.accept("+-")
		if !lx.containsChar(digits) {
			lx.errorf("Bad number syntax: %q", lx.stream())
		}
		lx.acceptRun(digits)
	}

	lx.emit(Number)

	return LexText
}

func LexEOL(lx *lexer) StateFn {
	lx.next()    //? Step over the newline
	lx.consume() //? Ignore it, it should not be included in the item value
	lx.emit(EOL)
	return LexText
}

func LexStatement(lx *lexer) StateFn {
	var statement = string(lx.stream())

	lx.emit(Keyword) //? Emit the keyword

	lx.skipRun(" ")

	switch statement {
	case "if":
		return LexCondition
	case "else":
		if lx.accept("if") {
			return LexStatement
		}
		return LexScope
	}

	return lx.errorf("Bad statement syntax")
}

func LexCondition(lx *lexer) StateFn {
	return LexText
}

func LexScope(lx *lexer) StateFn {
	return LexText
}
