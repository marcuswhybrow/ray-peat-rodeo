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

func EnsureMapExists[A any](testInstance map[string]A) map[string]A {
	if testInstance == nil {
		return map[string]A{}
	}
	return testInstance
}
