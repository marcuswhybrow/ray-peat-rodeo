package utils

// Takes a slice ts and returns, at most, the first i elements.
func AtMost[T any](ts []T, i int) []T {
	return ts[:min(i, len(ts))]
}
