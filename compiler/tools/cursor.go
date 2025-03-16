package tools

import (
	"conveycode/compiler/utils"
	"fmt"
	"io"
)

func CursorTests(content []rune) {
	var cursor = NewCursor(content)

	fmt.Printf("%d: '%s'\n", cursor.Pos, string(cursor.Peak(0)))
	cursor.Seek(41)
	fmt.Printf("%d: '%s'\n", cursor.Pos, string(cursor.Peak(0)))
	fmt.Printf("%d: '%s'\n", cursor.Pos-1, string(cursor.Read()))
	fmt.Printf("%d: '%s'\n", cursor.Pos, string(cursor.PeakRange(-2, 3)))
	fmt.Println(cursor)
}

type Cursor struct {
	// The content that the cursor is running though
	content []rune

	// The position of the cursor relative to the full content
	Pos int

	// The line the cursor is currently on
	Line int
}

func NewCursor(content []rune) *Cursor {
	for i, c := range content {
		if c == '\r' && content[i+1] == '\n' {
			content[i] = '\n'
			content, _ = utils.Splice(content, i, 1)
		}
	}

	return &Cursor{
		content: content,
		Pos:     0,
		Line:    1,
	}
}

func (this *Cursor) String() string {
	return fmt.Sprintf("Content Length: %d\nPosition: %d\nLine: %d", len(this.content), this.Pos, this.Line)
}

func (this *Cursor) validateOffset(offset ...int) error {
	for i := range offset {
		var index = this.Pos + offset[i]
		// fmt.Println(offset[i], index)

		if index < 0 || index >= len(this.content) {
			return io.EOF
		}
	}

	return nil
}

func (this *Cursor) throw(err error) {
	fmt.Println(err.Error())
}

// Gets the character at the relative offset of the current cursor position and returns it
//
// If the offset is out of range, either the first or the last character is returned instead
// depending on if the offset was positive or negative
func (this *Cursor) Peak(offset int) rune {
	if err := this.validateOffset(offset); err != nil {
		this.throw(err)

		index := this.Pos

		if offset > this.Pos {
			index = len(this.content) - 1
		} else if offset < this.Pos {
			index = 0
		}

		return this.content[index]
	}

	return this.content[this.Pos+offset]
}

// Get all characters within the range of the start and end index relative to the current cursor position
func (this *Cursor) PeakRange(start int, end int) []rune {
	if err := this.validateOffset(start); err != nil {
		this.throw(err)
		start = this.Pos
	}

	if err := this.validateOffset(end); err != nil {
		this.throw(err)
		end = this.Pos
	}

	return this.content[(this.Pos + start):(this.Pos + end)]
}

// Consume the current character and return it
func (this *Cursor) Read() (char rune) {
	char = this.Peak(0)
	this.Seek(1)
	return
}

// Seek the cursors position relative to its current position.
//
// If the offset is out of range, the cursors position will remain the same and the function returns false
func (this *Cursor) Seek(offset int) bool {
	if err := this.validateOffset(offset); err != nil {
		this.throw(err)
		return false
	}

	absOffset, signed := utils.Abs(offset)
	direction := utils.If(signed, -1, 1)

	for range absOffset {
		if this.content[this.Pos] == '\n' {
			this.Line += direction
		}

		this.Pos += direction
	}

	return true
}
