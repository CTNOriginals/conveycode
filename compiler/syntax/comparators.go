package syntax

var Comparators = syntaxMap{
	"==":  "equal",
	"!=":  "notEqual",
	"<":   "lessThan",
	"<=":  "lessThanEq",
	">":   "greaterThan",
	">=":  "greaterThanEq",
	"===": "strictEqual",
	// "true": "always", //?? Should this be included?
}
