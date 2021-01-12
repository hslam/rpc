// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	clientAlpha   = 0.8
	clientTick    = time.Millisecond * 100
	dialTimeout   = time.Minute
	clientLatency = int64(dialTimeout)
	emptyString   = ""
)

var errTarget = errors.New("There is no alive target")

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

type waiter struct {
	seq  uint64
	err  error
	Done chan *waiter
}

func (w *waiter) done() {
	select {
	case w.Done <- w:
	default:
	}
}

// Client is an RPC client.
//
// The Client's Transport typically has internal state (cached connections),
// so Clients should be reused instead of created as
// needed. Clients are safe for concurrent use by multiple goroutines.
//
// A Client is higher-level than a RoundTripper such as Transport.
type Client struct {
	lock        sync.Mutex
	targets     map[string]*target
	list        []*target
	minHeap     []*target
	last        []string
	seq         uint64
	pending     map[uint64]*waiter
	alives      uint32
	pos         int
	lastTime    time.Time
	Director    func() (target string)
	Transport   RoundTripper
	Scheduling  Scheduling
	Tick        time.Duration
	Alpha       float64
	DialTimeout time.Duration
	waiterPool  *sync.Pool
	donePool    *sync.Pool
	done        chan struct{}
	closed      uint32
}

// NewClient returns a new RPC Client.
func NewClient(opts *Options, targets ...string) *Client {
	client := &Client{
		Tick:        clientTick,
		Alpha:       clientAlpha,
		DialTimeout: dialTimeout,
		done:        make(chan struct{}, 1),
		pending:     make(map[uint64]*waiter),
		waiterPool:  &sync.Pool{New: func() interface{} { return &waiter{} }},
		donePool:    &sync.Pool{New: func() interface{} { return make(chan *waiter, 10) }},
	}
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
	go client.run()
	return client
}

