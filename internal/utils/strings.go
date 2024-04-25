package utils

// Returns string s if i does not equal one
func Pluralise(i int, s string) string {
	if i != 1 {
		return s
	}
	return ""
}
