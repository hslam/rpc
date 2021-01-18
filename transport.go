// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	//DefaultMaxConnsPerHost is the default value of Transport's MaxConnsPerHost.
	DefaultMaxConnsPerHost = 1
	//DefaultMaxIdleConnsPerHost is the default value of Transport's MaxIdleConnsPerHost.
	DefaultMaxIdleConnsPerHost = 1
	//DefaultKeepAlive is the default value of Transport's KeepAlive.
	DefaultKeepAlive = 90 * time.Second
	//DefaultIdleConnTimeout is the default value of Transport's IdleConnTimeout.
	DefaultIdleConnTimeout = 60 * time.Second

	defaultRunTicker = time.Second
)

var (
	// MaxConnsPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	MaxConnsPerHost = DefaultMaxConnsPerHost
	// MaxIdleConnsPerHost controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	// KeepAlive specifies the maximum amount of time keeping the active connections in the Transport's conns.
	KeepAlive = DefaultKeepAlive
	// IdleConnTimeout specifies the maximum amount of time keeping the idle connections  in the Transport's idleConns.
	IdleConnTimeout = DefaultIdleConnTimeout

	runTicker = defaultRunTicker

	// ErrDial is returned when dial failed.
	ErrDial = errors.New("dial failed")
)

// RoundTripper is an interface representing the ability to execute a
// single RPC transaction, obtaining the Response for a given Request.
type RoundTripper interface {
	RoundTrip(addr string, call *Call) *Call
	Go(addr, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
	Call(addr, serviceMethod string, args interface{}, reply interface{}) error
	CallWithContext(addr string, ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
	Watch(addr, key string) (Watcher, error)
	Ping(addr string) error
	Close() error
}

//Transport defines the struct of transport
type Transport struct {
	connsMu             sync.Mutex
	once                sync.Once
	idleConns           map[string]*connQueue
	conns               map[string]*conns
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	KeepAlive           time.Duration
	IdleConnTimeout     time.Duration
	now                 time.Time
	Network             string
	Codec               string
	Dial                func(network, address, codec string) (*Conn, error)
	Options             *Options
	DialWithOptions     func(address string, opts *Options) (*Conn, error)
	running             bool
	done                chan bool
	closed              uint32
	ticker              time.Duration
}

// DefaultTransport is a default RPC transport.
var DefaultTransport = &Transport{
	MaxConnsPerHost:     MaxConnsPerHost,
	MaxIdleConnsPerHost: MaxIdleConnsPerHost,
	KeepAlive:           KeepAlive,
	IdleConnTimeout:     IdleConnTimeout,
	Network:             "tcp",
	Codec:               "json",
	Dial:                Dial,
	DialWithOptions:     DialWithOptions,
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (t *Transport) RoundTrip(addr string, call *Call) *Call {
	done := call.Done
	done = checkDone(done)
	call.Done = done
	conn, err := t.getConn(addr)
	if err != nil {
		call.Error = err
		call.done()
		return call
	}
	conn.RoundTrip(call)
	conn.lastTime = t.now
	checkPersistConnErr(call.Error, conn)
	return call
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (t *Transport) Go(addr, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	conn, err := t.getConn(addr)
	if err != nil {
		call := new(Call)
		call.ServiceMethod = serviceMethod
		call.Args = args
		call.Reply = reply
		done = checkDone(done)
		call.Done = done
		call.Error = err
		call.done()
		return call
	}
	call := conn.Go(serviceMethod, args, reply, done)
	conn.lastTime = t.now
	checkPersistConnErr(call.Error, conn)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (t *Transport) Call(addr, serviceMethod string, args interface{}, reply interface{}) error {
	conn, err := t.getConn(addr)
	if err != nil {
		return err
	}
	err = conn.Call(serviceMethod, args, reply)
	conn.lastTime = t.now
	checkPersistConnErr(err, conn)
	return err
}

// CallWithContext acts like Call but takes a context.
func (t *Transport) CallWithContext(addr string, ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	conn, err := t.getConn(addr)
	if err != nil {
		return err
	}
	err = conn.CallWithContext(ctx, serviceMethod, args, reply)
	conn.lastTime = t.now
	checkPersistConnErr(err, conn)
	return err
}

// Watch returns the Watcher.
func (t *Transport) Watch(addr, key string) (Watcher, error) {
	conn, err := t.getConn(addr)
	if err != nil {
		return nil, err
	}
	watcher, err := conn.Watch(key)
	conn.lastTime = t.now
	checkPersistConnErr(err, conn)
	return watcher, err
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (t *Transport) Ping(addr string) error {
	conn, err := t.getConn(addr)
	if err != nil {
		return err
	}
	err = conn.Ping()
	conn.lastTime = t.now
	checkPersistConnErr(err, conn)
	return err
}

func checkPersistConnErr(err error, pc *persistConn) {
	if err == ErrShutdown {
		pc.mu.Lock()
		pc.alive = false
		pc.mu.Unlock()
		pc.Close()
	}
}

func (t *Transport) getConn(addr string) (pc *persistConn, err error) {
	if len(addr) == 0 {
		return nil, ErrDial
	}
	t.connsMu.Lock()
	defer t.connsMu.Unlock()
	if !t.running {
		t.once.Do(func() {
			t.idleConns = make(map[string]*connQueue)
			t.conns = make(map[string]*conns)
			t.done = make(chan bool, 10)
			if t.MaxConnsPerHost < 1 {
				t.MaxConnsPerHost = MaxConnsPerHost
			}
			if t.MaxIdleConnsPerHost < 1 {
				t.MaxIdleConnsPerHost = MaxIdleConnsPerHost
			} else if t.MaxIdleConnsPerHost > t.MaxConnsPerHost {
				t.MaxIdleConnsPerHost = t.MaxConnsPerHost
			}
			if t.KeepAlive <= 0 {
				t.KeepAlive = KeepAlive
			}
			if t.IdleConnTimeout <= 0 {
				t.IdleConnTimeout = IdleConnTimeout
			}
			t.now = time.Now()
			if t.Dial == nil {
				t.Dial = Dial
			}
			if t.DialWithOptions == nil {
				t.DialWithOptions = DialWithOptions
			}
			t.running = true
			go t.run()
		})
	}
	if cs, ok := t.conns[addr]; ok {
		if len(cs.Conns) < t.MaxConnsPerHost {
			if cq, ok := t.idleConns[addr]; ok && cq.Length() > 0 {
				pc = cq.Dequeue()
				pc.lastTime = t.now
				pc.mu.Lock()
				if !pc.alive {
					pc.mu.Unlock()
					if pc, err = t.newPersistConn(addr); err != nil {
						return nil, err
					}
				} else {
					pc.mu.Unlock()
				}
			} else if pc, err = t.newPersistConn(addr); err != nil {
				return nil, err
			}
			cs.Append(pc)
			t.conns[addr] = cs
			return pc, nil
		}
		cursor := cs.Cursor()
		pc = cs.Conns[cursor]
		pc.mu.Lock()
		if !pc.alive {
			pc.mu.Unlock()
			if pc, err = t.newPersistConn(addr); err != nil {
				return nil, err
			}
			cs.Conns[cursor] = pc
			return
		}
		pc.mu.Unlock()
		return
	}
	if cq, ok := t.idleConns[addr]; ok && cq.Length() > 0 {
		pc = cq.Dequeue()
		pc.lastTime = time.Now()
	} else {
		if pc, err = t.newPersistConn(addr); err != nil {
			return nil, err
		}
	}
	cs := &conns{addr: addr}
	cs.Append(pc)
	t.conns[addr] = cs
	return pc, nil
}

func (t *Transport) newPersistConn(addr string) (*persistConn, error) {
	var conn *Conn
	var err error
	if t.Options != nil {
		conn, err = t.DialWithOptions(addr, t.Options)
	} else {
		conn, err = t.Dial(t.Network, addr, t.Codec)
	}
	if err != nil {
		return nil, ErrDial
	}
	return &persistConn{
		Conn:     conn,
		alive:    true,
		lastTime: time.Now(),
	}, nil
}

func (t *Transport) run() {
	if t.ticker <= 0 {
		t.ticker = runTicker
	}
	ticker := time.NewTicker(t.ticker)
	for {
		select {
		case <-ticker.C:
			t.now = time.Now()
			t.connsMu.Lock()
			for _, cs := range t.conns {
				length := len(cs.Conns)
				for i := 0; i < length; i++ {
					pc := cs.Conns[i]
					if pc.lastTime.Add(t.KeepAlive).Before(time.Now()) && pc.NumCalls() == 0 {
						cs.Delete(i)
						i--
						length--
						if cq, ok := t.idleConns[cs.addr]; ok {
							if !t.idleConns[cs.addr].Enqueue(pc) {
								pc.Close()
							}
						} else {
							cq = newConnQueue(t.MaxIdleConnsPerHost, cs.addr)
							cq.Enqueue(pc)
							t.idleConns[cs.addr] = cq
						}
					} else {
						pc.Ping()
					}
				}
				if len(cs.Conns) == 0 {
					delete(t.conns, cs.addr)
				}
			}
			for _, cq := range t.idleConns {
				length := cq.Length()
				for i := 0; i < length; i++ {
					if cq.Rear().value.lastTime.Add(t.IdleConnTimeout).Before(time.Now()) {
						pc := cq.Dequeue()
						pc.Close()
					} else {
						cq.Rear().value.Ping()
					}
				}
				if cq.Length() == 0 {
					delete(t.idleConns, cq.addr)
				}
			}

			t.connsMu.Unlock()
		case <-t.done:
			ticker.Stop()
			return
		}
	}
}

// CloseIdleConnections closes the idle connections.
func (t *Transport) CloseIdleConnections() {
	t.connsMu.Lock()
	defer t.connsMu.Unlock()
	for _, cs := range t.conns {
		length := len(cs.Conns)
		for i := 0; i < length; i++ {
			pc := cs.Conns[i]
			if pc.NumCalls() == 0 {
				cs.Delete(i)
				i--
				length--
				pc.Close()
			}
		}
		if len(cs.Conns) == 0 {
			delete(t.conns, cs.addr)
		}
	}

	for _, cq := range t.idleConns {
		length := cq.Length()
		for i := 0; i < length; i++ {
			pc := cq.Dequeue()
			pc.Close()
		}
		delete(t.idleConns, cq.addr)
	}
}

// Close closes the all connections.
func (t *Transport) Close() error {
	if !atomic.CompareAndSwapUint32(&t.closed, 0, 1) {
		return nil
	}
	t.connsMu.Lock()
	defer t.connsMu.Unlock()
	if !t.running {
		return nil
	}
	for _, cs := range t.conns {
		length := len(cs.Conns)
		for i := 0; i < length; i++ {
			pc := cs.Conns[i]
			pc.Close()
		}
	}
	t.conns = make(map[string]*conns)
	for _, cq := range t.idleConns {
		length := cq.Length()
		if length > 0 {
			for i := 0; i < length; i++ {
				pc := cq.Dequeue()
				pc.Close()
			}
		}
		delete(t.idleConns, cq.addr)
	}
	t.idleConns = make(map[string]*connQueue)
	close(t.done)
	return nil
}

type persistConn struct {
	*Conn
	mu       sync.Mutex
	alive    bool
	lastTime time.Time
}

type conns struct {
	addr   string
	Conns  []*persistConn
	cursor int
}

func (c *conns) Cursor() int {
	c.cursor++
	if c.cursor > len(c.Conns)-1 {
		c.cursor = 0
	}
	return c.cursor
}

func (c *conns) Append(pc *persistConn) {
	c.Conns = append(c.Conns, pc)
}

func (c *conns) Delete(cursor int) {
	copy(c.Conns[cursor:], c.Conns[cursor+1:])
	c.Conns = c.Conns[:len(c.Conns)-1]
}

type node struct {
	value    *persistConn
	previous *node
	next     *node
}

type connQueue struct {
	addr     string
	front    *node
	rear     *node
	length   int
	capacity int
}

func newConnQueue(capacity int, addr string) *connQueue {
	front := &node{
		value:    nil,
		previous: nil,
	}

	rear := &node{
		value:    nil,
		previous: front,
	}

	front.next = rear
	return &connQueue{
		addr:     addr,
		front:    front,
		rear:     rear,
		capacity: capacity,
	}
}
func (q *connQueue) Length() int {
	return q.length
}

func (q *connQueue) Capacity() int {
	return q.capacity
}

func (q *connQueue) Front() *node {
	if q.length == 0 {
		return nil
	}
	return q.front.next
}

func (q *connQueue) Rear() *node {
	if q.length == 0 {
		return nil
	}
	return q.rear.previous
}

func (q *connQueue) Enqueue(value *persistConn) bool {
	if q.length == q.capacity || value == nil {
		return false
	}
	node := &node{
		value: value,
	}
	if q.length == 0 {
		q.front.next = node
	}
	node.previous = q.rear.previous
	node.next = q.rear
	q.rear.previous.next = node
	q.rear.previous = node
	q.length++
	return true
}

func (q *connQueue) Dequeue() *persistConn {
	if q.length == 0 {
		return nil
	}
	result := q.front.next
	q.front.next = result.next
	result.next = nil
	result.previous = nil
	q.length--
	return result.value
}
