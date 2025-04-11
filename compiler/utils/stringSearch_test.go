package utils

import (
	"conveycode/internal"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type stringSearch_cases struct {
	valid  string
	index  int
	expect any
}

func Test_stringSearch(t *testing.T) {
	var cases = map[string][]stringSearch_cases{
		"C:/conveycode/compiler\\parser/construct_or.go:420({0x000015200?, 0x0?})": {
			{
				valid:  internal.WordCharacters + "./\\",
				index:  2,
				expect: "/conveycode/compiler\\parser/construct_or.go",
			},
			{
				valid:  internal.WordCharacters,
				index:  30,
				expect: "construct_or",
			},
			{
				valid:  internal.Numbers + ":",
				index:  3,
				expect: ":420",
			},
			{
				valid:  "({})x,? " + internal.Numbers,
				index:  0,
				expect: "420({0x000015200?, 0x0?})",
			},
			{
				valid:  "!",
				index:  0,
				expect: nil,
			},
		},
	}

	for heystack, caseList := range cases {
		Convey(fmt.Sprintf("The filepath is '%s'", heystack), t, func() {
			for _, caseValues := range caseList {
				Convey(fmt.Sprintf("Index: %d, Valid: '%s'", caseValues.index, caseValues.valid), func() {
					var start, end = GetValidStringRange(heystack, caseValues.valid, caseValues.index)
					if caseValues.expect != nil {
						So(heystack[start:end], ShouldEqual, caseValues.expect)
					} else {
						So(start, ShouldEqual, 0)
						So(end, ShouldEqual, 0)
					}
				})
			}
		})
	}
}