// Update updates targets.
func (c *Client) Update(targets ...string) {
	m := make(map[string]*target)
	for _, address := range targets {
		if len(address) > 0 {
			if _, ok := m[address]; !ok {
				m[address] = &target{address: address, latency: clientLatency}
			}
		}
	}
	c.lock.Lock()
	c.targets = m
	c.list = list{}
	c.minHeap = list{}
	c.last = []string{}
	c.lock.Unlock()
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (c *Client) RoundTrip(call *Call) *Call {
	address, target, err := c.director()
	if err != nil {
		return c.transport().RoundTrip("", call)
	}
	if len(address) > 0 {
		return c.transport().RoundTrip(address, call)
	}
	return c.transport().RoundTrip(target.address, call)
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	address, target, err := c.director()
	if err != nil {
		return err
	}
	if len(address) > 0 {
		return c.transport().Call(address, serviceMethod, args, reply)
	}
	start := time.Now()
	err = c.transport().Call(target.address, serviceMethod, args, reply)
	target.Update(c.Alpha, int64(time.Now().Sub(start)), err)
	return err
}

// CallTimeout acts like Call but takes a timeout.
func (c *Client) CallTimeout(serviceMethod string, args interface{}, reply interface{}, timeout time.Duration) error {
	address, target, err := c.director()
	if err != nil {
		return err
	}
	if len(address) > 0 {
		return c.transport().CallTimeout(address, serviceMethod, args, reply, timeout)
	}
	start := time.Now()
	err = c.transport().CallTimeout(target.address, serviceMethod, args, reply, timeout)
	target.Update(c.Alpha, int64(time.Now().Sub(start)), err)
	return err
}

// Watch returns the Watcher.
func (c *Client) Watch(key string) (Watcher, error) {
	address, target, err := c.director()
	if err != nil {
		return c.transport().Watch("", key)
	}
	if len(address) > 0 {
		return c.transport().Watch(address, key)
	}
	start := time.Now()
	watcher, err := c.transport().Watch(target.address, key)
	target.Update(c.Alpha, int64(time.Now().Sub(start)), err)
	return watcher, err
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (c *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	address, target, err := c.director()
	if err != nil {
		return c.transport().Go("", serviceMethod, args, reply, done)
	}
	if len(address) > 0 {
		return c.transport().Go(address, serviceMethod, args, reply, done)
	}
	return c.transport().Go(target.address, serviceMethod, args, reply, done)
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (c *Client) Ping() error {
	address, target, err := c.director()
	if err != nil {
		return c.transport().Ping("")
	}
	if len(address) > 0 {
		return c.transport().Ping(address)
	}
	start := time.Now()
	err = c.transport().Ping(target.address)
	target.Update(c.Alpha, int64(time.Now().Sub(start)), err)
	return err
}

// Close closes the all connections.
func (c *Client) Close() (err error) {
	if c.Transport != nil {
		err = c.Transport.Close()
	}
	if atomic.CompareAndSwapUint32(&c.closed, 0, 1) {
		if c.done != nil {
			close(c.done)
		}
		c.lock.Lock()
		for seq, w := range c.pending {
			delete(c.pending, seq)
			w.err = ErrShutdown
			w.done()
		}
		c.lock.Unlock()
	}
	return
}

func (c *Client) transport() RoundTripper {
	if c.Transport == nil {
		panic("The transport is nil")
	}
	return c.Transport
}

func (c *Client) director() (address string, t *target, err error) {
	if c.Director != nil {
		address = c.Director()
		if len(address) > 0 {
			return address, nil, nil
		}
	}
	c.lock.Lock()
	if len(c.list) > 0 {
		address, t, err = c.schedule()
		c.lock.Unlock()
		return
	}
	c.lock.Unlock()
	done := c.donePool.Get().(chan *waiter)
	w := c.waiterPool.Get().(*waiter)
	w.Done = done
	c.wait(w)
	timer := time.NewTimer(c.DialTimeout)
	runtime.Gosched()
	select {
	case <-w.Done:
		timer.Stop()
		resetWaiterDone(done)
		c.donePool.Put(done)
		err = w.err
		*w = waiter{}
		c.waiterPool.Put(w)
		if err == nil {
			c.lock.Lock()
			address, t, err = c.schedule()
			c.lock.Unlock()
		}
	case <-timer.C:
		resetWaiterDone(done)
		c.donePool.Put(done)
		seq := w.seq
		*w = waiter{}
		c.waiterPool.Put(w)
		err = ErrTimeout
		c.lock.Lock()
		delete(c.pending, seq)
		c.lock.Unlock()
	}
	return
}

func (c *Client) wait(s *waiter) {
	if atomic.LoadUint32(&c.closed) == 1 {
		s.err = ErrShutdown
		s.done()
		return
	}
	c.lock.Lock()
	s.seq = c.seq
	c.seq++
	c.pending[s.seq] = s
	c.lock.Unlock()
}

func (c *Client) schedule() (string, *target, error) {
	if len(c.list) == 1 {
		return c.list[0].address, nil, nil
	}
	if len(c.list) > 1 {
		var t *target
		switch c.Scheduling {
		case RoundRobinScheduling:
			t = c.list[c.pos]
			c.pos = (c.pos + 1) % len(c.list)
		case RandomScheduling:
			pos := rand.Intn(len(c.list))
			t = c.list[pos]
		case LeastTimeScheduling:
			now := time.Now()
			if c.lastTime.Add(c.Tick).Before(now) {
				c.lastTime = now
				t = c.list[c.pos]
				c.pos = (c.pos + 1) % len(c.list)
			} else {
				minHeap(c.minHeap)
				t = c.minHeap[0]
			}
		default:
			t = c.list[c.pos]
			c.pos = (c.pos + 1) % len(c.list)
		}
		return emptyString, t, nil
	}
	return emptyString, nil, ErrDial
}

func (c *Client) run() {
	ticker := time.NewTicker(clientTick)
	for {
		c.detect()
		select {
		case <-ticker.C:
		case <-c.done:
			ticker.Stop()
			return
		}
	}
}

func (c *Client) detect() {
	c.lock.Lock()
	for _, t := range c.targets {
		if t.alive == false {
			go c.check(t)
		}
	}
	c.lock.Unlock()
}

func (c *Client) check(t *target) (alive bool) {
	if c.Transport == nil {
		return false
	}
	err := c.transport().Ping(t.address)
	c.lock.Lock()
	alive = t.Alive(err)
	var l = list{}
	var addrs = []string{}
	for _, t := range c.targets {
		if t.alive {
			l = append(l, t)
			addrs = append(addrs, t.address)
		} else {
			t.Update(c.Alpha, clientLatency, errTarget)
		}
	}
	if len(l) > 0 {
		sort.Strings(addrs)
		if !reflect.DeepEqual(addrs, c.last) {
			c.last = addrs
			minHeap := make(list, len(l))
			copy(minHeap, l)
			c.list = l
			c.minHeap = minHeap
			c.pos = 0
		}
		for seq, w := range c.pending {
			delete(c.pending, seq)
			w.done()
		}
	} else {
		c.list = list{}
		c.minHeap = list{}
		c.last = []string{}
	}
	c.lock.Unlock()
	return
}

type target struct {
	address string
	latency int64
	alive   bool
}

func (t *target) Update(alpha float64, new int64, err error) {
	old := atomic.LoadInt64(&t.latency)
	if !t.Alive(err) {
		atomic.StoreInt64(&t.latency, clientLatency)
	} else if old >= clientLatency {
		atomic.StoreInt64(&t.latency, new)
	} else {
		atomic.StoreInt64(&t.latency, int64(float64(old)*alpha+float64(new)*(1-alpha)))
	}
}

func (t *target) Alive(err error) bool {
	if err == ErrDial {
		t.alive = false
	} else {
		t.alive = true
	}
	return t.alive
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

func resetWaiterDone(done chan *waiter) {
	for len(done) > 0 {
		onceWaiterDone(done)
	}
}

func onceWaiterDone(done chan *waiter) {
	select {
	case <-done:
	default:
	}
}
