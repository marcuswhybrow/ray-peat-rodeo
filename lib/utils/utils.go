package utils

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReturnOrPanic[A any](a A, err error) A {
	PanicOnErr(err)
	return a
}
