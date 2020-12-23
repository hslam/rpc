// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	reverseProxyAlpha   = 0.2
	reverseProxyTick    = time.Millisecond * 100
	reverseProxyLatency = int64(time.Minute)
)

// Select represents the selecting algorithms.
type Select int

const (
	//RoundRobin uses the Round Robin algorithm to load balance traffic.
	RoundRobin Select = iota
	//Random randomly selects the target server.
	Random
	//LeastTime selects the target server with the lowest latency.
	LeastTime
)

// ReverseProxy is an RPC Handler that takes an incoming request and
// sends it to another server, proxying the response back to the
// client.
type ReverseProxy struct {
	lock      sync.Mutex
	targets   []*target
	pos       int
	lastTime  time.Time
	Tick      time.Duration
	Select    Select
	Transport RoundTripper
}

// NewReverseProxy returns a new ReverseProxy that routes
// requests to the targets.
func NewReverseProxy(targets ...string) *ReverseProxy {
	if len(targets) == 0 {
		panic("The targets is nil")
	}
	l := make([]*target, len(targets))
	for i := 0; i < len(targets); i++ {
		l[i] = &target{address: targets[i], latency: reverseProxyLatency}
	}
	return &ReverseProxy{targets: l, Tick: reverseProxyTick}
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (c *ReverseProxy) RoundTrip(call *Call) *Call {
	return c.transport().RoundTrip(c.target().address, call)
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *ReverseProxy) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if c.Select == LeastTime && len(c.targets) > 1 {
		start := time.Now()
		t := c.target()
		err := c.transport().Call(t.address, serviceMethod, args, reply)
		t.update(int64(time.Now().Sub(start)))
		return err
	}
	return c.transport().Call(c.target().address, serviceMethod, args, reply)
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (c *ReverseProxy) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	return c.transport().Go(c.target().address, serviceMethod, args, reply, done)
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (c *ReverseProxy) Ping() error {
	if c.Select == LeastTime && len(c.targets) > 1 {
		start := time.Now()
		t := c.target()
		err := c.transport().Ping(t.address)
		t.update(int64(time.Now().Sub(start)))
		return err
	}
	return c.transport().Ping(c.target().address)
}

func (c *ReverseProxy) transport() RoundTripper {
	if c.Transport == nil {
		panic("The transport is nil")
	}
	return c.Transport
}

func (c *ReverseProxy) target() *target {
	if len(c.targets) == 1 {
		return c.targets[0]
	}
	if len(c.targets) > 1 {
		var t *target
		switch c.Select {
		case RoundRobin:
			c.lock.Lock()
			t = c.targets[c.pos]
			c.pos = (c.pos + 1) % len(c.targets)
			c.lock.Unlock()
		case Random:
			pos := rand.Intn(len(c.targets))
			t = c.targets[pos]
		case LeastTime:
			now := time.Now()
			c.lock.Lock()
			if c.lastTime.Add(c.Tick).Before(now) {
				c.lastTime = now
				t = c.targets[c.pos]
				c.pos = (c.pos + 1) % len(c.targets)
			} else {
				minHeap(c.targets)
				t = c.targets[0]
			}
			c.lock.Unlock()
		}
		return t
	}
	panic("The targets is nil")
}

type target struct {
	address string
	latency int64
}

func (t *target) update(new int64) {
	old := atomic.LoadInt64(&t.latency)
	if old >= reverseProxyLatency {
		atomic.StoreInt64(&t.latency, new)
	} else {
		atomic.StoreInt64(&t.latency, int64(float64(new)*reverseProxyAlpha+float64(old)*(1-reverseProxyAlpha)))
	}
}

type list []*target

func (l list) Len() int { return len(l) }
func (l list) Less(i, j int) bool {
	return atomic.LoadInt64(&l[i].latency) < atomic.LoadInt64(&l[j].latency)
}
func (l list) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func minHeap(h list) {
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		heapDown(h, i, n)
	}
}

func heapDown(h list, i, n int) bool {
	parent := i
	for {
		leftChild := 2*parent + 1
		if leftChild >= n || leftChild < 0 { // leftChild < 0 after int overflow
			break
		}
		lessChild := leftChild
		if rightChild := leftChild + 1; rightChild < n && h.Less(rightChild, leftChild) {
			lessChild = rightChild
		}
		if !h.Less(lessChild, parent) {
			break
		}
		h.Swap(parent, lessChild)
		parent = lessChild
	}
	return parent > i
}
