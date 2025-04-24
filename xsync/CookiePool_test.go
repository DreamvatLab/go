package xsync

import (
	"net/http"
	"testing"
	"time"
)

func TestCookiePool_Basic(t *testing.T) {
	pool := NewSyncCookiePool()

	// Test Get and Put
	c1 := pool.GetCookie()

	// Set some data
	c1.Name = "test"
	c1.Value = "value"
	c1.Domain = "example.com"

	// Put back and get again
	pool.PutCookie(c1)
	c2 := pool.GetCookie()

	// Check if the cookie is reset
	if c2.Name != "" || c2.Value != "" || c2.Domain != "" {
		t.Error("Cookie was not properly reset")
	}
}

func TestCookiePool_Concurrent(t *testing.T) {
	pool := NewSyncCookiePool()
	done := make(chan bool)

	// Run 10 goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				c := pool.GetCookie()
				c.Name = "test"
				c.Value = "value"
				time.Sleep(time.Millisecond)
				pool.PutCookie(c)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestCookiePool_EdgeCases(t *testing.T) {
	pool := NewSyncCookiePool()

	// Test nil pointer
	pool.PutCookie(nil) // Should not panic

	// Test all fields are reset
	c := pool.GetCookie()
	c.Name = "test"
	c.Value = "value"
	c.Domain = "example.com"
	c.Path = "/"
	c.Expires = time.Now()
	c.MaxAge = 100
	c.HttpOnly = true
	c.Secure = true
	c.SameSite = http.SameSiteStrictMode
	c.Raw = "raw"
	c.RawExpires = "raw-expires"
	c.Unparsed = []string{"test"}

	pool.PutCookie(c)
	c2 := pool.GetCookie()

	// Check all fields are reset
	if c2.Name != "" || c2.Value != "" || c2.Domain != "" || c2.Path != "" ||
		!c2.Expires.IsZero() || c2.MaxAge != 0 || c2.HttpOnly || c2.Secure ||
		c2.SameSite != 0 || c2.Raw != "" || c2.RawExpires != "" || len(c2.Unparsed) != 0 {
		t.Error("Not all cookie fields were properly reset")
	}
}
