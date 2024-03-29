package must

// apache 2.0 antlabs
func TakeOneDiscardBool[T any](v T, ok bool) T {
	if !ok {
		panic("ok is false")
	}
	return v
}

func TakeOne[T any](v T, err error) T {
	if err != nil {
		panic(err.Error())
	}
	return v
}

func TakeTwo[T, U any](a T, b U, err error) (T, U) {
	if err != nil {
		panic(err.Error())
	}
	return a, b
}

func TakeThree[T, U, V any](a T, b U, c V, err error) (T, U, V) {
	if err != nil {
		panic(err.Error())
	}

	return a, b, c
}

func TakeOneErr[T any](v T, err error) error {
	return err
}

func TakeOneBool[T any](v T, ok bool) bool {
	return ok
}

func TakeTwoErr[T, U any](a T, b U, err error) error {
	return err
}

func TakeThreeErr[T, U, V any](a T, b U, c V, err error) error {
	return err
}
