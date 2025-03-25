package lexer

import (
	"conveycode/compiler/tokenizer"
	"fmt"
	"slices"
	"strings"

	"github.com/TwiN/go-color"
)

type lexer struct {
	Blocks chan block
	tokens tokenizer.TokenList
	items  []item
	start  int
	pos    int
}

func Lex(tokens tokenizer.TokenList) (ret lexer) {
	ret = lexer{
		tokens: tokens,
		Blocks: make(chan block),
	}

	go ret.run()

	return ret
}

func (this *lexer) run() {
	var state StateFn
	for state = LexText; state != nil; {
		state = state(this)
	}

	close(this.Blocks)

	//* Cleanup
	this.tokens = nil
	this.items = nil
}

func (this *lexer) token() tokenizer.Token {
	return this.tokens[this.pos]
}

func (this *lexer) length() int {
	return this.pos - this.start
}

func (this *lexer) isEOF() bool {
	return this.token().Typ == tokenizer.EOF
}

func (this *lexer) next() (ret tokenizer.Token) {
	ret = this.token()

	if this.isEOF() {
		return
	}

	this.pos++
	return ret
}

func (this *lexer) backup() {
	if this.length() <= 0 {
		this.errorf("Can not backup past start")
		return
	}
	this.pos--
}

func (this *lexer) peek() (tok tokenizer.Token) {
	this.next()
	tok = this.token()
	this.backup()
	return
}

func (this *lexer) peekBack() (tok tokenizer.Token) {
	if this.pos == 0 {
		return tokenizer.Token{
			Typ: tokenizer.EOF,
		}
	}

	return this.tokens[this.pos-1]
}

func (this *lexer) consume() {
	this.start = this.pos
}

// func (this *lexer) accept(valid ...tokenizer.TokenType) bool {
// 	if slices.Contains(valid, this.next().Typ) {
// 		return true
// 	}

// 	this.backup()
// 	return false
// }

// func (this *lexer) acceptRun(valid ...tokenizer.TokenType) {
// 	for this.accept(valid...) {
// 	}
// }

func (this *lexer) acceptContent(valid string) bool {
	if string(this.next().Val) == valid {
		return true
	}

	this.backup()
	return false
}

// Accept until f returns true
func (this *lexer) acceptUntilFunc(f func(token tokenizer.Token) bool) bool {
	for {
		if this.isEOF() {
			return false
		}

		if f(this.next()) {
			return true
		}
	}
}

func (this *lexer) expect(valid ...tokenizer.TokenType) (bool, StateFn) {
	var token = this.next()
	if !slices.Contains(valid, token.Typ) {
		var validString []string = make([]string, len(valid))

		for i, typ := range valid {
			validString[i] = typ.String()
		}

		return false, this.errorf("expected type [%s] but found %s", strings.Join(validString, ", "), token.String())
	}

	return true, nil
}

func (this *lexer) expectSequence(seq [][]tokenizer.TokenType) (bool, StateFn) {
	for _, valid := range seq {
		if valid, err := this.expect(valid...); !valid {
			return false, err
		}
	}

	return true, nil
}

func (this *lexer) stream() (tokens []tokenizer.Token) {
	tokens = make([]tokenizer.Token, this.length())
	copy(tokens, this.tokens[this.start:this.pos])
	return tokens
}

func (this *lexer) wrapScope() bool {
	var openBracket = this.peekBack().Typ
	var closeBracket tokenizer.TokenType

	switch openBracket {
	case tokenizer.RoundL:
		closeBracket = tokenizer.RoundR
	case tokenizer.SquareL:
		closeBracket = tokenizer.SquareR
	case tokenizer.CurlyL:
		closeBracket = tokenizer.CurlyR
	}

	var depth = 0
	return this.acceptUntilFunc(func(token tokenizer.Token) bool {
		if token.Typ == closeBracket {
			if depth == 0 {
				return true
			}

			depth--
		} else if token.Typ == openBracket {
			depth++
		}

		return false
	})
}

func (this *lexer) emitItem(typ itemType) {
	this.items = append(this.items, NewItem(typ, this.stream()...))
	this.consume()
}
func (this *lexer) emitBlock(typ blockType) {
	this.Blocks <- block{
		Typ:   typ,
		Items: this.items,
	}

	//? Reset the arrays length to 0, freeing up the slices contents
	//? this doesnt make a new slice,
	//? it just keeps the values there and marks those memory adresses as free to override
	this.items = this.items[:0]
	this.consume()
}

//#region Utils

// Returns all the values of the tokens in sequence
// highlighting the section between the current start and pos of the lexer
func (this *lexer) getLocationHighlight() string {
	var within = func(index int) bool {
		return index >= this.start && index < this.pos
	}

	var stream []string = make([]string, len(this.tokens))
	for i, token := range this.tokens {
		if token.Typ == tokenizer.EOL {
			stream[i] = "\n"
			continue
		}
		stream[i] = token.ColoredValue()
	}

	var str string = ""
	for i, part := range stream {
		if i > 0 && stream[i-1] != "\n" {
			if within(i-1) && within(i) {
				str += color.InBlackOverWhite(" ")
			} else {
				str += " "
			}
		}

		if within(i) || (this.length() == 0 && i == this.pos) {
			str += color.InUnderline(color.WhiteBackground + part)
		} else {
			str += color.InUnderline(part)
		}
	}

	return str
}

func (this *lexer) errorf(format string, args ...any) StateFn {
	var message = fmt.Sprintf(color.InRed(format), args...)
	message += fmt.Sprintf("\n-- %s %d:%d (%d) --\n%s", color.InYellow("LOCATION"), this.start, this.pos, this.length(), this.getLocationHighlight())

	this.Blocks <- block{
		Typ: BlockError,
	}

	fmt.Println(message)

	return nil
}

//#endregion

//() Item constructor
//| Select items over channel
//- Normal item
//< Add it to the current item stream
//- End of Transmission item
//< Send the current item stream over the the blocks channel

// The starting token contains "var"
// Make a new block and pass the type "Variable"
// The block calls the item constructor to create an item with type "Keyword" and puts in the "var" token
// The item returns
// The block now calls the item constructor for the identifier
// After that its for the operator
// and finally for the value
// After that, the block is sent back over the channel to the lexer
