package parser

import (
	"conveycode/compiler/lexer"
	"fmt"
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
	variables variableDefinitions
	methods   methodDefinitions
	children  []Scope
	//?? parent? string
}

var scopeCount int
var GlobalScope Scope

func GetScopeID() (id int) {
	id = scopeCount
	scopeCount += 1
	return id
}

func NewScope(context scopeContext) Scope {
	return Scope{
		id:        GetScopeID(),
		context:   context,
		variables: variableDefinitions{},
		methods:   methodDefinitions{},
	}
}

func InitializeScope() {
	GlobalScope.Clear()

	scopeCount = 0
	GlobalScope = NewScope(Global)
}

func (this Scope) String() string {
	return fmt.Sprintf(
		"%s {\n\tVariables: %d,\n\tMethods: %d,\n\tChildren: %d,\n}",
		this.GetLabelPrefix(),
		len(this.variables),
		len(this.methods),
		len(this.children),
	)
}

// #region Psuh
func (this *Scope) PushVariable(block lexer.Block) {
	this.variables[block.GetIdentifier()] = NewVariableDefinition(block, *this)
}
func (this *Scope) PushMethod(block lexer.Block) (def methodDefinition) {
	def = NewMethodDefinition(block)
	this.methods[block.GetIdentifier()] = def
	return def
}
func (this *Scope) PushChild(child Scope) {
	this.children = append(this.children, child)
}

//#endregion

func (this Scope) ContainsVariable(ident string) bool {
	var _, exists = this.variables[ident]
	return exists
}
func (this Scope) ContainsMethod(block lexer.Block) bool {
	var _, exists = this.methods[block.GetIdentifier()]
	return exists
}

func (this Scope) GetVariableByIdentifier(ident string) variableDefinition {
	if block, exists := this.variables[ident]; exists {
		return block
	}

	return NewVariableDefinition(lexer.NewBlock(lexer.BlockError), this)
}
func (this Scope) GetMethodByIdentifier(ident string) methodDefinition {
	if def, exists := this.methods[ident]; exists {
		return def
	}

	return NewMethodDefinition(lexer.NewBlock(lexer.BlockError))
}

func (this Scope) GetLabelPrefix() string {
	return fmt.Sprintf("%s%d", this.context, this.id)
}

func (this Scope) GetIdentifierLabel(ident string) string {
	return fmt.Sprintf("%s_%s", this.GetLabelPrefix(), ident)
}

func (this Scope) forEachChild(f func(child Scope)) {
	for _, child := range this.children {
		f(child)
	}
	for _, child := range this.children {
		child.forEachChild(f)
	}
}
func (this Scope) forEachMethod(f func(ident string, def methodDefinition)) {
	for ident, def := range this.methods {
		f(ident, def)
	}

	for _, child := range this.children {
		child.forEachMethod(f)
	}
}

func (this Scope) getParentScope() (parent Scope) {
	parent = GlobalScope
	GlobalScope.forEachChild(func(child Scope) {
		if this.id == child.id {
			return
		}
	})

	return parent
}

func (this Scope) getParentMethodScope() (parent Scope) {
	parent = this

	// fmt.Println(parent.context)
	for parent.context != Method && parent.context != Global {
		fmt.Println(parent.context)
		parent = parent.getParentScope()
	}

	return parent
}

// Clean up all memory contained in the scope
func (this *Scope) Clear() {
	clear(this.variables)
	clear(this.methods)
	for _, child := range this.children {
		child.Clear()
	}
}
