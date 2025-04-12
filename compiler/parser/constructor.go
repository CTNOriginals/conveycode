package parser

import (
	"conveycode/compiler/lexer"
	"conveycode/compiler/syntax"
	"conveycode/compiler/tokenizer"
	"conveycode/compiler/utils"
	"conveycode/constents"
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

type Constructor = func(block lexer.Block, scope *Scope) (instructions []Instruction)

type ConstructorMap = map[lexer.BlockType]Constructor

var constructors ConstructorMap
var methodDefinitionBodies []Instruction

var instructionSpacing = true //? for debugging readability

func init() {
	constructors = ConstructorMap{
		lexer.Assignment: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
			var prefix = "set"
			var itemIdent = block.FindItemByType(lexer.Identifier)
			var valDef = NewValueDefinitionFromBlock(block, lexer.Value, scope)
			var label = scope.GetIdentifierLabel(itemIdent.ValueString())

			//- Validation
			if block.FindItemByType(lexer.Keyword).Typ != lexer.ItemError {
				if scope.ContainsVariable(itemIdent.ValueString()) {
					var itemIdent = block.FindItemByType(lexer.Identifier)
					var origVar = scope.GetVariableByIdentifier(itemIdent.ValueString())
					panic(block.ErrorF(
						itemIdent,
						"Variable '%s' has already been declared at %s:%s",
						color.InBlue(itemIdent.ValueString()),
						color.InYellow(origVar.block.BlockLine()),
						color.InYellow(origVar.block.FindItemByType(lexer.Identifier).ItemColumn()),
					))
				} else {
					scope.PushVariable(block)
				}
			} else if !scope.ContainsVariable(itemIdent.ValueString()) {
				var itemIdent = block.FindItemByType(lexer.Identifier)
				panic(block.ErrorF(
					itemIdent,
					"Variable '%s' has not been declared",
					color.InBlue(itemIdent.ValueString()),
				))
			}

			//- Is the value a math operation of some sort?
			if valDef.item.Tokens.Contains(tokenizer.Operator) {
				for _, token := range valDef.item.Tokens {
					if token.Typ == tokenizer.Operator {
						ope, ok := syntax.MathOperators[string(token.Val)]

						if !ok {
							fmt.Printf(color.InRed("Math operator '%s' is not yet defined in MathOperatorStrings\n"), string(token.Val))
							continue
						}

						var instruction = NewInstruction("op", ope, label)
						instruction.Push(valDef.valuesOnly().Tokens.ValuesAsStringArray()...)

						instructions = append(instructions, instruction)
					}
				}

				return instructions
			}

			return []Instruction{NewInstruction(prefix, label, valDef.valuesOnly().ValueString())}
		},
		lexer.BuiltIn: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
			var command = string(block.FindItemByType(lexer.Command).Tokens.FindTokenByType(tokenizer.Command).Val)
			var valueDef = NewValueDefinitionFromBlock(block, lexer.Arguments, scope)
			var args = valueDef.rawValues()

			for _, arg := range args {
				instructions = append(instructions, NewInstruction(syntax.Commands[command], arg))
			}

			return instructions
		},
		lexer.Statement: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
			var prefix = fmt.Sprintf("%s%d", scope.context, scope.id)
			var jumps []Instruction
			var exitLabel string
			var exitInstruction Instruction

			var subScope *Scope
			var label string
			var comparable string
			var subInstructions []Instruction

			for i := len(block.Items) - 1; i >= 0; i-- {
				var item = block.Items[i]

				switch item.Typ {
				case lexer.Scope:
					subScope = NewScope(Statement)
					scope.PushChild(subScope)

					var blocks = lexer.Lex(item.Tokens).Construct()
					var prs = Parse(blocks, subScope)
					subInstructions = Construct(prs)

					label = fmt.Sprintf("%s_%s%d", prefix, Statement, subScope.id)
					if i == len(block.Items)-1 {
						exitLabel = fmt.Sprintf("%s_%s", label, "exit")
						exitInstruction = NewInstruction("jump", exitLabel, "always", utils.If(instructionSpacing, "\n", ""))
					}
				case lexer.Condition:
					var comparator = ""
					var valueDef = NewValueDefinition(block, item, scope)
					var values = valueDef.rawValues()

					for _, token := range valueDef.item.Tokens {
						if token.Typ == tokenizer.Operator {
							if comparator != "" {
								panic(block.ErrorF(item, "Multi conditional statements are not supported yet"))
							}
							comparator = string(token.Val)
						}
					}

					if len(values) != 2 {
						panic(block.ErrorF(item, "Conditional statements expect 2 comparable values but received %s", color.InRed(len(values))))
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
						instructions = append(instructions, exitInstruction)
					}
				}
			}

			// TODO Save some space here by checing if the condition does not include a "==="
			// TODO and if so, invert the last condition and point it to the exit
			// TODO this way, the last conditional (not else) statement can just flow through to the next line if "true"
			jumps = append(jumps, exitInstruction)

			instructions = append(jumps, instructions...)
			instructions = append(instructions, NewInstruction(exitLabel+":"))

			return instructions
		},
		lexer.Method: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
			var identItem = block.FindItemByType(lexer.Identifier)

			//- Is this method already defined anywhere?
			if scope.ContainsMethod(block) {
				var orig = scope.GetMethodByIdentifier(identItem.ValueString())
				panic(block.ErrorF(
					identItem,
					"Method '%s' has already be defined on line %s",
					color.InPurple(identItem.ValueString()),
					color.InYellow(orig.block.FindItemByType(lexer.Identifier).ItemLine()),
				))
			}

			var def = scope.PushMethod(block)

			var body = def.block.FindItemByType(lexer.Scope)
			var prs = ParseItem(body, def.scope)
			var methodInstructions = Construct(prs)

			methodDefinitionBodies = append(methodDefinitionBodies, NewInstruction(def.getMethodLabel(*scope)+":"))
			methodDefinitionBodies = append(methodDefinitionBodies, methodInstructions...)
			methodDefinitionBodies = append(methodDefinitionBodies, NewInstruction(
				"set",
				"@counter",
				fmt.Sprintf("%s_%s", def.scope.GetLabelPrefix(), "caller-adress"),
			))

			return instructions
		},
		lexer.Call: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
			var ident = block.GetIdentifier()

			var def = scope.GetMethodByIdentifier(ident)
			var method = def.block

			if method.IsError() {
				panic(block.ErrorF(
					block.FindItemByType(lexer.Arguments),
					"Unknown method '%s'",
					color.InPurple(ident),
				))
			}

			if block.GetArity() != method.GetArity() {
				panic(block.ErrorF(
					block.FindItemByType(lexer.Arguments),
					"Method '%s%s' at line %s:\nExpected %s parameters, but received %s instead",
					color.InPurple(ident),
					color.InBlue(method.FindItemByType(lexer.Parameters).ValueString()),
					color.InYellow(method.BlockLine()),
					color.InGreen(method.GetArity()),
					color.InRed(block.GetArity()),
				))
			}

			var paramLabels = def.getParamLabels()
			for i, arg := range block.GetArguments() {
				instructions = append(instructions, NewInstruction("set", paramLabels[i], arg))
			}

			instructions = append(instructions,
				NewInstruction("set", def.scope.GetIdentifierLabel("caller-adress"), constents.ReturnLinePlaceholder),
				NewInstruction("jump", def.getMethodLabel(*scope), "always"),
			)

			return instructions
		},
		lexer.Return: func(block lexer.Block, scope *Scope) (instructions []Instruction) {
			var parentScope = scope.getSurroundingMethod()
			var valueDef = NewValueDefinitionFromBlock(block, lexer.Value, parentScope)
			// var mockAssignment = fmt.Sprintf("var %s = %s", "return", valueDef.item.Tokens.Stream())
			var mockAssignment = CreateMockAssignment(block.BlockFile(), false, "return", valueDef.item.Tokens)

			return Construct(Parse([]lexer.Block{mockAssignment}, valueDef.scope))
		},
	}
}

