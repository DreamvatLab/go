package xtask

import (
	"errors"
	"testing"
	"time"
)

func TestParallelRunSlice(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		input     []int
		processor func(int) (interface{}, error)
		want      []interface{}
		wantErr   bool
	}{
		{
			name:  "basic processing",
			limit: 2,
			input: []int{1, 2, 3, 4},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			want:    []interface{}{2, 4, 6, 8},
			wantErr: false,
		},
		{
			name:  "error handling",
			limit: 2,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				if n == 2 {
					return nil, errors.New("error on 2")
				}
				return n, nil
			},
			want:    []interface{}{1, nil, 3},
			wantErr: true,
		},
		{
			name:  "empty slice",
			limit: 2,
			input: []int{},
			processor: func(n int) (interface{}, error) {
				return n, nil
			},
			want:    []interface{}{},
			wantErr: false,
		},
		{
			name:  "zero limit",
			limit: 0,
			input: []int{1, 2},
			processor: func(n int) (interface{}, error) {
				return n, nil
			},
			want:    []interface{}{1, 2},
			wantErr: false,
		},
		{
			name:  "panic recovery",
			limit: 2,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				if n == 2 {
					panic("test panic")
				}
				return n, nil
			},
			want:    []interface{}{1, nil, 3},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := ParallelRunSlice(tt.limit, tt.input, tt.processor, nil)

			// Check length
			if len(results) != len(tt.input) {
				t.Errorf("ParallelRunSlice() returned %d results, want %d", len(results), len(tt.input))
				return
			}

			// Check results
			for i, result := range results {
				if tt.wantErr && result.Error != nil {
					continue // Error expected
				}

				if result.Error != nil {
					t.Errorf("ParallelRunSlice()[%d] error = %v, wantErr %v", i, result.Error, tt.wantErr)
					continue
				}

				if result.Result != tt.want[i] {
					t.Errorf("ParallelRunSlice()[%d] = %v, want %v", i, result.Result, tt.want[i])
				}
			}
		})
	}
}

func TestParallelRunSliceConcurrency(t *testing.T) {
	// Test to ensure concurrent execution
	limit := 2
	input := []int{0, 1, 2, 3} // Initialize with proper values
	processed := make(chan int, 4)
	done := make(chan struct{})

	// Start a goroutine to collect results
	go func() {
		processedCount := 0
		for n := range processed {
			processedCount++
			if processedCount > len(input) {
				t.Errorf("More elements processed than input size, got value %d", n)
			}
		}
		if processedCount != len(input) {
			t.Errorf("Expected %d elements to be processed, got %d", len(input), processedCount)
		}
		close(done)
	}()

	results := ParallelRunSlice(limit, input, func(n int) (interface{}, error) {
		time.Sleep(100 * time.Millisecond) // Simulate work
		processed <- n
		return n, nil
	}, nil)

	// Close the processed channel after all processing is done
	close(processed)

	// Wait for result collection with timeout
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out waiting for results")
	}

	// Check results
	for i, result := range results {
		if result.Error != nil {
			t.Errorf("Unexpected error at index %d: %v", i, result.Error)
		}
		if result.Result != i {
			t.Errorf("Expected result %d at index %d, got %v", i, i, result.Result)
		}
	}
}

func TestParallelRunSliceBatchCallback(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		input       []int
		processor   func(int) (interface{}, error)
		wantBatches int
	}{
		{
			name:  "exact batch size",
			limit: 2,
			input: []int{1, 2, 3, 4},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			wantBatches: 2,
		},
		{
			name:  "partial batch",
			limit: 3,
			input: []int{1, 2, 3, 4, 5},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			wantBatches: 2,
		},
		{
			name:  "single batch",
			limit: 5,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			wantBatches: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batches := 0
			completedCounts := make([]int, 0)
			batchCallback := func(batchIndex, completedCount int) {
				batches++
				completedCounts = append(completedCounts, completedCount)
			}

			results := ParallelRunSlice(tt.limit, tt.input, tt.processor, batchCallback)

			// Verify batch count
			if batches != tt.wantBatches {
				t.Errorf("Expected %d batches, got %d", tt.wantBatches, batches)
			}

			// Verify completed counts are in order
			for i := 1; i < len(completedCounts); i++ {
				if completedCounts[i] <= completedCounts[i-1] {
					t.Errorf("Completed counts should be increasing, got %v", completedCounts)
					break
				}
			}

			// Verify final completed count matches input length
			if len(completedCounts) > 0 && completedCounts[len(completedCounts)-1] != len(tt.input) {
				t.Errorf("Final completed count should be %d, got %d", len(tt.input), completedCounts[len(completedCounts)-1])
			}

			// Verify results
			for i, result := range results {
				if result.Error != nil {
					t.Errorf("Unexpected error at index %d: %v", i, result.Error)
				}
				if result.Result != tt.input[i]*2 {
					t.Errorf("Expected result %d at index %d, got %v", tt.input[i]*2, i, result.Result)
				}
			}
		})
	}
}

func TestParallelRunSliceBatchCallbackWithErrors(t *testing.T) {
	input := []int{1, 2, 3, 4}
	limit := 2
	batches := 0
	completedCounts := make([]int, 0)

	batchCallback := func(batchIndex, completedCount int) {
		batches++
		completedCounts = append(completedCounts, completedCount)
	}

	results := ParallelRunSlice(limit, input, func(n int) (interface{}, error) {
		if n == 2 {
			return nil, errors.New("test error")
		}
		return n * 2, nil
	}, batchCallback)

	// Verify batch count (should still get callbacks even with errors)
	if batches != 2 {
		t.Errorf("Expected 2 batches, got %d", batches)
	}

	// Verify completed counts
	if len(completedCounts) != 2 {
		t.Errorf("Expected 2 completed counts, got %d", len(completedCounts))
	}

	// Verify final completed count
	if completedCounts[len(completedCounts)-1] != len(input) {
		t.Errorf("Final completed count should be %d, got %d", len(input), completedCounts[len(completedCounts)-1])
	}

	// Verify results
	for i, result := range results {
		if i == 1 { // index of value 2
			if result.Error == nil {
				t.Error("Expected error at index 1")
			}
		} else {
			if result.Error != nil {
				t.Errorf("Unexpected error at index %d: %v", i, result.Error)
			}
			if result.Result != input[i]*2 {
				t.Errorf("Expected result %d at index %d, got %v", input[i]*2, i, result.Result)
			}
		}
	}
}
