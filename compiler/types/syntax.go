package types

type syntaxMap = map[string]string

var MathOperators = syntaxMap{
	"+": "add",
	"-": "sub",
	"*": "mul",
	"/": "div",
	"%": "mod",
}

var Comparator = syntaxMap{
	"==":  "equal",
	"!=":  "notEqual",
	"<":   "lessThan",
	"<=":  "lessThanEq",
	">":   "greaterThan",
	">=":  "greaterThanEq",
	"===": "strictEqual",
	// "true": "always", //?? Should this be included?
}

var Commands = syntaxMap{
	"print": "print",
	"flush": "printflush",
}
