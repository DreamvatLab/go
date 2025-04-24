package xsync

import (
	"net/http"
	"sync"
	"time"
)

// ICookiePool defines the interface for a http.Cookie pool.
// It provides methods to get and put cookies back to the pool.
type ICookiePool interface {
	// GetCookie returns a cookie from the pool.
	// The returned cookie will be reset to its zero values.
	GetCookie() *http.Cookie
	// PutCookie puts a cookie back to the pool.
	// The cookie should not be used after being put back.
	PutCookie(*http.Cookie)
}

// syncCookiePool implements ICookiePool using sync.Pool
type syncCookiePool struct {
	pool *sync.Pool
}

// NewSyncCookiePool creates a new cookie pool.
func NewSyncCookiePool() ICookiePool {
	return &syncCookiePool{
		pool: &sync.Pool{
			New: func() interface{} {
				return &http.Cookie{}
			},
		},
	}
}

// GetCookie returns a cookie from the pool.
func (x *syncCookiePool) GetCookie() *http.Cookie {
	return x.pool.Get().(*http.Cookie)
}

// PutCookie puts a cookie back to the pool.
// The cookie should not be used after being put back.
func (x *syncCookiePool) PutCookie(c *http.Cookie) {
	if c == nil {
		return
	}
	// Reset all fields to zero values
	c.Domain = ""
	c.Expires = time.Time{}
	c.HttpOnly = false
	c.MaxAge = 0
	c.Name = ""
	c.Path = ""
	c.Raw = ""
	c.RawExpires = ""
	c.SameSite = 0
	c.Secure = false
	c.Value = ""
	c.Unparsed = nil
	x.pool.Put(c)
}
