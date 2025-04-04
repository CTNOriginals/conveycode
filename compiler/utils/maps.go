package utils

func Keys[K comparable, V any](m map[K]V) (ret []K) {
	ret = make([]K, len(m))

	i := 0
	for key := range m {
		ret[i] = key
		i++
	}

	return ret
}

func Values[K comparable, V any](m map[K]V) (ret []V) {
	ret = make([]V, len(m))

	i := 0
	for key := range m {
		ret[i] = m[key]
		i++
	}

	return ret
}
