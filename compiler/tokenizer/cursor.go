package tokenizer

import (
	"conveycode/compiler/types"
	"fmt"
	"io"
	"slices"
)

func CursorTests(content []rune) {
	var cursor = NewCursor(content)

	fmt.Println(len(cursor.Content))

	// fmt.Printf("%d: '%s'\n", cursor.Pos, string(cursor.Peek()))
	// cursor.Seek(39)
	// fmt.Printf("%d/%d: '%s'\n", cursor.Pos-12, cursor.Pos, string(cursor.ReadN(12)))
	// fmt.Printf("%d: '%v'\n", cursor.Pos-1, cursor.Read())
	// fmt.Printf("%d: '%v'\n", cursor.Pos-1, cursor.Read())

	// cursor.Pos = 0
	// cursor.Seek(115)
	// fmt.Printf("%d/%d: '%s'\n", cursor.Pos-20, cursor.Pos, string(cursor.ReadN(20)))

	// cursor.Seek(55)
	// cursor.Seek(1) // Go past the first quote
	// fmt.Println(string(cursor.ReadUntil('"')))

	// cursor.Seek(18)
	// fmt.Println(string(cursor.ReadUntilFunc(func(c rune) bool {
	// 	return !unicode.IsDigit(c)
	// })))

	// cursor.Seek(24)
	// fmt.Println(string(cursor.Peek()), cursor.getColumn())

	fmt.Println(cursor)
}

type Cursor struct {
	// The Content that the cursor is running though
	Content []rune

	// The position of the cursor relative to the full content
	Pos int

	// The line the cursor is currently on
	Line int

	// The column number the cursor is currently on
	Column int

	// Wether the end of the file has been reached
	//
	// Used by Cursor.Read() to know when to return EOT.
	// Set to true by Cursor.Read() once Cursor.Consume() returns false.
	// Set to false once Cursor.Seek() is used with a negative number.
	EOF bool
}

func NewCursor(content []rune) *Cursor {
	content = slices.DeleteFunc(content, func(e rune) bool { return e == '\r' })

	return &Cursor{
		Content: content,
		Pos:     0,
		Line:    1,
		Column:  1,
		EOF:     false,
	}
}

func (cur *Cursor) String() string {
	return fmt.Sprintf("Content Length: %d\nPosition: %d\nLine: %d\nEOF: %t", len(cur.Content), cur.Pos, cur.Line, cur.EOF)
}

// Seek the cursors position relative to its current position.
//
// If the offset is out of range, the cursors position will remain the same and the function returns false
func (cur *Cursor) Seek(offset uint) bool {
	if err := cur.validateOffset(int(offset)); err != nil {
		return false
	}

	for range offset {
		if cur.Content[cur.Pos] == '\n' {
			cur.Column = 1
			cur.Line++
		} else {
			cur.Column++
		}

		cur.Pos++
	}

	return true
}

// Peek at an offset relative to the cursors current position
func (cur *Cursor) PeekOffset(offset int) rune {
	if err := cur.validateOffset(offset); err != nil {
		return types.EOT
	}

	return cur.Content[cur.Pos+offset]
}

// Returns the character at the current position of the cursor
func (cur *Cursor) Peek() rune {
	return cur.PeekOffset(0)
}

// Returns the character at the position ahead of the cursor
func (cur *Cursor) PeekNext() rune {
	return cur.PeekOffset(1)
}

// Returns the character at the position ahead of the cursor
func (cur *Cursor) PeekPrev() rune {
	return cur.PeekOffset(-1)
}

// Returns the current character and consumes it
func (cur *Cursor) Read() (char rune) {
	if cur.EOF {
		return types.EOT
	}

	char = cur.Peek()

	if !cur.Seek(1) {
		cur.EOF = true
	}

	return
}

// Read n of characters from the current cursor position and returns all of them
func (cur *Cursor) ReadN(n int) (list []rune) {
	for range n {
		list = append(list, cur.Read())
	}
	return
}

// Reads until f(c) returns true or the EOF is reached
func (cur *Cursor) ReadUntilFunc(f func(c rune) bool) (list []rune) {
	for !f(cur.Peek()) && !cur.EOF {
		list = append(list, cur.Read())
	}

	return
}

//#region Private

func (cur *Cursor) validateOffset(offset ...int) error {
	for i := range offset {
		var index = cur.Pos + offset[i]

		if index < 0 || index >= len(cur.Content) {
			// fmt.Printf("%d: Offset out of range (pos: %d + off: %d[%d] = %d)\n", cur.Pos, cur.Pos, offset[i], i, index)
			return io.EOF
		}
	}

	return nil
}

//#endregion