func Construct(prs *parser) (instructions []Instruction) {
	defer func() {
		if errMsg := recover(); errMsg != nil {
			fmt.Println(errMsg)
			utils.PrintStackTrace(9)
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

	if prs.scope.context == Global {
		instructions = append(instructions, NewInstruction("end"))
	}

	instructions = append(instructions, methodDefinitionBodies...)

	//? Check each line and replace the <RETURN_LINE> parts with the correct line number
	for i, instruction := range instructions {
		var lineNum = i + 1
		for j, part := range instruction.Parts {
			if part == constents.ReturnLinePlaceholder {
				instruction.Parts[j] = fmt.Sprint(lineNum + 2)
			}
		}
	}

	return instructions
}

func ValuesToTokenList(file string, values ...string) tokenizer.TokenList {
	return tokenizer.TokenizeContent(file, strings.Join(values, " "))
}
func CreateMockAssignment(file string, newVar bool, ident string, tokens tokenizer.TokenList) (block lexer.Block) {
	block = lexer.NewBlock(lexer.Assignment)

	block.Items = []lexer.Item{
		{Typ: lexer.Identifier,
			Tokens: append(tokenizer.NewTokenList(), tokenizer.NewToken(file, tokenizer.Text, []rune(ident), 0, 5)),
		},
		{Typ: lexer.Operator,
			Tokens: append(tokenizer.NewTokenList(), tokenizer.NewToken(file, tokenizer.Text, []rune("="), 0, len(ident)+2)),
		},
		{Typ: lexer.Value,
			Tokens: tokens,
		},
	}

	if newVar {
		block.Items = append([]lexer.Item{lexer.Item{Typ: lexer.Keyword,
			Tokens: append(tokenizer.NewTokenList(), tokenizer.NewToken(file, tokenizer.Text, []rune("var"), 0, 1)),
		}}, block.Items...)
	}

	return block
}
