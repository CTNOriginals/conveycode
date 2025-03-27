package lexer

import (
	"conveycode/compiler/tokenizer"
	"fmt"
	"slices"
	"strings"

	"github.com/TwiN/go-color"
)

// The lexer works like a curser.
// In its normal state, the start and pos are the same value.
// This is like a carrot cursor sitting in between characters without selecting any.
// When the pos is greater then the start, thats when the cursor turns into a selection
//
// When the start is 0 and the pos is 0, the lexer doesnt have any content selected.
// When the start is 0 and the pos is 1, the lexer is selecting the first token
//
// pos is the index of the next token, this token is not currently selected but will be after the next read
type lexer struct {
	tokens tokenizer.TokenList
	Blocks chan Block
	State  StateFn
	items  []item
	start  int
	pos    int
}

// #region Core
func Lex(tokens tokenizer.TokenList) (lx *lexer) {
	lx = &lexer{
		tokens: tokens,
		Blocks: make(chan Block, 2),
		State:  LexText,
	}

	return lx
}

func (this *lexer) NextBlock() Block {
	for {
		select {
		case block := <-this.Blocks:
			return block
		default:
			//? Nil Pointer: If you explicitly assign nil to a pointer
			//? and then try to dereference it, you’ll encounter the same error.
			//? Ensure that your pointer is assigned a valid memory address before using it.
			//* https://edwinsiby.medium.com/runtime-error-invalid-memory-address-or-nil-pointer-dereference-golang-dd4a58ab7536

			var state StateFn = this.State

			if state == nil {
				this.State = nil
			} else {
				this.State = state(this)
			}
		}
	}
}

//#endregion

func (this *lexer) length() int {
	return this.pos - this.start
}

func (this *lexer) peek() tokenizer.Token {
	if this.isEOF() {
		return tokenizer.Token{
			Typ: tokenizer.EOF,
		}
	}
	return this.tokens[this.pos]
}

// The token that was read most recently
func (this *lexer) current() tokenizer.Token {
	return this.tokens[this.start+max(0, this.length()-1)]
}

func (this *lexer) isEOF() bool {
	return this.pos >= len(this.tokens)
}

// Advances the position forward
func (this *lexer) advance() {
	this.pos++
}

// Reads the current token and advances the position conuming the token it returns
func (this *lexer) read() (ret tokenizer.Token) {
	ret = this.peek()

	if this.isEOF() {
		return
	}

	this.advance()
	return ret
}

// func (this *lexer) peek() tokenizer.Token {
// 	if this.pos+1 >= len(this.tokens) {
// 		return tokenizer.Token{
// 			Typ: tokenizer.EOF,
// 		}
// 	}

// 	return this.tokens[this.pos]
// }

func (this *lexer) peekBack() tokenizer.Token {
	if this.pos == 0 {
		return tokenizer.Token{
			Typ: tokenizer.EOF,
		}
	}

	return this.tokens[this.pos-1]
}

func (this *lexer) stream() (tokens []tokenizer.Token) {
	return this.tokens[this.start:this.pos]
}

func (this *lexer) consume() {
	this.start = this.pos
}
func (this *lexer) reset() {
	this.pos = this.start
}

func (this *lexer) backup() {
	if this.length() <= 0 {
		this.errorf("Can not backup past start")
		return
	}
	this.pos--
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
	if string(this.read().Val) == valid {
		return true
	}

	this.backup()
	return false
}

// Accept until f returns true or EOF is reached
func (this *lexer) acceptUntilFunc(f func(token tokenizer.Token) bool) bool {
	for {
		if f(this.read()) {
			return !this.isEOF()
		}
	}
}

func (this *lexer) expect(valid ...tokenizer.TokenType) (bool, StateFn) {
	var token = this.read()
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
	this.Blocks <- Block{
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
func (this *lexer) getLocationHighlight() (str string) {
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

	if this.length() == 0 {
		stream = slices.Insert(stream, this.start, color.InWhiteOverBlue("ⵊ"))
	}

	for i, part := range stream {
		if i > 0 && stream[i-1] != "\n" {
			if within(i-1) && within(i) {
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

	this.Blocks <- Block{
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
