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
	reverseProxyAlpha   = 0.8
	reverseProxyTick    = time.Millisecond * 100
	reverseProxyLatency = int64(time.Minute)
	emptyString         = ""
)

// Scheduling represents the scheduling algorithms.
type Scheduling int

const (
	//RoundRobinScheduling uses the Round Robin algorithm to load balance traffic.
	RoundRobinScheduling Scheduling = iota
	//RandomScheduling randomly selects the target server.
	RandomScheduling
	//LeastTimeScheduling selects the target server with the lowest latency.
	LeastTimeScheduling
)

// ReverseProxy is an RPC Handler that takes an incoming request and
// sends it to another server, proxying the response back to the
// client.
type ReverseProxy struct {
	lock       sync.Mutex
	targets    map[string]*target
	list       []*target
	pos        int
	lastTime   time.Time
	Director   func() (target string)
	Transport  RoundTripper
	Scheduling Scheduling
	Tick       time.Duration
	Alpha      float64
}

// NewReverseProxy returns a new ReverseProxy that routes
// requests to the targets.
func NewReverseProxy(targets ...string) *ReverseProxy {
	if len(targets) > 0 {
		m := make(map[string]*target)
		l := make([]*target, 0, len(targets))
		for _, address := range targets {
			if len(address) > 0 {
				if _, ok := m[address]; !ok {
					t := &target{address: address, latency: reverseProxyLatency}
					m[address] = t
					l = append(l, t)
				}
			}
		}
		if len(l) > 0 {
			return &ReverseProxy{targets: m, list: l, Tick: reverseProxyTick, Alpha: reverseProxyAlpha}
		}
	}
	return &ReverseProxy{}
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (c *ReverseProxy) RoundTrip(call *Call) *Call {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().RoundTrip(address, call)
	}
	return c.transport().RoundTrip(target.address, call)
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *ReverseProxy) Call(serviceMethod string, args interface{}, reply interface{}) error {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().Call(address, serviceMethod, args, reply)
	}
	start := time.Now()
	err := c.transport().Call(target.address, serviceMethod, args, reply)
	c.update(target, int64(time.Now().Sub(start)))
	return err
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (c *ReverseProxy) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().Go(address, serviceMethod, args, reply, done)
	}
	return c.transport().Go(target.address, serviceMethod, args, reply, done)
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (c *ReverseProxy) Ping() error {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().Ping(address)
	}
	start := time.Now()
	err := c.transport().Ping(target.address)
	c.update(target, int64(time.Now().Sub(start)))
	return err
}

func (c *ReverseProxy) transport() RoundTripper {
	if c.Transport == nil {
		panic("The transport is nil")
	}
	return c.Transport
}

func (c *ReverseProxy) target() (string, *target) {
	if c.Director != nil {
		address := c.Director()
		if len(address) > 0 {
			return address, nil
		}
	}
	if len(c.list) == 1 {
		return c.list[0].address, nil
	}
	if len(c.list) > 1 {
		var t *target
		switch c.Scheduling {
		case RoundRobinScheduling:
			c.lock.Lock()
			t = c.list[c.pos]
			c.pos = (c.pos + 1) % len(c.list)
			c.lock.Unlock()
		case RandomScheduling:
			pos := rand.Intn(len(c.list))
			t = c.list[pos]
		case LeastTimeScheduling:
			now := time.Now()
			c.lock.Lock()
			if c.lastTime.Add(c.Tick).Before(now) {
				c.lastTime = now
				t = c.list[c.pos]
				c.pos = (c.pos + 1) % len(c.list)
			} else {
				minHeap(c.list)
				t = c.list[0]
			}
			c.lock.Unlock()
			return emptyString, t
		default:
			c.lock.Lock()
			t = c.list[c.pos]
			c.pos = (c.pos + 1) % len(c.list)
			c.lock.Unlock()
		}
		return t.address, nil
	}
	panic("The targets is nil")
}

func (c *ReverseProxy) update(t *target, new int64) {
	old := atomic.LoadInt64(&t.latency)
	if old >= reverseProxyLatency {
		atomic.StoreInt64(&t.latency, new)
	} else {
		atomic.StoreInt64(&t.latency, int64(float64(old)*c.Alpha+float64(new)*(1-c.Alpha)))
	}
}

type target struct {
	address string
	latency int64
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
