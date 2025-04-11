package utils

import "regexp"

var (
	NonWordCharacter, _    = regexp.Compile(`\W`)
	WordCharacterStream, _ = regexp.Compile(`\w+`)
	NonNumberCharacter, _  = regexp.Compile(`\D`)
)
