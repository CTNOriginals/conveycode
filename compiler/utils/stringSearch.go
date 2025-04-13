package utils

import (
	"slices"
	"strings"
)

// Check if an array contains any item from the search array
func ContainsListItem[T comparable](items []T, searchList []T) bool {
	for _, search := range searchList {
		if slices.Contains(items, search) {
			return true
		}
	}

	return false
}

// Gets the first possible range of characters that that are contained in valid
//
// The return values make up the index range of the substring heystack[start:end].
// If start and end both return 0, there was no match
func GetValidStringRange(heystack string, valid string, startIndex int) (start int, end int) {
	start = -1

	for i := startIndex; i < len(heystack); i++ {
		var char = rune(heystack[i])

		if start == -1 && strings.ContainsRune(valid, char) {
			start = i
		} else if start != -1 && !strings.ContainsRune(valid, char) {
			return start, i
		}
	}

	if start == -1 {
		return 0, 0
	}

	return start, len(heystack) //? Dont substract 1 to account for slices end range non-inclusivity
}

func ValidateString(heystack string, valid string) bool {
	start, end := GetValidStringRange(heystack, valid, 0)
	return start == 0 && end == len(heystack)
}
