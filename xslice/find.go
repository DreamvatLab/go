package xslice

import (
	"runtime"
)

// FindItem searches for the first element in the slice that matches the given predicate function.
// For large datasets (>= 10,000 elements), it automatically switches to concurrent search.
// Returns the index and value of the found element, or (-1, zero value) if not found.
func FindItem[T any](slice []T, fnFilter func(T) bool) (int, T) {
	const threshold = 10_000

	if len(slice) < threshold {
		return FindItemLinear(slice, fnFilter)
	}
	return FindItemConcurrent(slice, fnFilter)
}

// FindItemLinear performs a sequential search through the slice.
// Returns the index and value of the first matching element, or (-1, zero value) if not found.
func FindItemLinear[T any](slice []T, fnFilter func(T) bool) (int, T) {
	for i := 0; i < len(slice); i++ {
		if fnFilter(slice[i]) {
			return i, slice[i]
		}
	}
	var zero T
	return -1, zero
}

// FindItemConcurrent performs a parallel search using multiple goroutines.
// The number of workers is determined by runtime.GOMAXPROCS.
// Falls back to linear search if GOMAXPROCS < 2.
// Returns the index and value of the first matching element, or (-1, zero value) if not found.
func FindItemConcurrent[T any](slice []T, fnFilter func(T) bool) (int, T) {
	Workers := runtime.GOMAXPROCS(0)
	if Workers < 2 {
		return FindItemLinear(slice, fnFilter)
	}

	type Result struct {
		Index int
		Value T
	}

	ch := make(chan Result, Workers)
	ChunkSize := (len(slice) + Workers - 1) / Workers

	for W := 0; W < Workers; W++ {
		Start := W * ChunkSize
		End := Start + ChunkSize
		if End > len(slice) {
			End = len(slice)
		}

		go func(Start, End int) {
			for I := Start; I < End; I++ {
				if fnFilter(slice[I]) {
					ch <- Result{I, slice[I]}
					return
				}
			}
			ch <- Result{-1, *new(T)}
		}(Start, End)
	}

	for I := 0; I < Workers; I++ {
		res := <-ch
		if res.Index != -1 {
			return res.Index, res.Value
		}
	}

	var Zero T
	return -1, Zero
}
