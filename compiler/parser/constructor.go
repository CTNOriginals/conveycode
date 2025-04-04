package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/types"
	"fmt"
	"slices"

	"github.com/TwiN/go-color"
)

type Constructor = func(block lexer.Block) (instructions []Instruction)

type ConstructorMap = map[lexer.BlockType]Constructor

var constructors = ConstructorMap{
	lexer.Assignment: func(block lexer.Block) (instructions []Instruction) {
		var prefix = "set"
		var ident = block.FindItemByType(lexer.Identifier)
		var item = block.FindItemByType(lexer.Value)

		//- Is the value a calculation of some sort?
		if item.Tokens.Contains(tokenizer.Operator) {
			return mathValueAssignment(block)
		}

		return []Instruction{NewInstruction(prefix, ident.Tokens.JoinValues(""), item.Tokens.JoinValues(" "))}
	},
	lexer.BuiltIn: func(block lexer.Block) (instructions []Instruction) {
		var command = string(block.FindItemByType(lexer.Command).Tokens.FindTokenByType(tokenizer.Command).Val)
		var args = block.FindItemByType(lexer.Arguments)

		instructions = append(instructions, NewInstruction(types.Commands[command]))

		for _, token := range args.Tokens {
			if token.Typ == tokenizer.Seperator {
				instructions = append(instructions, NewInstruction(types.Commands[command]))
			} else if slices.Contains(tokenizer.ValueTokenTypes, token.Typ) {
				instructions[len(instructions)-1].Push(string(token.Val))
			}
		}

		return instructions
	},
}

func mathValueAssignment(block lexer.Block) (instructions []Instruction) {
	var itemIdent = block.FindItemByType(lexer.Identifier)
	var itemValue = block.FindItemByType(lexer.Value)

	//TODO Support multy math operations, like z = x + (2 * y) - 4
	for i, token := range itemValue.Tokens {
		if token.Typ == tokenizer.Operator {
			ope, ok := types.MathOperators[string(token.Val)]

			if !ok {
				fmt.Printf(color.InRed("Math operator '%s' is not yet defined in MathOperatorStrings\n"), string(token.Val))
				continue
			}

			instructions = append(instructions, NewInstruction("op", ope, itemIdent.Tokens.JoinValues(""), string(itemValue.Tokens[i-1].Val), string(itemValue.Tokens[i+1].Val)))
		}
	}

	return instructions
}
