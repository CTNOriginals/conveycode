package lexer

import (
	"conveycode/compiler/tokenizer"
	"conveycode/internal"
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
	Call
	Return
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
		"Call",
		"Return",
		"BuiltIn",
		"BlockEOF",
		"BlockError",
	}[this-1]
}

type Block struct {
	Typ   BlockType
	Items []Item
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

func (this Block) IsError() bool {
	return this.Typ == BlockError
}

// Finds and returns the first item of typ
//
// Returns new item of type ItemError if the typ was not present within this block
func (this Block) FindItemByType(typ ItemType) Item {
	for _, item := range this.Items {
		if item.Typ == typ {
			return item
		}
	}

	return NewItem(ItemError)
}

// Returns the first identifier found in the block if present
func (this Block) GetIdentifier() (item string) {
	if this.FindItemByType(Identifier).IsError() {
		return internal.StringError
	}
	return this.FindItemByType(Identifier).ValueString()
}

func (this Block) GetItemsOfType(typ ItemType) (items []Item) {
	for _, item := range this.Items {
		if item.Typ == typ {
			items = append(items, item)
		}
	}

	return items
}

func (this Block) getArgumentHolderType() ItemType {
	switch this.Typ {
	case Method:
		return Parameters
	case Call:
		return Arguments
	}

	return ItemError
}

// Get the number of parameters/arguments this block holds.
// Can be both a function call and a function definition
func (this Block) GetArity() (count int) {
	var typ = this.getArgumentHolderType()
	var fieldItem = this.FindItemByType(typ)

	if fieldItem.IsError() || fieldItem.Tokens[1].Typ == tokenizer.RoundR {
		return 0
	}

	for _, token := range fieldItem.Tokens {
		if token.Typ.IsValue() {
			count++
		}
	}

	return count
}

func (this Block) GetArguments() (args []string) {
	var typ = this.getArgumentHolderType()
	var argItem = this.FindItemByType(typ)

	for _, token := range argItem.Tokens {
		if token.Typ.IsValue() {
			args = append(args, string(token.Val))
		}
	}

	return args
}

func (this Block) FirstItem() Item {
	if len(this.Items) == 0 {
		return NewItem(ItemError)
	}

	return this.Items[0]
}

func (this Block) BlockFile() string {
	return this.FirstItem().ItemFile()
}

// Returns the line number this block starts on
func (this Block) BlockLine() int {
	return this.Items[0].Tokens[0].Line
}

// Returns the column number this block starts on
func (this Block) BlockColumn() int {
	return this.Items[0].Tokens[0].Column
}

// #region Error handling
func (this Block) GetErrorPrefix(col int) string {
	return fmt.Sprintf(color.InRed("ERROR %s:%s:"), color.InYellow(this.BlockLine()), color.InYellow(col))
}

func (this Block) ErrorF(sourceItem Item, format string, args ...any) string {
	var message = fmt.Sprintf(format, args...)
	if sourceItem.IsError() {
		return fmt.Sprintf("%s\n", message)
	}

	var location = fmt.Sprintf("%s:%s:%s", color.InCyan(sourceItem.ItemFile()), color.InYellow(sourceItem.ItemLine()), color.InYellow(sourceItem.ItemColumn()))
	return fmt.Sprintf("%s\n\t%s", message, location)
}

//#endregion
