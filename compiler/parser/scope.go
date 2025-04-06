package parser

import (
	"conveycode/compiler/lexer"
	"slices"

	"github.com/TwiN/go-color"
)

type scopeContext string

const (
	_ scopeContext = ""

	Global          = "global"
	IfStatement     = "if"
	ElseIfStatement = "elseif"
	ElseStatement   = "else"
	ForLoop         = "for"
	WhileLoop       = "while"
	Method          = "method"
)

type Scope struct {
	context      scopeContext
	Instructions []Instruction
	Variables    []lexer.Block
	Methods      []string
	Children     []Scope
	//?? parent? string
}

var GlobalScope = Scope{
	context: Global,
}

func (this *Scope) PushVariable(blocks ...lexer.Block) {
	this.Variables = append(this.Variables, blocks...)
}

func (this Scope) ContainsVariable(block lexer.Block) bool {
	for _, b := range this.Variables {
		if b.FindItemByType(lexer.Identifier).Compare(block.FindItemByType(lexer.Identifier)) {
			return true
		}
	}

	return false
}
func (this Scope) ContainsMethod(ident string) bool {
	return slices.Contains(this.Methods, ident)
}

func (this Scope) GetVariableByIdentifier(ident string) lexer.Block {
	for _, block := range this.Variables {
		if block.FindItemByType(lexer.Identifier).ValueString() == ident {
			return block
		}
	}

	return lexer.NewBlock(lexer.BlockError)
}

// #region Error Logging
func (this Scope) DuplicateVariable(block lexer.Block) string {
	var itemIdent = block.FindItemByType(lexer.Identifier)
	var origVar = this.GetVariableByIdentifier(itemIdent.ValueString())
	return block.ErrorF(
		itemIdent,
		"Variable '%s' has already been declared at %s:%s",
		color.InBlue(itemIdent.ValueString()),
		color.InYellow(origVar.BlockLine()),
		color.InYellow(origVar.FindItemByType(lexer.Identifier).ItemColumn()),
	)
}

func (this Scope) UndeclaredVariable(block lexer.Block) string {
	var itemIdent = block.FindItemByType(lexer.Identifier)
	return block.ErrorF(
		itemIdent,
		"Variable '%s' has not been declared",
		color.InBlue(itemIdent.ValueString()),
	)
}

//#endregion
