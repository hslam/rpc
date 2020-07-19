package rpc

import (
	"sync"
	"time"
)

const (
	//DefaultMaxConnsPerHost defines the max conn per host
	DefaultMaxConnsPerHost = 1
	//DefaultMaxIdleConnsPerHost is the default value of Transport's MaxIdleConnsPerHost.
	DefaultMaxIdleConnsPerHost = 1
	DefaultKeepAlive           = 90 * time.Second
	DefaultIdleConnTimeout     = 60 * time.Second
)

var (
	//MaxConnsPerHost defines the max conn per host
	MaxConnsPerHost = DefaultMaxConnsPerHost
	//MaxIdleConnsPerHost is the default value of Transport's MaxIdleConnsPerHost.
	MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	KeepAlive           = DefaultKeepAlive
	IdleConnTimeout     = DefaultIdleConnTimeout
)

type RoundTripper interface {
	Call(addr, serviceMethod string, args interface{}, reply interface{}) error
	Go(addr, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
	Ping(addr string) error
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
	Dial                func(network, address, codec string) (*Client, error)
	Options             *Options
	DialWithOptions     func(address string, opts *Options) (*Client, error)
	running             bool
	done                chan bool
}

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

func (t *Transport) Call(addr, serviceMethod string, args interface{}, reply interface{}) error {
	client, err := t.getConn(addr)
	if err != nil {
		return err
	}
	err = client.Call(serviceMethod, args, reply)
	client.lastTime = t.now
	if err == ErrShutdown {
		client.mu.Lock()
		client.alive = false
		client.mu.Unlock()
		client.Close()
	}
	return err
}

func (t *Transport) Go(addr, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	client, err := t.getConn(addr)
	if err != nil {
		return nil
	}
	call := client.Go(serviceMethod, args, reply, done)
	client.lastTime = t.now
	if err == ErrShutdown {
		client.mu.Lock()
		client.alive = false
		client.mu.Unlock()
		client.Close()
	}
	return call
}

func (t *Transport) Ping(addr string) error {
	client, err := t.getConn(addr)
	if err != nil {
		return err
	}
	err = client.Ping()
	client.lastTime = t.now
	if err != nil {
		client.mu.Lock()
		client.alive = false
		client.mu.Unlock()
		client.Close()
	}
	return err
}

func (t *Transport) getConn(addr string) (pc *persistConn, err error) {
	t.connsMu.Lock()
	defer t.connsMu.Unlock()
	if !t.running {
		t.once.Do(func() {
			t.running = true
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
			go t.run()
		})
	}
	if cs, ok := t.conns[addr]; ok {
		if len(cs.Conns) < t.MaxConnsPerHost {
			if cq, ok := t.idleConns[addr]; ok && cq.Length() > 0 {
				pc = cq.Dequeue()
				pc.lastTime = time.Now()
			} else {
				var client *Client
				if t.Options != nil {
					client, err = t.DialWithOptions(addr, t.Options)
				} else {
					client, err = t.Dial(t.Network, addr, t.Codec)
				}
				if err != nil {
					return nil, err
				}
				pc = newPersistConn(client)
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
			var client *Client
			if t.Options != nil {
				client, err = t.DialWithOptions(addr, t.Options)
			} else {
				client, err = t.Dial(t.Network, addr, t.Codec)
			}
			if err != nil {
				return nil, err
			}
			pc = newPersistConn(client)
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
		var client *Client
		if t.Options != nil {
			client, err = t.DialWithOptions(addr, t.Options)
		} else {
			client, err = t.Dial(t.Network, addr, t.Codec)
		}
		if err != nil {
			return nil, err
		}
		pc = newPersistConn(client)
	}
	cs := &conns{addr: addr}
	cs.Append(pc)
	t.conns[addr] = cs
	return pc, nil
}
func (t *Transport) run() {
	for {
		select {
		case <-time.After(time.Second):
			t.now = time.Now()
			t.connsMu.Lock()
			for _, cs := range t.conns {
				length := len(cs.Conns)
				for i := 0; i < length; i++ {
					pc := cs.Conns[i]
					if pc.lastTime.Add(t.KeepAlive).Before(time.Now()) && pc.NumCalls() == 0 {
						pc.lastTime = time.Now()
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
					}
				}
				if cq.Length() == 0 {
					delete(t.idleConns, cq.addr)
				}
			}

			t.connsMu.Unlock()
		case <-t.done:
			return
		}
	}
}
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

type conns struct {
	addr   string
	Conns  []*persistConn
	cursor int
}

func (c *conns) Cursor() int {
	c.cursor += 1
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

type persistConn struct {
	*Client
	mu       sync.Mutex
	alive    bool
	lastTime time.Time
}

func newPersistConn(client *Client) *persistConn {
	return &persistConn{
		Client:   client,
		alive:    true,
		lastTime: time.Now(),
	}
}

type node struct {
	value    *persistConn
	previous *node
	next     *node
}

func (n *node) Value() *persistConn {
	return n.value
}

func (n *node) Set(value *persistConn) {
	n.value = value
}

func (n *node) Previous() *node {
	return n.previous
}

func (n *node) Next() *node {
	return n.next
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
