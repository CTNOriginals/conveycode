package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"

	"github.com/TwiN/go-color"
)

type valueDefinition struct {
	block lexer.Block
	scope *Scope

	item lexer.Item
}

func NewValueDefinition(block lexer.Block, item lexer.Item, scope *Scope) valueDefinition {
	return valueDefinition{
		block: block,
		scope: scope,
		item:  item,
	}
}

func NewValueDefinitionFromBlock(block lexer.Block, valueType lexer.ItemType, scope *Scope) valueDefinition {
	return valueDefinition{
		block: block,
		scope: scope,
		item:  block.FindItemByType(valueType),
	}
}

// Converts all items of valueType into a valueDefinition and returns them
func ConvertAllToValueDefinition(block lexer.Block, valueType lexer.ItemType, scope *Scope) (defList []valueDefinition) {
	for _, item := range block.GetItemsOfType(valueType) {
		defList = append(defList, NewValueDefinition(block, item, scope))
	}

	return defList
}

func (this valueDefinition) validate() (validated lexer.Item) {
	for i := 0; i < len(this.item.Tokens); i++ {
		var token = this.item.Tokens[i]

		if token.Typ != tokenizer.Text {
			validated.Tokens = append(validated.Tokens, token)
			continue
		}

		var ident = string(token.Val)

		if this.scope.identIsVariable(ident) {
			token.Val = []rune(this.scope.getVariableOrigin(ident).scope.GetIdentifierLabel(ident))
		} else if this.scope.identIsMethod(ident) {
			token.Val = []rune(this.scope.getMethodOrigin(ident).scope.GetIdentifierLabel("return"))
			var matchIndex = this.item.Tokens.GetMatchingBracketIndex(i)
			if matchIndex != -1 {
				i = matchIndex
			}
		} else {
			panic(this.block.ErrorF(
				this.item,
				"Unknown identifier '%s'",
				color.InPurple(ident),
			))
		}

		validated.Tokens = append(validated.Tokens, token)
	}

	return validated
}

// Returns all validated values and excludes any non-value item (like operators or seperators)
func (this valueDefinition) valuesOnly() (values lexer.Item) {
	for _, token := range this.validate().Tokens {
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

// Returns the instruction lines that need to prepend the actual
// values for all of them to work and contain the correct things
func (this valueDefinition) getValueInstructions() (instructions []Instruction) {
	for i := 0; i < len(this.item.Tokens); i++ {
		var token = this.item.Tokens[i]

		if token.Typ != tokenizer.Text {
			continue
		}

		var ident = string(token.Val)

		//TODO Support multi math opperations
		if this.scope.identIsMethod(ident) {
			var endArgs = this.item.Tokens.GetMatchingBracketIndex(i + 1)

			if endArgs == -1 {
				panic(this.block.ErrorF(
					this.item,
					"Incorrect method call '%s'",
					string(token.Val),
				))
			}

			var block = CreateMockMethodCall(token.File, token, this.item.Tokens[i+1:endArgs+1])
			var prs = Parse([]lexer.Block{block}, this.scope)
			instructions = append(instructions, Construct(prs)...)

			i = endArgs
			continue
		}
	}
	return instructions
}
