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
		want      []*TaskResult
	}{
		{
			name:  "empty slice",
			limit: 2,
			input: []int{},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			want: []*TaskResult{},
		},
		{
			name:  "single element",
			limit: 2,
			input: []int{1},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
			},
		},
		{
			name:  "multiple elements",
			limit: 2,
			input: []int{1, 2, 3, 4},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: 4, Error: nil},
				{Result: 6, Error: nil},
				{Result: 8, Error: nil},
			},
		},
		{
			name:  "with error",
			limit: 2,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				if n == 2 {
					return nil, errors.New("test error")
				}
				return n * 2, nil
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: nil, Error: errors.New("test error")},
				{Result: 6, Error: nil},
			},
		},
		{
			name:  "with panic",
			limit: 2,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				if n == 2 {
					panic("test panic")
				}
				return n * 2, nil
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: nil, Error: nil}, // Will be checked separately for panic
				{Result: 6, Error: nil},
			},
		},
		{
			name:  "zero limit",
			limit: 0,
			input: []int{1, 2},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: 4, Error: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParallelRunSlice(tt.limit, tt.input, tt.processor)

			if len(got) != len(tt.want) {
				t.Errorf("ParallelRunSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, result := range got {
				if tt.name == "with panic" && i == 1 {
					// Special case for panic test
					if result.Error == nil {
						t.Errorf("ParallelRunSlice() panic case = nil, want non-nil error")
					}
					continue
				}

				if result.Error != nil && tt.want[i].Error != nil {
					if result.Error.Error() != tt.want[i].Error.Error() {
						t.Errorf("ParallelRunSlice() error = %v, want %v", result.Error, tt.want[i].Error)
					}
				} else if (result.Error != nil) != (tt.want[i].Error != nil) {
					t.Errorf("ParallelRunSlice() error presence = %v, want %v", result.Error != nil, tt.want[i].Error != nil)
				}

				if result.Result != tt.want[i].Result {
					t.Errorf("ParallelRunSlice() result = %v, want %v", result.Result, tt.want[i].Result)
				}
			}
		})
	}
}

