package utils

import "sync"

// Runs func in parallel for each entry in slice and awaits all results
func Parallel[Item, Result any](slice []Item, f func(Item) Result) []Result {
	var waitGroup sync.WaitGroup

	count := len(slice)
	results := make(chan Result, count)
	waitGroup.Add(count)

	for _, item := range slice {
		go func(i Item) {
			defer waitGroup.Done()
			results <- f(i)
		}(item)
	}

	waitGroup.Wait()
	close(results)

	allResults := []Result{}
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}
