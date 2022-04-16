package must

func One[T any](v T, err error) T {
	if err != nil {
		panic(err.Error())
	}
	return v
}

func Two[T, U any](a T, b U, err error) (T, U) {
	if err != nil {
		panic(err.Error())
	}
	return a, b
}

func Three[T, U, V any](a T, b U, c V, err error) (T, U, V) {
	if err != nil {
		panic(err.Error())
	}

	return a, b, c
}
