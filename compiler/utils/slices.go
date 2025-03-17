package utils

func Splice[T comparable](slice []T, start int, count int) (remaining []T, deleted []T) {
	remaining = make([]T, len(slice)-count)
	deleted = make([]T, count)

	deleteCount := 0
	for i, item := range slice {
		if i >= start && i < (start+count) {
			deleted[deleteCount] = item
			deleteCount++
			continue
		}

		remaining[i-deleteCount] = item
	}

	return
}
