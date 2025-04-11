package internal

import (
	"path/filepath"
	"runtime"
	"strings"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	RootPath = filepath.Join(filepath.Dir(b), "../")
)

var ProjectName = strings.Split(RootPath, "\\")[len(strings.Split(RootPath, "\\"))-1]
