package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/types"
	"fmt"
	"runtime"
	"slices"

	"github.com/TwiN/go-color"
)

type Constructor = func(block lexer.Block, scope *Scope) (instructions []Instruction)

type ConstructorMap = map[lexer.BlockType]Constructor

var constructors = ConstructorMap{
	lexer.Assignment: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
		var prefix = "set"
		var itemIdent = block.FindItemByType(lexer.Identifier)
		var itemValue = block.FindItemByType(lexer.Value)

		if block.FindItemByType(lexer.Keyword).Typ != lexer.ItemError {
			if scope.ContainsVariable(block) {
				panic(scope.DuplicateVariable(block))
			} else {
				scope.PushVariable(block)
			}
		} else if !scope.ContainsVariable(block) {
			panic(scope.UndeclaredVariable(block))
		}

		//- Is the value a math operation of some sort?
		if itemValue.Tokens.Contains(tokenizer.Operator) {
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

		return []Instruction{NewInstruction(prefix, itemIdent.Tokens.JoinValues(""), itemValue.Tokens.JoinValues(" "))}
	},
	lexer.BuiltIn: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
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
	lexer.Statement: func(block lexer.Block, scope *Scope) (instructions []Instruction) {

		return instructions
	},
}

func Construct(prs *parser) (instructions []Instruction) {
	defer func() {
		if errMsg := recover(); errMsg != nil {
			fmt.Println(errMsg)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, false)
			fmt.Printf("%s", buf)
		}
	}()

	for _, block := range prs.blocks {
		var constructor, ok = constructors[block.Typ]

		//- Does the block type have a constructor defined
		if !ok {
			continue
		}

		instructions = append(instructions, constructor(block, prs.scope)...)
	}

	return instructions
}
