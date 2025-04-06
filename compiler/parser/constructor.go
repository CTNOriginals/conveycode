package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/syntax"
	"conveycode/compiler/tokenizer"
	"fmt"
	"runtime"
	"slices"

	"github.com/TwiN/go-color"
)

type Constructor = func(block lexer.Block, scope *Scope) (instructions []Instruction)

type ConstructorMap = map[lexer.BlockType]Constructor

var constructors ConstructorMap

func init() {
	constructors = ConstructorMap{
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
						ope, ok := syntax.MathOperators[string(token.Val)]

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

			instructions = append(instructions, NewInstruction(syntax.Commands[command]))

			for _, token := range args.Tokens {
				if token.Typ == tokenizer.Seperator {
					instructions = append(instructions, NewInstruction(syntax.Commands[command]))
				} else if slices.Contains(tokenizer.ValueTokenTypes, token.Typ) {
					instructions[len(instructions)-1].Push(string(token.Val))
				}
			}

			return instructions
		},
		lexer.Statement: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
			var prefix = fmt.Sprintf("%s%d", scope.context, scope.id)
			var jumps []Instruction
			var exitLabel string

			var subScope Scope
			var label string
			var comparable string
			var subInstructions []Instruction

			for i := len(block.Items) - 1; i >= 0; i-- {
				var item = block.Items[i]

				switch item.Typ {
				case lexer.Scope:
					var blocks = lexer.Lex(item.Tokens).Construct()
					var prs = Parse(blocks, &subScope)
					subInstructions = Construct(prs)

					subScope = NewScope(Statement)
					scope.Children = append(scope.Children, subScope)

					label = fmt.Sprintf("%s_%s%d", prefix, Statement, subScope.id)
					if i == len(block.Items)-1 {
						exitLabel = label
					}
				case lexer.Condition:
					var comparator = ""
					var values []string

					for _, token := range item.Tokens {
						if token.Typ == tokenizer.Operator {
							if comparator != "" {
								panic(block.ErrorF(item, "Multi conditional statements are not supported yet"))
							}
							comparator = string(token.Val)
						} else if slices.Contains(tokenizer.ValueTokenTypes, token.Typ) {
							if len(values) >= 2 {
								panic(block.ErrorF(item, "Conditional statements cant compare more then 2 values"))
							}
							values = append(values, string(token.Val))
						}
					}

					if len(values) < 2 {
						panic(block.ErrorF(item, "Conditional statements require 2 comparable values"))
					}

					comparable = fmt.Sprintf("%s %v %v", syntax.Comparators[string(comparator)], values[0], values[1])
				case lexer.Keyword:
					if item.ValueString() != "else" {
						jumps = append([]Instruction{NewInstruction("jump", label, comparable)}, jumps...)
						instructions = append(instructions, NewInstruction(label+":"))
					}

					instructions = append(instructions, subInstructions...)

					if i > 0 {
						//? The first (bottom in compiled result) statement doesnt need an exit jump as its already at the bottom of everything
						instructions = append(instructions, NewInstruction("jump", exitLabel, "always"))
					}
				}
			}

			instructions = append(jumps, instructions...)
			instructions = append(instructions, NewInstruction(exitLabel+":"))

			return instructions
		},
	}
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
