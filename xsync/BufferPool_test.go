package xsync

import (
	"testing"
	"time"
)

func TestBufferPool_Basic(t *testing.T) {
	pool := NewSyncBufferPool(10)

	// Test Get and Put
	b1 := pool.GetBuffer()
	if b1.Cap() < 10 {
		t.Errorf("Expected capacity >= 10, got %d", b1.Cap())
	}

	// Write some data
	b1.WriteString("test")

	// Put back and get again
	pool.PutBuffer(b1)
	b2 := pool.GetBuffer()

	// Check if the buffer is empty
	if b2.Len() != 0 {
		t.Error("Buffer was not properly reset")
	}
}

func TestBufferPool_Concurrent(t *testing.T) {
	pool := NewSyncBufferPool(100)
	done := make(chan bool)

	// Run 10 goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				b := pool.GetBuffer()
				b.WriteString("test")
				time.Sleep(time.Millisecond)
				pool.PutBuffer(b)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestBufferPool_EdgeCases(t *testing.T) {
	pool := NewSyncBufferPool(10)

	// Test nil pointer
	pool.PutBuffer(nil) // Should not panic

	// Test zero size - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for zero size")
		}
	}()
	NewSyncBufferPool(0)

	// Test negative size - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for negative size")
		}
	}()
	NewSyncBufferPool(-1)
}
