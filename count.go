package rpc

import (
	"sync/atomic"
)

type count struct {
	v int64
}

func (c *count) add(delta int64) int64 {
	return atomic.AddInt64(&c.v, delta)
}

func (c *count) load() int64 {
	return atomic.LoadInt64(&c.v)
}
