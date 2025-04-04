package types

type syntaxMap = map[string]string

var MathOperators = syntaxMap{
	"+": "add",
	"-": "sub",
	"*": "mul",
	"/": "div",
	"%": "mod",
}

var Commands = syntaxMap{
	"print": "print",
	"flush": "printflush",
}
