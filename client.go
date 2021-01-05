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
	clientAlpha   = 0.8
	clientTick    = time.Millisecond * 100
	clientLatency = int64(time.Minute)
	emptyString   = ""
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

// Client is an RPC client.
//
// The Client's Transport typically has internal state (cached connections),
// so Clients should be reused instead of created as
// needed. Clients are safe for concurrent use by multiple goroutines.
//
// A Client is higher-level than a RoundTripper such as Transport.
type Client struct {
	lock       sync.Mutex
	targets    map[string]*target
	list       []*target
	minHeap    []*target
	pos        int
	lastTime   time.Time
	Director   func() (target string)
	Transport  RoundTripper
	Scheduling Scheduling
	Tick       time.Duration
	Alpha      float64
}

// NewClient returns a new RPC Client.
func NewClient(opts *Options, targets ...string) *Client {
	client := &Client{Tick: clientTick, Alpha: clientAlpha}
	if opts != nil {
		if opts.NewCodec == nil && opts.NewHeaderEncoder == nil && opts.Codec == "" {
			panic("need opts.NewCodec, opts.NewHeaderEncoder or opts.Codec")
		}
		if opts.NewSocket == nil && opts.Network == "" {
			panic("need opts.NewSocket or opts.Network")
		}
		client.Transport = &Transport{Options: opts}
	}
	if len(targets) > 0 {
		client.Update(targets...)
	}
	return client
}

// Update updates targets.
func (c *Client) Update(targets ...string) {
	m := make(map[string]*target)
	l := make([]*target, 0, len(targets))
	for _, address := range targets {
		if len(address) > 0 {
			if _, ok := m[address]; !ok {
				t := &target{address: address, latency: clientLatency}
				m[address] = t
				l = append(l, t)
			}
		}
	}
	c.lock.Lock()
	c.targets = m
	c.list = l
	c.lock.Unlock()
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (c *Client) RoundTrip(call *Call) *Call {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().RoundTrip(address, call)
	}
	return c.transport().RoundTrip(target.address, call)
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().Call(address, serviceMethod, args, reply)
	}
	start := time.Now()
	err := c.transport().Call(target.address, serviceMethod, args, reply)
	c.update(target, int64(time.Now().Sub(start)), err)
	return err
}

// CallTimeout acts like Call but takes a timeout.
func (c *Client) CallTimeout(serviceMethod string, args interface{}, reply interface{}, timeout time.Duration) error {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().CallTimeout(address, serviceMethod, args, reply, timeout)
	}
	start := time.Now()
	err := c.transport().CallTimeout(target.address, serviceMethod, args, reply, timeout)
	c.update(target, int64(time.Now().Sub(start)), err)
	return err
}

// Watch returns the Watcher.
func (c *Client) Watch(key string) (Watcher, error) {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().Watch(address, key)
	}
	start := time.Now()
	watcher, err := c.transport().Watch(target.address, key)
	c.update(target, int64(time.Now().Sub(start)), err)
	return watcher, err
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (c *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().Go(address, serviceMethod, args, reply, done)
	}
	return c.transport().Go(target.address, serviceMethod, args, reply, done)
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (c *Client) Ping() error {
	address, target := c.target()
	if len(address) > 0 {
		return c.transport().Ping(address)
	}
	start := time.Now()
	err := c.transport().Ping(target.address)
	c.update(target, int64(time.Now().Sub(start)), err)
	return err
}

// Close closes the all connections.
func (c *Client) Close() (err error) {
	if c.Transport != nil {
		err = c.Transport.Close()
	}
	return
}

func (c *Client) transport() RoundTripper {
	if c.Transport == nil {
		panic("The transport is nil")
	}
	return c.Transport
}

func (c *Client) target() (string, *target) {
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
		var pos int
		switch c.Scheduling {
		case RoundRobinScheduling:
			c.lock.Lock()
			t = c.list[c.pos]
			pos = c.pos
			c.pos = (c.pos + 1) % len(c.list)
			c.lock.Unlock()
		case RandomScheduling:
			pos = rand.Intn(len(c.list))
			t = c.list[pos]
		case LeastTimeScheduling:
			now := time.Now()
			c.lock.Lock()
			if c.lastTime.Add(c.Tick).Before(now) {
				c.lastTime = now
				t = c.list[c.pos]
				pos = c.pos
				c.pos = (c.pos + 1) % len(c.list)
			} else {
				if len(c.minHeap) == 0 {
					c.minHeap = make([]*target, len(c.list))
					copy(c.minHeap, c.list)
				}
				minHeap(c.minHeap)
				t = c.minHeap[0]
			}
			c.lock.Unlock()
		default:
			c.lock.Lock()
			t = c.list[c.pos]
			pos = c.pos
			c.pos = (c.pos + 1) % len(c.list)
			c.lock.Unlock()
		}
		cnt := len(c.list)
	retry:
		cnt--
		if !t.alive && cnt > 0 {
			if !c.check(t) {
				pos = (pos + 1) % len(c.list)
				t = c.list[pos]
				goto retry
			}
		}
		return emptyString, t
	}
	panic("The targets is nil")
}

func (c *Client) check(t *target) bool {
	if c.Transport == nil {
		return false
	}
	return c.alive(t, c.transport().Ping(t.address))
}

func (c *Client) update(t *target, new int64, err error) {
	old := atomic.LoadInt64(&t.latency)
	if !c.alive(t, err) {
		atomic.StoreInt64(&t.latency, clientLatency)
	} else if old >= clientLatency {
		atomic.StoreInt64(&t.latency, new)
	} else {
		atomic.StoreInt64(&t.latency, int64(float64(old)*c.Alpha+float64(new)*(1-c.Alpha)))
	}
}

func (c *Client) alive(t *target, err error) bool {
	if err == ErrDial {
		t.alive = false
	} else {
		t.alive = true
	}
	return t.alive
}

type target struct {
	address string
	latency int64
	alive   bool
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
