package tools

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

func (this *Cursor) String() string {
	return fmt.Sprintf("Content Length: %d\nPosition: %d\nLine: %d\nEOF: %t", len(this.Content), this.Pos, this.Line, this.EOF)
}

// Seek the cursors position relative to its current position.
//
// If the offset is out of range, the cursors position will remain the same and the function returns false
func (this *Cursor) Seek(offset uint) bool {
	if err := this.validateOffset(int(offset)); err != nil {
		return false
	}

	if this.EOF {
		this.EOF = false
	}

	for range offset {
		if this.Content[this.Pos] == '\n' {
			this.Column = 1
			this.Line++
		} else {
			this.Column++
		}

		this.Pos++
	}

	return true
}

// Consumes the current character and seeks next
func (this *Cursor) Consume() bool {
	return this.Seek(1)
}

// Peek at an offset relative to the cursors current position
func (this *Cursor) PeekOffset(offset int) rune {
	if err := this.validateOffset(offset); err != nil {
		return types.EOT
	}

	return this.Content[this.Pos+offset]
}

// Returns the character at the current position of the cursor
func (this *Cursor) Peek() rune {
	return this.PeekOffset(0)
}

// Returns the character at the position ahead of the cursor
func (this *Cursor) PeekNext() rune {
	return this.PeekOffset(1)
}

// Returns the character at the position ahead of the cursor
func (this *Cursor) PeekPrev() rune {
	return this.PeekOffset(-1)
}

// Returns the current character and consumes it
func (this *Cursor) Read() (char rune) {
	if this.EOF {
		return types.EOT
	}

	char = this.Peek()

	if !this.Consume() {
		this.EOF = true
	}

	return
}

// Read n of characters from the current cursor position and returns all of them
func (this *Cursor) ReadN(n int) (list []rune) {
	for range n {
		list = append(list, this.Read())
	}
	return
}

// Reads until f(c) returns true or the EOF is reached
func (this *Cursor) ReadUntilFunc(f func(c rune) bool) (list []rune) {
	var char = this.Read()

	for !f(char) && !this.EOF {
		list = append(list, char)
		char = this.Read()
	}

	return
}

//#region Private

func (this *Cursor) validateOffset(offset ...int) error {
	for i := range offset {
		var index = this.Pos + offset[i]

		if index < 0 || index >= len(this.Content) {
			// fmt.Printf("%d: Offset out of range (pos: %d + off: %d[%d] = %d)\n", this.Pos, this.Pos, offset[i], i, index)
			return io.EOF
		}
	}

	return nil
}

func (this *Cursor) getColumn() (col int) {
	col = 0
	char := this.Content[this.Pos] //? Cant use seek because seek calls this

	for char != '\n' && (this.Pos-(col+1)) > 0 { //? line feed
		col++
		char = this.Content[this.Pos-col]
	}

	return col
}

// func (this *Cursor) throw(err error) {
// 	fmt.Println(err.Error())
// }

//#endregion
