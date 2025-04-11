package utils

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type file_cases struct {
	path   string
	dir    string
	file   string
	name   string
	ext    string
	line   int
	column int
	full   string
	trail  string
}

func Test_filePath(t *testing.T) {
	var cases = map[string]file_cases{
		"C:/foo/bar\\wip/fileName.xten:420:84({0x000015200?, 0x0?})": {
			path:   "C:/foo/bar/wip",
			dir:    "wip",
			file:   "fileName.xten",
			name:   "fileName",
			ext:    "xten",
			line:   420,
			column: 84,
			full:   "C:/foo/bar/wip/fileName.xten",
			trail:  "({0x000015200?, 0x0?})",
		},
		"fileName.xten": {
			path:  "",
			dir:   "",
			file:  "fileName.xten",
			name:  "fileName",
			ext:   "xten",
			line:  0,
			full:  "fileName.xten",
			trail: "",
		},
	}

	for rawPath, expect := range cases {
		Convey(fmt.Sprintf("File Path: '%s'", rawPath), t, func() {
			var parsed = ParseFilePath(rawPath)

			So(parsed.Path, ShouldEqual, expect.path)
			So(parsed.Dir, ShouldEqual, expect.dir)
			So(parsed.File, ShouldEqual, expect.file)
			So(parsed.Name, ShouldEqual, expect.name)
			So(parsed.Ext, ShouldEqual, expect.ext)
			So(parsed.Line, ShouldEqual, expect.line)
			So(parsed.Column, ShouldEqual, expect.column)
			So(parsed.Full, ShouldEqual, expect.full)
			So(parsed.Trail, ShouldEqual, expect.trail)
		})
	}
}
