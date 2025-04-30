package xtask

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallelRun(t *testing.T) {
	// Test case 1: Normal execution with multiple functions
	t.Run("Normal execution", func(t *testing.T) {
		func1 := func() (interface{}, error) {
			return "result1", nil
		}
		func2 := func() (interface{}, error) {
			return 42, nil
		}
		func3 := func() (interface{}, error) {
			return true, nil
		}

		results := ParallelRun(3, func1, func2, func3)

		assert.Equal(t, 3, len(results))
		assert.Equal(t, "result1", results[0].Result)
		assert.Nil(t, results[0].Error)
		assert.Equal(t, 42, results[1].Result)
		assert.Nil(t, results[1].Error)
		assert.Equal(t, true, results[2].Result)
		assert.Nil(t, results[2].Error)
	})

	// Test case 2: Empty function list
	t.Run("Empty function list", func(t *testing.T) {
		results := ParallelRun(0)
		assert.Equal(t, 0, len(results))
	})

	// Test case 3: Functions with errors
	t.Run("Functions with errors", func(t *testing.T) {
		func1 := func() (interface{}, error) {
			return nil, errors.New("error1")
		}
		func2 := func() (interface{}, error) {
			return nil, errors.New("error2")
		}

		results := ParallelRun(2, func1, func2)

		assert.Equal(t, 2, len(results))
		assert.Nil(t, results[0].Result)
		assert.Equal(t, "error1", results[0].Error.Error())
		assert.Nil(t, results[1].Result)
		assert.Equal(t, "error2", results[1].Error.Error())
	})

	// Test case 4: Functions with panic
	t.Run("Functions with panic", func(t *testing.T) {
		func1 := func() (interface{}, error) {
			return "result1", nil
		}
		func2 := func() (interface{}, error) {
			panic("panic in func2")
		}
		func3 := func() (interface{}, error) {
			return 42, nil
		}

		results := ParallelRun(3, func1, func2, func3)

		assert.Equal(t, 3, len(results))
		assert.Equal(t, "result1", results[0].Result)
		assert.Nil(t, results[0].Error)
		assert.Nil(t, results[1].Result)
		assert.Contains(t, results[1].Error.Error(), "panic recovered: panic in func2")
		assert.Equal(t, 42, results[2].Result)
		assert.Nil(t, results[2].Error)
	})

	// Test case 5: Functions with different execution times
	t.Run("Functions with different execution times", func(t *testing.T) {
		func1 := func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return "result1", nil
		}
		func2 := func() (interface{}, error) {
			time.Sleep(50 * time.Millisecond)
			return 42, nil
		}
		func3 := func() (interface{}, error) {
			time.Sleep(150 * time.Millisecond)
			return true, nil
		}

		startTime := time.Now()
		results := ParallelRun(3, func1, func2, func3)
		elapsedTime := time.Since(startTime)

		// The total execution time should be close to the longest function execution time
		// (150ms) plus some overhead, but less than the sum of all execution times (300ms)
		assert.True(t, elapsedTime < 300*time.Millisecond)
		assert.True(t, elapsedTime > 150*time.Millisecond)

		assert.Equal(t, 3, len(results))
		assert.Equal(t, "result1", results[0].Result)
		assert.Nil(t, results[0].Error)
		assert.Equal(t, 42, results[1].Result)
		assert.Nil(t, results[1].Error)
		assert.Equal(t, true, results[2].Result)
		assert.Nil(t, results[2].Error)
	})
}
