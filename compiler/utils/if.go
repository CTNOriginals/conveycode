package utils

func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
