package xsync

import (
	"bytes"
	"sync"
)

// IBufferPool defines the interface for a bytes.Buffer pool.
// It provides methods to get and put buffers back to the pool.
type IBufferPool interface {
	// GetBuffer returns a buffer from the pool.
	GetBuffer() *bytes.Buffer
	// PutBuffer puts a buffer back to the pool.
	// The buffer should not be used after being put back.
	PutBuffer(*bytes.Buffer)
}

// syncBufferPool implements IBufferPool using sync.Pool
type syncBufferPool struct {
	pool *sync.Pool
}

// NewSyncBufferPool creates a new buffer pool with the specified initial size.
// The buf_size parameter must be positive.
func NewSyncBufferPool(buf_size int) IBufferPool {
	if buf_size <= 0 {
		panic("buffer size must be positive")
	}

	return &syncBufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				var b bytes.Buffer
				b.Grow(buf_size)
				return &b
			},
		},
	}
}

// GetBuffer returns a buffer from the pool.
func (x *syncBufferPool) GetBuffer() *bytes.Buffer {
	return x.pool.Get().(*bytes.Buffer)
}

// PutBuffer puts a buffer back to the pool.
// The buffer should not be used after being put back.
func (bp *syncBufferPool) PutBuffer(b *bytes.Buffer) {
	if b == nil {
		return
	}
	b.Reset()
	bp.pool.Put(b)
}
