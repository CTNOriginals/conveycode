package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/utils"
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

type variableDefinition struct {
	block lexer.Block
	scope Scope

	ident string
	value valueDefinition
}

type variableDefinitions map[string]variableDefinition

func (this variableDefinitions) String() (str string) {
	if len(this) == 0 {
		return fmt.Sprintf("Variables %s { }", color.InRed("EMPTY"))
	}

	var lines []string
	for key, val := range this {
		lines = append(lines, fmt.Sprintf("%s: %s", color.InBlue(key), val.value.item.Tokens.JoinValues(" ")))
	}

	var head = fmt.Sprintf("Variables %s", color.InCyan(this[utils.Keys(this)[0]].scope.GetLabelPrefix()))
	return fmt.Sprintf("%s {\n\t%s\n}", head, strings.Join(lines, "\n\t"))
}

func NewVariableDefinition(block lexer.Block, scope Scope) variableDefinition {
	return variableDefinition{
		block: block,
		scope: scope,

		ident: block.GetIdentifier(),
		value: NewValueDefinitionFromBlock(block, lexer.Value, scope),
	}
}

func (this variableDefinition) IsError() bool {
	return this.block.IsError()
}
