package utils

import "slices"

// Check if an array contains any item from the search array
func ContainsListItem(items []string, searchList []string) bool {
	for _, search := range searchList {
		if slices.Contains(items, search) {
			return true
		}
	}
	return false
}
