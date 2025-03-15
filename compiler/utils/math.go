package utils

func Max[T int | float64](a T, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

// Always returns a positive number, signed is true if the number was negative
func Abs[T int | float64](n T) (num T, signed bool) {
	if n < 0 {
		return n - n - n, true
	}

	return n, false
}
