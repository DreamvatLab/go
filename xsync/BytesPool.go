package xsync

import (
	"sync"
)

// IBytesPool defines the interface for a bytes slice pool.
// It provides methods to get and put byte slices back to the pool.
type IBytesPool interface {
	// GetBytes returns a byte slice from the pool.
	GetBytes() *[]byte
	// PutBytes puts a byte slice back to the pool.
	// The slice should not be used after being put back.
	PutBytes(*[]byte)
}

// syncBytesPool implements IBytesPool using sync.Pool
type syncBytesPool struct {
	pool *sync.Pool
	// initialSize is the size of new byte slices
	initialSize int
}

// NewSyncBytesPool creates a new bytes pool with the specified initial size.
// The buf_size parameter must be positive.
func NewSyncBytesPool(buf_size int) IBytesPool {
	if buf_size <= 0 {
		panic("buffer size must be positive")
	}

	return &syncBytesPool{
		pool: &sync.Pool{
			New: func() interface{} {
				b := make([]byte, buf_size)
				return &b
			},
		},
		initialSize: buf_size,
	}
}

// GetBytes returns a byte slice from the pool.
func (bp *syncBytesPool) GetBytes() *[]byte {
	return bp.pool.Get().(*[]byte)
}

// PutBytes puts a byte slice back to the pool.
// The slice should not be used after being put back.
func (bp *syncBytesPool) PutBytes(b *[]byte) {
	if b == nil {
		return
	}
	// Zero out the slice
	copy(*b, make([]byte, len(*b)))

	bp.pool.Put(b)
}
