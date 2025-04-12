package constents

import (
	"path/filepath"
	"runtime"
	"strings"
)

// String Characters
const (
	Alphabet           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers            = "1234567890"
	AlphaNumaric       = Alphabet + Numbers
	WordCharacters     = AlphaNumaric + "_"
	FileNameCharacters = WordCharacters + " -."
)

// String Flags
const (
	StringError           = "<ERROR>"
	ReturnLinePlaceholder = "<RETURN_LINE>"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	RootPath    = filepath.Join(filepath.Dir(b), "../")
	ProjectName = strings.Split(RootPath, "\\")[len(strings.Split(RootPath, "\\"))-1]
)
