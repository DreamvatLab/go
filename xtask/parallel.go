package xtask

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// ParallelRun executes multiple tasks concurrently with a specified concurrency limit.
// It returns a slice of TaskResult containing the results and errors from each task.
// If limit is less than or equal to 0, it uses the number of available CPU cores as the limit.
func ParallelRun(limit int, tasks ...func() (interface{}, error)) []*TaskResult {
	// Set default concurrency limit to number of CPU cores if not specified
	if limit <= 0 {
		limit = runtime.GOMAXPROCS(0)
	}

	// Initialize result slice and synchronization primitives
	actionCount := len(tasks)
	results := make([]*TaskResult, actionCount)
	var wg sync.WaitGroup
	taskIndice := make(chan int) // Channel for distributing task indices to workers

	// Start worker goroutines
	for w := 0; w < limit; w++ {
		go func() {
			// Each worker processes tasks from the channel until it's closed
			for i := range taskIndice {
				func(index int) {
					// Ensure task completion is signaled and handle panics
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							// Store panic information in the result
							results[index] = &TaskResult{
								Result: nil,
								Error:  fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack()),
							}
						}
					}()

					// Execute the task and store its result
					result, err := tasks[index]()
					results[index] = &TaskResult{
						Result: result,
						Error:  err,
					}
				}(i)
			}
		}()
	}

	// Distribute tasks to workers
	for i := 0; i < actionCount; i++ {
		wg.Add(1)       // Increment wait group counter before sending task
		taskIndice <- i // Send task index to worker
	}
	close(taskIndice) // Signal that no more tasks will be sent

	// Wait for all tasks to complete
	wg.Wait()
	return results
}
