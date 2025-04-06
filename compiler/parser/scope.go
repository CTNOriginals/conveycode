package parser

import (
	"conveycode/compiler/lexer"
	"slices"

	"github.com/TwiN/go-color"
)

type scopeContext string

const (
	_ scopeContext = ""

	Global    = "global"
	Statement = "statement"
	Method    = "method"
)

type Scope struct {
	id        int
	context   scopeContext
	Variables []lexer.Block
	Methods   []string
	Children  []Scope
	//?? parent? string
}

var scopeCount int
var GlobalScope Scope

func GetScopeID() (id int) {
	id = scopeCount
	scopeCount += 1
	return id
}

func InitializeScope() {
	GlobalScope.Clear()

	scopeCount = 0
	GlobalScope = Scope{
		id:      GetScopeID(),
		context: Global,
	}
}

func NewScope(context scopeContext) Scope {
	return Scope{
		id:      GetScopeID(),
		context: context,
	}
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

// Clean up all memory contained in the scope
func (this *Scope) Clear() {
	this.Variables = this.Variables[:0]
	this.Methods = this.Methods[:0]
	for _, child := range this.Children {
		child.Clear()
	}
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
