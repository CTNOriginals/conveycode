package lexer

import (
	"conveycode/compiler/tokenizer"
	"fmt"

	"github.com/TwiN/go-color"
)

type lexer struct {
	Blocks chan block
	tokens tokenizer.TokenList
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
	this.start = 16
	this.pos = this.start + 18
	// fmt.Println(this.getLocationHighlight())

	this.errorf("Test error")

	var state StateFn
	for state = LexText; state != nil; {
		state = state(this)
	}

	close(this.Blocks)
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
	return
}

func (this *lexer) backup() {
	if this.length() <= 0 {
		//! ERROR
		return
	}
	this.pos--
}

func (this *lexer) peek() (ret tokenizer.Token) {
	this.next()
	ret = this.token()
	this.backup()
	return
}

func (this *lexer) consume() {
	this.start = this.pos
}

// Returns all the values of the tokens in sequence
// highlighting the section between the current start and pos of the lexer
func (this *lexer) getLocationHighlight() string {
	var within = func(index int) bool {
		return index >= this.start && index <= this.pos
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
			if within(i) {
				str += color.InBlackOverWhite(" ")
			} else {
				str += " "
			}
		}

		if within(i) {
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
