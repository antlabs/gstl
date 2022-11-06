package ifop

// apache 2.0 antlabs
func If[T any](cond bool, t T) (zero T) {
	if cond {
		return t
	}
	return
}

func IfElse[T any](cond bool, ifVal T, elseVal T) T {
	if cond {
		return ifVal
	}
	return elseVal
}
