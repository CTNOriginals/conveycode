package utils

func Splice[T comparable](s []T, start int, count int) (ret []T, del []T) {
	ret = make([]T, len(s)-count)
	del = make([]T, count)

	deleteCount := 0
	for i, item := range s {
		if i >= start && i < (start+count) {
			del[deleteCount] = item
			deleteCount++
			continue
		}

		ret[i-deleteCount] = item
	}

	return
}
