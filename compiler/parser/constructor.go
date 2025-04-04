package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"
	"fmt"

	"github.com/TwiN/go-color"
)

type Constructor = func(block lexer.Block) []Instruction

type ConstructorMap = map[lexer.BlockType]Constructor

var constructors = ConstructorMap{
	lexer.Assignment: func(block lexer.Block) []Instruction {
		var prefix = "set"
		var ident = block.FindItem(lexer.Identifier)
		var item = block.FindItem(lexer.Value)

		//- Is the value a calculation of some sort?
		if item.Tokens.Contains(tokenizer.Operator) {
			return mathValueAssignment(block)
		}

		return []Instruction{NewInstruction(prefix, ident.Tokens.JoinValues(""), item.Tokens.JoinValues(" "))}
	},
}

func mathValueAssignment(block lexer.Block) (instructions []Instruction) {
	var itemIdent = block.FindItem(lexer.Identifier)
	var itemValue = block.FindItem(lexer.Value)

	//TODO Support multy math operations, like z = x + (2 * y) - 4
	for i, token := range itemValue.Tokens {
		if token.Typ == tokenizer.Operator {
			ope, ok := MathOperatorStrings[string(token.Val)]

			if !ok {
				fmt.Printf(color.InRed("Math operator '%s' is not yet defined in MathOperatorStrings\n"), string(token.Val))
				continue
			}

			instructions = append(instructions, NewInstruction("op", ope, itemIdent.Tokens.JoinValues(""), string(itemValue.Tokens[i-1].Val), string(itemValue.Tokens[i+1].Val)))
		}
	}

	return instructions
}
