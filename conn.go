// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"github.com/hslam/socket"
	"io"
	"runtime"
	"sync"
	"time"
)

// ErrTimeout is returned after the timeout,.
var ErrTimeout = errors.New("timeout")

// ErrShutdown is returned when the connection is shut down.
var ErrShutdown = errors.New("The connection is shut down")

// ErrWatch is returned when the watch is existed.
var ErrWatch = errors.New("The watch is existed")

// Call represents an active RPC.
type Call struct {
	upgrade       *upgrade
	Value         []byte
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	CallError     bool
	Error         error
	Done          chan *Call
	watcher       *watcher
}

func (call *Call) done() {
	select {
	case call.Done <- call:
	default:
	}
	call.watch()
}

func (call *Call) watch() {
	if call.watcher != nil {
		if len(call.Value) > 0 || call.Error != nil {
			e := getEvent()
			e.Value = call.Value
			e.Error = call.Error
			call.watcher.trigger(e)
		}
	}
}

// Conn represents an RPC Conn.
// There may be multiple outstanding Calls associated
// with a single Conn, and a Conn may be used by
// multiple goroutines simultaneously.
type Conn struct {
	codec         ClientCodec
	reqMutex      sync.Mutex
	ctx           Context
	mutex         sync.Mutex
	seq           uint64
	pending       map[uint64]*Call
	watchs        map[string]*Call
	callPool      *sync.Pool
	donePool      *sync.Pool
	upgradePool   *sync.Pool
	upgradeBuffer []byte
	closing       bool
	shutdown      bool
}

// NewClientCodecFunc is the function to make a new ClientCodec by socket.Messages.
type NewClientCodecFunc func(messages socket.Messages) ClientCodec

// NewConn returns a new Conn to handle requests to the
// set of services at the other end of the connection.
// It adds a buffer to the write side of the connection so
// the header and payload are sent as a unit.
//
// The read and write halves of the connection are serialized independently,
// so no interlocking is required. However each half may be accessed
// concurrently so the implementation of conn should protect against
// concurrent reads or concurrent writes.
func NewConn() *Conn {
	return &Conn{
		pending:       make(map[uint64]*Call),
		watchs:        make(map[string]*Call),
		callPool:      &sync.Pool{New: func() interface{} { return &Call{} }},
		donePool:      &sync.Pool{New: func() interface{} { return make(chan *Call, 10) }},
		upgradePool:   &sync.Pool{New: func() interface{} { return &upgrade{} }},
		upgradeBuffer: make([]byte, 1024),
	}
}

// Dial connects to an RPC server at the specified network address.
func (conn *Conn) Dial(s socket.Socket, address string, New NewClientCodecFunc) (*Conn, error) {
	c, err := s.Dial(address)
	if err != nil {
		return nil, err
	}
	conn.codec = New(c.Messages())
	go conn.read()
	return conn, nil
}

// NewConnWithCodec is like NewConn but uses the specified
// codec to encode requests and decode responses.
func NewConnWithCodec(codec ClientCodec) *Conn {
	if codec == nil {
		return nil
	}
	c := NewConn()
	c.codec = codec
	go c.read()
	return c
}

func (conn *Conn) getUpgrade() *upgrade {
	return conn.upgradePool.Get().(*upgrade)
}

func (conn *Conn) putUpgrade(u *upgrade) {
	u.Reset()
	conn.upgradePool.Put(u)
}

