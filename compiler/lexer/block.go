package lexer

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

type BlockType int

const (
	_ BlockType = iota

	BlockText

	Assignment
	Statement
	Method
	BuiltIn

	BlockEOF
	BlockError
)

func (this BlockType) String() string {
	return [...]string{
		"BlockText",
		"Assignment",
		"Statement",
		"Method",
		"BuiltIn",
		"BlockEOF",
		"BlockError",
	}[this-1]
}

type Block struct {
	Typ   BlockType
	Items []item
}

func NewBlock(typ BlockType) Block {
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
