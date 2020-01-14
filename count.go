package rpc

import (
	"sync/atomic"
)

type Count struct {
	v int64
}

func (c *Count) add(delta int64) int64 {
	return atomic.AddInt64(&c.v, delta)
}

func (c *Count) load() int64 {
	return atomic.LoadInt64(&c.v)
}
