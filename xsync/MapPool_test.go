package xsync

import (
	"strconv"
	"testing"
	"time"
)

func TestMapPool_Basic(t *testing.T) {
	pool := NewSyncMapPool()

	// Test Get and Put
	m1 := pool.GetMap()

	// Add some data
	(*m1)["key1"] = "value1"
	(*m1)["key2"] = "value2"

	// Put back and get again
	pool.PutMap(m1)
	m2 := pool.GetMap()

	// Check if the map is empty
	if len(*m2) != 0 {
		t.Error("Map was not properly cleared")
	}
}

func TestMapPool_Concurrent(t *testing.T) {
	pool := NewSyncMapPool()
	done := make(chan bool)

	// Run 10 goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				m := pool.GetMap()
				(*m)["key"] = "value"
				time.Sleep(time.Millisecond)
				pool.PutMap(m)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestMapPool_EdgeCases(t *testing.T) {
	pool := NewSyncMapPool()

	// Test nil pointer
	pool.PutMap(nil) // Should not panic

	// Test map with many entries
	m := pool.GetMap()
	for i := 0; i < 1000; i++ {
		(*m)[strconv.Itoa(i)] = i
	}

	pool.PutMap(m)
	m2 := pool.GetMap()

	// Check if the map is empty
	if len(*m2) != 0 {
		t.Error("Map with many entries was not properly cleared")
	}

	// Test nested maps
	m3 := pool.GetMap()
	(*m3)["nested"] = map[string]interface{}{
		"key": "value",
	}
	pool.PutMap(m3)
	m4 := pool.GetMap()
	if len(*m4) != 0 {
		t.Error("Map with nested maps was not properly cleared")
	}
}
