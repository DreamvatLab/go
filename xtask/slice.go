package xtask

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// ParallelRunSlice concurrently processes each element in the input slice using a fixed number of worker goroutines.
// The function provides controlled concurrency by limiting the number of simultaneous workers.
//
// Parameters:
//   - limit: Maximum number of concurrent workers. If <= 0, uses runtime.GOMAXPROCS(0)
//   - slice: Input slice of type T to be processed
//   - processor: Function that processes a single element of type T and returns (interface{}, error)
//
// Returns:
//   - []*TaskResult: Slice of results in the same order as the input slice
//
// Each TaskResult contains:
//   - Result: The processed output (interface{})
//   - Error: Any error that occurred during processing, including panic recoveries
//
// Example:
//
//	results := ParallelRunSlice(4, []int{1, 2, 3}, func(n int) (interface{}, error) {
//	    return n * 2, nil
//	})
func ParallelRunSlice[T any](limit int, slice []T, processor func(T) (interface{}, error)) []*TaskResult {
	if limit <= 0 {
		limit = runtime.GOMAXPROCS(0)
	}

	taskCount := len(slice)
	results := make([]*TaskResult, taskCount)
	var wg sync.WaitGroup
	taskIndices := make(chan int)

	// Start worker goroutines
	for w := 0; w < limit; w++ {
		go func() {
			for index := range taskIndices {
				func(i int) {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							results[i] = &TaskResult{
								Result: nil,
								Error:  fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack()),
							}
						}
					}()

					result, err := processor(slice[i])
					results[i] = &TaskResult{
						Result: result,
						Error:  err,
					}
				}(index)
			}
		}()
	}

	// Send task indices to workers
	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		taskIndices <- i
	}
	close(taskIndices)

	// Wait for all tasks to finish
	wg.Wait()
	return results
}
