package syntax

import (
	"conveycode/constents"
	"slices"
)

var Commands = []string{
	"read",
	"write",
	"draw",
	"print",
	"printchar",

	"format",
	"drawflush",
	"printflush",
	"getlink",
	"control",
	"radar",
	"sensor",

	"lookup",
	"packcolor",

	"wait",
	"stop",
	"end",

	"ubind",
	"ucontrol",
	"uradar",
	"ulocate",
}

var CommandAliases = map[string][]string{
	"printflush": {"flush", "pflush"},
	"drawflush":  {"dflush"},
	"printchar":  {"printc"},
}

func IsCommand(input string) bool {
	if slices.Contains(Commands, input) {
		return true
	}

	for _, aliases := range CommandAliases {
		if slices.Contains(aliases, input) {
			return true
		}
	}

	return false
}

func GetCommandPrefix(input string) string {
	if slices.Contains(Commands, input) {
		return input
	}

	for key, aliases := range CommandAliases {
		if slices.Contains(aliases, input) {
			return key
		}
	}

	return constents.StringError
}

/* All of mindustries commands and default values
read result cell1 0
write result cell1 0
draw clear 0 0 0 0 0 0
print "frog"
printchar 65
format "frog"
drawflush display1
printflush message1
getlink result 0
control enabled block1 0 0 0 0
radar enemy any any distance turret1 1 result
sensor result block1 @copper
set result 0
op add result a b
lookup item result 0
packcolor result 1 0 0 1
wait 0.5
stop
end
jump -1 notEqual x false
ubind @poly
ucontrol move 0 0 0 0 0
uradar enemy any any distance 0 1 result
ulocate building core true @copper outx outy found building
*/
