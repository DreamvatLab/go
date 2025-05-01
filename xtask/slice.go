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
//   - batchCallback: Optional function that will be called after each batch of tasks completes
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
//	}, func(batchIndex int, completedCount int) {
//	    fmt.Printf("Batch %d completed, total completed: %d\n", batchIndex, completedCount)
//	})
func ParallelRunSlice[T any](limit int, slice []T, processor func(T) (interface{}, error), batchCallback func(batchIndex, completedCount int)) []*TaskResult {
	if limit <= 0 {
		limit = runtime.GOMAXPROCS(0)
	}

	taskCount := len(slice)
	results := make([]*TaskResult, taskCount)
	var wg sync.WaitGroup
	taskIndices := make(chan int)
	completedTasks := make(chan int, taskCount)
	var callbackWg sync.WaitGroup

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
						completedTasks <- i
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

	// Start batch tracking goroutine if callback is provided
	if batchCallback != nil {
		callbackWg.Add(1)
		go func() {
			defer callbackWg.Done()
			completedCount := 0
			batchIndex := 0
			for range completedTasks {
				completedCount++
				if completedCount%limit == 0 || completedCount == taskCount {
					batchCallback(batchIndex, completedCount)
					batchIndex++
				}
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
	close(completedTasks)

	// Wait for callback goroutine to finish
	if batchCallback != nil {
		callbackWg.Wait()
	}

	return results
}
