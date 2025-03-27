package lexer

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

type blockType int

const (
	_ blockType = iota

	BlockText

	Assignment
	Statement
	Method
	BuiltIn

	BlockError
)

func (this blockType) String() string {
	return [...]string{
		"BlockText",
		"Assignment",
		"Statement",
		"Method",
		"BuiltIn",
		"BlockError",
	}[this-1]
}

type Block struct {
	Typ   blockType
	Items []item
}

func NewBlock(typ blockType) Block {
	return Block{
		Typ: typ,
	}
}

func (this Block) String() (str string) {
	var itemString = make([]string, len(this.Items))
	for i, item := range this.Items {
		itemString[i] = item.String()
	}
	return fmt.Sprintf("%s\n %s", color.InBlue(color.Bold+this.Typ.String()), strings.Join(itemString, " "))
}

// func NewBlock(typ blockType) (ret block, ) {

// }

// func (this *block) construct(channel chan item) {

// }
