package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/utils"
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

type methodDefinition struct {
	block lexer.Block

	ident      string
	parameters []string
	returnVar  variableDefinition
	scope      *Scope
}

type methodDefinitions map[string]methodDefinition

func (this methodDefinitions) String() (str string) {
	if len(this) == 0 {
		return fmt.Sprintf("Methods %s { }", color.InRed("EMPTY"))
	}

	var lines []string
	for key, val := range this {
		lines = append(lines, fmt.Sprintf("%s: %s", color.InBlue(key), fmt.Sprintf("(%+s)", strings.Join(val.parameters, ", "))))
	}

	var head = fmt.Sprintf("Methods %s", color.InCyan(this[utils.Keys(this)[0]].scope.GetLabelPrefix()))
	return fmt.Sprintf("%s {\n\t%s\n}", head, strings.Join(lines, "\n\t"))
}

func NewMethodDefinition(block lexer.Block) (def methodDefinition) {
	def = methodDefinition{
		block:      block,
		ident:      block.GetIdentifier(),
		parameters: block.GetArguments(),
		scope:      NewScope(Method),
	}

	var nullToken = ValuesToTokenList(block.BlockFile(), "null")

	def.returnVar = NewVariableDefinition(CreateMockAssignment(block.BlockFile(), true, "return", nullToken), def.scope)
	def.scope.variables[def.returnVar.ident] = def.returnVar

	var mockBlocks []lexer.Block
	for _, param := range def.parameters {
		//TODO add default param values instead of null
		mockBlocks = append(mockBlocks, CreateMockAssignment(block.BlockFile(), true, param, nullToken))
	}

	for _, block := range mockBlocks {
		if block.Typ == lexer.BlockEOF {
			continue
		}
		def.scope.PushVariable(block)
	}

	return def
}

func (this methodDefinition) getMethodLabel(scope Scope) (label string) {
	return scope.GetIdentifierLabel(this.scope.GetLabelPrefix())
}

func (this methodDefinition) getParamLabels() (labels []string) {
	for _, param := range this.parameters {
		labels = append(labels, this.scope.GetIdentifierLabel(param))
	}

	return labels
}
