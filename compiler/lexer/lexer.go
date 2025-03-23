package lexer

import (
	"conveycode/compiler/types"
	"fmt"
	"slices"
	"strings"

	"github.com/TwiN/go-color"
)

type lexer struct {
	content []rune
	start   int
	pos     int
	Items   chan item
}

type StateFn func(*lexer) StateFn

func Lex(content []rune) (lx lexer) {
	content = slices.DeleteFunc(content, func(e rune) bool { return e == '\r' })
	lx = lexer{
		content: content,
		Items:   make(chan item),
	}

	go lx.run()

	return lx
}

func (lx *lexer) run() {
	var state StateFn
	for state = LexText; state != nil; {
		state = state(lx)
	}
	close(lx.Items)
}

//#region Getters

// Returns the current character
func (lx *lexer) char() rune {
	return lx.content[lx.pos]
}

// Returns wether the current char is present in the haystack
func (lx *lexer) containsChar(haystack string) bool {
	return strings.ContainsRune(haystack, lx.char())
}

// Returns the character at the start position
func (lx *lexer) startChar() rune {
	return lx.content[lx.start]
}

// Returns the number of characters from start to pos
func (lx *lexer) length() int {
	return lx.pos - lx.start
}

// Returns wether the end of file has been reached
func (lx *lexer) isEOF() bool {
	return lx.pos >= len(lx.content)
}

// Returns the current character and andvances to the next
func (lx *lexer) next() (char rune) {
	if lx.isEOF() {
		return types.EOT
	}

	char = lx.char()
	lx.pos++
	return char
}

// Returns the next character from the current pos and without moving the position
func (lx *lexer) peek() rune {
	if lx.pos+1 >= len(lx.content) {
		return types.EOT
	}
	return lx.content[lx.pos+1]
}

// Returns the previous character from the current pos and without moving the position
func (lx *lexer) peekBack() rune {
	if lx.pos == 0 {
		return types.EOT
	}

	return lx.content[lx.pos-1]
}

// Returns the current stream of characters between the start and pos
func (lx *lexer) stream() []rune {
	return lx.content[lx.start:lx.pos]
}

//#endregion

//#region Actions

// Consumes the current location and resets the start to the current position
func (lx *lexer) consume() {
	lx.start = lx.pos
}

// Undo the most recent lx.next() call
func (lx *lexer) backup() {
	lx.pos -= 1
}

// Reset the position back to the start
func (lx *lexer) reset() {
	lx.pos = lx.start
}

// Emit the current stream of characters and send the through the channel
func (lx *lexer) emit(typ itemType) {
	lx.Items <- item{Typ: typ, Val: lx.stream()}
	lx.consume()
}

// Skip the current character and advance forward
func (lx *lexer) skip() {
	lx.next()
	lx.consume()
}

// Skips until the next character is not presit in the valid string
func (lx *lexer) skipRun(valid string) {
	lx.acceptRun(valid)
	lx.consume()
}

// Advances forward if the current character is present within the valid string and returns true.
// If the current character is not present within the valid string, it doesnt advance and returns false instead
func (lx *lexer) accept(valid string) bool {
	for strings.ContainsRune(valid, lx.next()) {
		return true
	}

	//? backup one to undo the invalid character read
	lx.backup()
	return false
}

// Accepts as many characters as it can until it hits an invalid character that is not present within the valid string
//
// Returns true if the the f function broke the loop, false if the EOF was reached
func (lx *lexer) acceptRun(valid string) bool {
	//? Loop until the the current character is not valid
	for {
		if lx.isEOF() {
			return false
		}
		if !lx.accept(valid) {
			return true
		}
	}
}

// Accepts as many characters as it can until it hits a character that is present within the valid string
//
// Returns true if the the f function broke the loop, false if the EOF was reached
func (lx *lexer) acceptUntil(valid string) bool {
	for {
		if lx.isEOF() {
			return false
		}

		if lx.accept(valid) {
			lx.backup()
			return true
		}
		lx.next()
	}
}

// Accepts all characters until the f function returns false
//
// Returns true if the the f function broke the loop, false if the EOF was reached
func (lx *lexer) acceptFunc(f func(char rune) bool) bool {
	for {
		if lx.isEOF() {
			return false
		}
		if f(lx.next()) {
			return true
		}
	}
}

// Test if all the characters in the current stream pass true if put through the f function
func (lx *lexer) test(f func(char rune) bool) bool {
	for _, char := range lx.stream() {
		if !f(char) {
			return false
		}
	}

	return true
}

//#endregion

//#region Logging

func (lx *lexer) getLocationHighlight() string {
	var content = strings.Split(string(lx.content), "")

	var posIndex = min(lx.pos, len(lx.content)-1) //? To prevent out of range index errors

	//? Makes the character that caused the error show up with a white background
	content[lx.start] = color.White + color.BlueBackground + content[lx.start]
	content[posIndex] = content[posIndex] + color.Reset

	return strings.Join(content, "")
}

func (lx *lexer) errorf(format string, args ...any) StateFn {
	var message = fmt.Sprintf(color.InRed(format), args...)
	message += fmt.Sprintf("\n-- %s %d:%d (%d) --\n%s", color.InYellow("LOCATION"), lx.start, lx.pos, lx.length(), lx.getLocationHighlight())

	var characters = make([]rune, len(message))

	for i, char := range message {
		characters[i] = char
	}

	lx.Items <- item{
		Typ: Error,
		Val: characters,
	}

	return nil
}

//#endregion