func (conn *Conn) write(call *Call) {
	conn.reqMutex.Lock()
	defer conn.reqMutex.Unlock()
	conn.mutex.Lock()
	if conn.shutdown || conn.closing {
		conn.mutex.Unlock()
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := conn.seq
	if call.upgrade.Watch == watch {
		if _, ok := conn.watchs[call.ServiceMethod]; ok {
			conn.mutex.Unlock()
			call.Error = ErrWatch
			call.done()
			return
		}
		conn.watchs[call.ServiceMethod] = call
	}
	conn.seq++
	conn.pending[seq] = call
	conn.mutex.Unlock()
	conn.ctx = Context{}
	conn.ctx.Seq = seq
	conn.ctx.upgrade = call.upgrade
	if call.upgrade.Heartbeat == heartbeat ||
		call.upgrade.Watch == watch ||
		call.upgrade.Watch == stopWatch ||
		call.upgrade.NoRequest == noRequest ||
		call.upgrade.NoResponse == noResponse {
		conn.ctx.Upgrade, _ = call.upgrade.Marshal(conn.upgradeBuffer)
	}
	conn.ctx.ServiceMethod = call.ServiceMethod
	err := conn.codec.WriteRequest(&conn.ctx, call.Args)
	if err != nil {
		conn.mutex.Lock()
		delete(conn.pending, seq)
		if call.upgrade.Watch == watch {
			delete(conn.watchs, call.ServiceMethod)
		}
		conn.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (conn *Conn) read() {
	var err error
	var ctx Context
	for err == nil {
		ctx = Context{}
		err = conn.codec.ReadResponseHeader(&ctx)
		if err != nil {
			break
		}
		seq := ctx.Seq
		conn.mutex.Lock()
		call := conn.pending[seq]
		delete(conn.pending, seq)
		conn.mutex.Unlock()
		switch {
		case call == nil:
			err = conn.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		case ctx.Error != "":
			call.Error = errors.New(ctx.Error)
			err = conn.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
			call.done()
		default:
			if len(ctx.value) > 0 {
				call.Value = ctx.value
			}
			u := call.upgrade
			if u.NoResponse == noResponse {
				if u.Heartbeat == heartbeat {
					call.done()
					conn.putUpgrade(u)
				} else if u.Watch == stopWatch {
					conn.mutex.Lock()
					if _, ok := conn.watchs[call.ServiceMethod]; ok {
						delete(conn.pending, seq)
					}
					delete(conn.watchs, call.ServiceMethod)
					conn.mutex.Unlock()
					call.done()
					conn.putUpgrade(u)
				} else if u.Watch == watch {
					call.Value = ctx.value
					conn.mutex.Lock()
					if _, ok := conn.watchs[call.ServiceMethod]; ok {
						conn.pending[seq] = call
					}
					conn.mutex.Unlock()
					if len(call.Value) == 0 {
						call.done()
					} else {
						call.watch()
					}
				}
				continue
			}
			conn.putUpgrade(u)
			err = conn.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	conn.reqMutex.Lock()
	conn.mutex.Lock()
	conn.shutdown = true
	closing := conn.closing
	if err == io.EOF {
		err = ErrShutdown
	}
	for _, call := range conn.pending {
		call.Error = err
		call.done()
	}
	for _, call := range conn.watchs {
		if call.watcher != nil {
			call.watcher.stop()
		}
	}
	conn.pending = make(map[uint64]*Call)
	conn.watchs = make(map[string]*Call)
	conn.mutex.Unlock()
	conn.reqMutex.Unlock()
	if err != io.EOF && !closing {
		logger.Allln("rpc: Conn protocol error:", err)
	}
}

// NumCalls returns the number of calls.
func (conn *Conn) NumCalls() (n uint64) {
	conn.mutex.Lock()
	p := uint64(len(conn.pending))
	n = uint64(len(conn.watchs))
	if p > n {
		n = p
	}
	conn.mutex.Unlock()
	return
}

// Close calls the underlying codec's Close method. If the connection is already
// shutting down, ErrShutdown is returned.
func (conn *Conn) Close() error {
	conn.mutex.Lock()
	if conn.closing {
		conn.mutex.Unlock()
		return nil
	}
	conn.closing = true
	conn.mutex.Unlock()
	return conn.codec.Close()
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (conn *Conn) RoundTrip(call *Call) *Call {
	done := call.Done
	done = checkDone(done)
	call.upgrade = conn.getUpgrade()
	call.Done = done
	conn.write(call)
	return call
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (conn *Conn) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = conn.getUpgrade()
	done = checkDone(done)
	call.Done = done
	conn.write(call)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (conn *Conn) Call(serviceMethod string, args interface{}, reply interface{}) error {
	done := conn.donePool.Get().(chan *Call)
	upgrade := conn.getUpgrade()
	call := conn.callPool.Get().(*Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = upgrade
	call.Done = done
	conn.write(call)
	runtime.Gosched()
	<-call.Done
	ResetDone(done)
	conn.donePool.Put(done)
	err := call.Error
	*call = Call{}
	conn.callPool.Put(call)
	return err
}

// CallTimeout acts like Call but takes a timeout.
func (conn *Conn) CallTimeout(serviceMethod string, args interface{}, reply interface{}, timeout time.Duration) error {
	if timeout <= 0 {
		return conn.Call(serviceMethod, args, reply)
	}
	done := conn.donePool.Get().(chan *Call)
	upgrade := conn.getUpgrade()
	call := conn.callPool.Get().(*Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = upgrade
	call.Done = done
	conn.write(call)
	timer := time.NewTimer(timeout)
	var err error
	runtime.Gosched()
	select {
	case <-call.Done:
		timer.Stop()
		ResetDone(done)
		conn.donePool.Put(done)
		err = call.Error
		*call = Call{}
		conn.callPool.Put(call)
	case <-timer.C:
		err = ErrTimeout
	}
	return err
}

// Watch returns the Watcher.
func (conn *Conn) Watch(key string) (Watcher, error) {
	watcher := &watcher{conn: conn, C: make(chan *event, 10), key: key, done: make(chan struct{}, 1)}
	upgrade := conn.getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Watch = watch
	done := conn.donePool.Get().(chan *Call)
	call := new(Call)
	call.ServiceMethod = key
	call.upgrade = upgrade
	call.Done = done
	call.watcher = watcher
	conn.write(call)
	var err error
	runtime.Gosched()
	<-call.Done
	err = call.Error
	if err != nil {
		watcher.stop()
		return nil, err
	}
	return watcher, err
}

// stopWatch stops the key watcher .
func (conn *Conn) stopWatch(key string) error {
	upgrade := conn.getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Watch = stopWatch
	done := conn.donePool.Get().(chan *Call)
	call := conn.callPool.Get().(*Call)
	call.ServiceMethod = key
	call.upgrade = upgrade
	call.Done = done
	conn.write(call)
	runtime.Gosched()
	<-call.Done
	ResetDone(done)
	conn.donePool.Put(done)
	err := call.Error
	*call = Call{}
	conn.callPool.Put(call)
	return err
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (conn *Conn) Ping() error {
	upgrade := conn.getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Heartbeat = heartbeat
	done := conn.donePool.Get().(chan *Call)
	call := conn.callPool.Get().(*Call)
	call.upgrade = upgrade
	call.Done = done
	conn.write(call)
	runtime.Gosched()
	<-call.Done
	ResetDone(done)
	conn.donePool.Put(done)
	err := call.Error
	*call = Call{}
	conn.callPool.Put(call)
	return err
}

// ResetDone resets the done.
func ResetDone(done chan *Call) {
	for len(done) > 0 {
		onceDone(done)
	}
}

func onceDone(done chan *Call) {
	select {
	case <-done:
	default:
	}
}

func checkDone(done chan *Call) chan *Call {
	if done == nil {
		return make(chan *Call, 10)
	}
	if cap(done) == 0 {
		panic("rpc: done channel is unbuffered")
	}
	return done
}