func TestParallelRunSliceWithBatchCallback(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		input         []int
		processor     func(int) (interface{}, error)
		onBatchDone   func([]*TaskResult, int, int) bool
		want          []*TaskResult
		expectBatches int
	}{
		{
			name:  "empty slice",
			limit: 2,
			input: []int{},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			onBatchDone: func(results []*TaskResult, batchIndex int, totalBatches int) bool {
				if totalBatches != 0 {
					t.Errorf("totalBatches = %v, want %v", totalBatches, 0)
				}
				return false
			},
			want:          []*TaskResult{},
			expectBatches: 0,
		},
		{
			name:  "single batch",
			limit: 4,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			onBatchDone: func(results []*TaskResult, batchIndex int, totalBatches int) bool {
				if totalBatches != 1 {
					t.Errorf("totalBatches = %v, want %v", totalBatches, 1)
				}
				return false
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: 4, Error: nil},
				{Result: 6, Error: nil},
			},
			expectBatches: 1,
		},
		{
			name:  "multiple batches",
			limit: 2,
			input: []int{1, 2, 3, 4},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			onBatchDone: func(results []*TaskResult, batchIndex int, totalBatches int) bool {
				if totalBatches != 2 {
					t.Errorf("totalBatches = %v, want %v", totalBatches, 2)
				}
				return false
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: 4, Error: nil},
				{Result: 6, Error: nil},
				{Result: 8, Error: nil},
			},
			expectBatches: 2,
		},
		{
			name:  "early stop",
			limit: 2,
			input: []int{1, 2, 3, 4},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			onBatchDone: func(results []*TaskResult, batchIndex int, totalBatches int) bool {
				if totalBatches != 2 {
					t.Errorf("totalBatches = %v, want %v", totalBatches, 2)
				}
				return batchIndex == 1 // Stop after second batch
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: 4, Error: nil},
				{Result: 6, Error: nil},
				{Result: 8, Error: nil},
			},
			expectBatches: 2,
		},
		{
			name:  "with error",
			limit: 2,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				if n == 2 {
					return nil, errors.New("test error")
				}
				return n * 2, nil
			},
			onBatchDone: func(results []*TaskResult, batchIndex int, totalBatches int) bool {
				if totalBatches != 2 {
					t.Errorf("totalBatches = %v, want %v", totalBatches, 2)
				}
				return false
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: nil, Error: errors.New("test error")},
				{Result: 6, Error: nil},
			},
			expectBatches: 2,
		},
		{
			name:  "with panic",
			limit: 2,
			input: []int{1, 2, 3},
			processor: func(n int) (interface{}, error) {
				if n == 2 {
					panic("test panic")
				}
				return n * 2, nil
			},
			onBatchDone: func(results []*TaskResult, batchIndex int, totalBatches int) bool {
				if totalBatches != 2 {
					t.Errorf("totalBatches = %v, want %v", totalBatches, 2)
				}
				return false
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: nil, Error: nil}, // Will be checked separately for panic
				{Result: 6, Error: nil},
			},
			expectBatches: 2,
		},
		{
			name:  "uneven batches",
			limit: 2,
			input: []int{1, 2, 3, 4, 5},
			processor: func(n int) (interface{}, error) {
				return n * 2, nil
			},
			onBatchDone: func(results []*TaskResult, batchIndex int, totalBatches int) bool {
				if totalBatches != 3 {
					t.Errorf("totalBatches = %v, want %v", totalBatches, 3)
				}
				return false
			},
			want: []*TaskResult{
				{Result: 2, Error: nil},
				{Result: 4, Error: nil},
				{Result: 6, Error: nil},
				{Result: 8, Error: nil},
				{Result: 10, Error: nil},
			},
			expectBatches: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParallelRunSliceWithBatchCallback(tt.limit, tt.input, tt.processor, tt.onBatchDone)

			if len(got) != len(tt.want) {
				t.Errorf("ParallelRunSliceWithBatchCallback() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, result := range got {
				if tt.name == "with panic" && i == 1 {
					// Special case for panic test
					if result.Error == nil {
						t.Errorf("ParallelRunSliceWithBatchCallback() panic case = nil, want non-nil error")
					}
					continue
				}

				if result.Error != nil && tt.want[i].Error != nil {
					if result.Error.Error() != tt.want[i].Error.Error() {
						t.Errorf("ParallelRunSliceWithBatchCallback() error = %v, want %v", result.Error, tt.want[i].Error)
					}
				} else if (result.Error != nil) != (tt.want[i].Error != nil) {
					t.Errorf("ParallelRunSliceWithBatchCallback() error presence = %v, want %v", result.Error != nil, tt.want[i].Error != nil)
				}

				if result.Result != tt.want[i].Result {
					t.Errorf("ParallelRunSliceWithBatchCallback() result = %v, want %v", result.Result, tt.want[i].Result)
				}
			}
		})
	}
}

// TestConcurrency tests the actual concurrency behavior
func TestConcurrency(t *testing.T) {
	const (
		limit     = 2
		taskCount = 4
		delay     = 100 * time.Millisecond
	)

	start := time.Now()
	results := ParallelRunSlice(limit, make([]int, taskCount), func(n int) (interface{}, error) {
		time.Sleep(delay)
		return n, nil
	})

	elapsed := time.Since(start)
	expectedMin := delay * (taskCount / limit)
	expectedMax := delay * (taskCount/limit + 1)

	if elapsed < expectedMin || elapsed > expectedMax {
		t.Errorf("ParallelRunSlice() took %v, expected between %v and %v", elapsed, expectedMin, expectedMax)
	}

	if len(results) != taskCount {
		t.Errorf("ParallelRunSlice() returned %d results, want %d", len(results), taskCount)
	}
}
