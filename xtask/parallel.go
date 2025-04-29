package xtask

import (
	"fmt"
	"runtime/debug"

	"github.com/DreamvatLab/go/xdto"
)

// ParallelRun executes multiple functions in parallel and collects their results
// Parameters:
//   - actions: multiple functions, each returning a result and a possible error
//
// Returns:
//   - a slice of ChannelResult containing all function execution results
func ParallelRun(actions ...func() (interface{}, error)) []*xdto.ChannelResult {
	actionCount := len(actions)

	// If no functions are provided, return an empty slice
	if actionCount == 0 {
		return []*xdto.ChannelResult{}
	}

	// Create a slice of channels to collect results from all functions
	channels := make([]chan *xdto.ChannelResult, actionCount)

	// Create a channel for each function
	for i := range actions {
		channels[i] = make(chan *xdto.ChannelResult, 1)
	}

	// Start all functions
	for i, action := range actions {
		go func(index int, act func() (interface{}, error)) {
			// Use defer to ensure the channel is always closed
			defer close(channels[index])

			// Recover from any panic
			defer func() {
				if r := recover(); r != nil {
					// Create an error with the panic information
					err := fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack())

					// Send the error to the channel
					channels[index] <- &xdto.ChannelResult{
						Result: nil,
						Error:  err,
					}
				}
			}()

			// Execute the function and get the result
			result, err := act()

			// Send the result to the channel
			channels[index] <- &xdto.ChannelResult{
				Result: result,
				Error:  err,
			}
		}(i, action)
	}

	// Collect all results
	results := make([]*xdto.ChannelResult, actionCount)
	for i, ch := range channels {
		results[i] = <-ch
	}

	return results
}
