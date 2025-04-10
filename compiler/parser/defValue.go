package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"

	"github.com/TwiN/go-color"
)

type valueDefinition struct {
	block lexer.Block
	scope Scope

	item lexer.Item
}

func NewValueDefinition(block lexer.Block, item lexer.Item, scope Scope) valueDefinition {
	return valueDefinition{
		block: block,
		scope: scope,
		item:  item,
	}
}

func NewValueDefinitionFromBlock(block lexer.Block, valueType lexer.ItemType, scope Scope) valueDefinition {
	return valueDefinition{
		block: block,
		scope: scope,
		item:  block.FindItemByType(valueType),
	}
}

// Converts all items of valueType into a valueDefinition and returns them
func ConvertAllToValueDefinition(block lexer.Block, valueType lexer.ItemType, scope Scope) (defList []valueDefinition) {
	for _, item := range block.GetItemsOfType(valueType) {
		defList = append(defList, NewValueDefinition(block, item, scope))
	}

	return defList
}

func (this valueDefinition) getVariableOrigin(ident string) (variableDef variableDefinition) {
	var parentScope = this.scope
	variableDef = parentScope.GetVariableByIdentifier(ident)

	for variableDef.IsError() && parentScope.context != Global {
		parentScope = parentScope.getParentScope()
		variableDef = parentScope.GetVariableByIdentifier(ident)
	}

	return variableDef
}

func (this valueDefinition) getValidated() (validated lexer.Item) {
	for _, token := range this.item.Tokens {
		if token.Typ != tokenizer.Text {
			validated.Tokens = append(validated.Tokens, token)
			continue
		}

		var variableDef = this.getVariableOrigin(string(token.Val))

		if variableDef.IsError() {
			panic(this.block.ErrorF(
				this.item,
				"Unknown identifier '%s'",
				color.InPurple(string(token.Val)),
			))
		}

		token.Val = []rune(this.scope.GetIdentifierLabel(string(token.Val)))
		validated.Tokens = append(validated.Tokens, token)
	}

	return validated
}

// Returns all validated values and excludes any non-value item (like operators or seperators)
func (this valueDefinition) valuesOnly() (values lexer.Item) {
	for _, token := range this.getValidated().Tokens {
		if token.Typ.IsValue() {
			values.Tokens = append(values.Tokens, token)
		}
	}

	return values
}

// Returns all the validated values from valuesOnly() as strings
func (this valueDefinition) rawValues() (values []string) {
	for _, token := range this.valuesOnly().Tokens {
		values = append(values, string(token.Val))
	}

	return values
}
