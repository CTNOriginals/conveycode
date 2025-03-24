package lexer

import "fmt"

type blockType int

const (
	_ blockType = iota

	Assignment
	Statement
	Method

	BlockError
)

func (this blockType) String() string {
	return [...]string{
		"Assignment",
		"Statement",
		"Method",
		"BlockError",
	}[this-1]
}

type block struct {
	Typ   blockType
	Items []item
}

func (this block) String() (str string) {
	return fmt.Sprintf("%s: %s", this.Typ, this.Items)
}

// func NewBlock(typ blockType) (ret block, ) {

// }

// func (this *block) construct(channel chan item) {

// }
