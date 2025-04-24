package xsync

import (
	"sync"
)

// IMapPool defines the interface for a map pool.
type IMapPool interface {
	GetMap() *map[string]interface{}
	PutMap(*map[string]interface{})
}

type syncMapPool struct {
	pool *sync.Pool
}

func NewSyncMapPool() IMapPool {
	return &syncMapPool{
		pool: &sync.Pool{
			New: func() interface{} {
				r := make(map[string]interface{})
				return &r
			},
		},
	}
}

func (x *syncMapPool) GetMap() *map[string]interface{} {
	return x.pool.Get().(*map[string]interface{})
}

func (x *syncMapPool) PutMap(m *map[string]interface{}) {
	if m == nil {
		return
	}
	// Clear all entries from the map
	for k := range *m {
		delete(*m, k)
	}
	x.pool.Put(m)
}
