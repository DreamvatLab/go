package xsync

import (
	"testing"
	"time"
)

func TestBytesPool_Basic(t *testing.T) {
	pool := NewSyncBytesPool(10)

	// Test Get and Put
	b1 := pool.GetBytes()
	if len(*b1) != 10 {
		t.Errorf("Expected length 10, got %d", len(*b1))
	}

	// Write some data
	(*b1)[0] = 1
	(*b1)[1] = 2

	// Put back and get again
	pool.PutBytes(b1)
	b2 := pool.GetBytes()

	// Check if the slice is zeroed
	if (*b2)[0] != 0 || (*b2)[1] != 0 {
		t.Error("Slice was not properly zeroed")
	}
}

func TestBytesPool_Concurrent(t *testing.T) {
	pool := NewSyncBytesPool(100)
	done := make(chan bool)

	// Run 10 goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				b := pool.GetBytes()
				// Write some data
				(*b)[0] = byte(j)
				time.Sleep(time.Millisecond)
				pool.PutBytes(b)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestBytesPool_EdgeCases(t *testing.T) {
	pool := NewSyncBytesPool(10)

	// Test nil pointer
	pool.PutBytes(nil) // Should not panic

	// Test zero size - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for zero size")
		}
	}()
	NewSyncBytesPool(0)

	// Test negative size - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for negative size")
		}
	}()
	NewSyncBytesPool(-1)
}
