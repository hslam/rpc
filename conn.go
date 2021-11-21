// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"errors"
	"github.com/hslam/buffer"
	"github.com/hslam/scheduler"
	"github.com/hslam/socket"
	"github.com/hslam/transition"
	"io"
	"runtime"
	"sync"
)

// ErrTimeout is returned after the timeout,.
var ErrTimeout = errors.New("timeout")

const shutdownMsg = "The connection is shut down"

// ErrShutdown is returned when the connection is shut down.
var ErrShutdown = errors.New(shutdownMsg)

// ErrWatch is returned when the watch is existed.
var ErrWatch = errors.New("The watch is existed")

var (
	callPool    = &sync.Pool{New: func() interface{} { return &Call{} }}
	donePool    = &sync.Pool{New: func() interface{} { return make(chan *Call, 10) }}
	ctxPool     = &sync.Pool{New: func() interface{} { return &Context{} }}
	upgradePool = &sync.Pool{New: func() interface{} { return &upgrade{} }}
)

func getUpgrade() *upgrade {
	return upgradePool.Get().(*upgrade)
}

func putUpgrade(u *upgrade) {
	u.Reset()
	upgradePool.Put(u)
}

func getContext() *Context {
	return ctxPool.Get().(*Context)
}

func putContext(ctx *Context) {
	ctx.Reset()
	ctxPool.Put(ctx)
}

// GetCall gets a call from the callPool.
func GetCall() *Call {
	call := callPool.Get().(*Call)
	call.Done = donePool.Get().(chan *Call)
	return call
}

// PutCall puts a call to the callPool.
func PutCall(call *Call) {
	if call.Done != nil {
		done := call.Done
		ResetDone(done)
		donePool.Put(done)
	}
	*call = Call{}
	callPool.Put(call)
}

// Call represents an active RPC.
type Call struct {
	Buffer        []byte
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
	codec          ClientCodec
	pipeMut        sync.Mutex
	mutex          sync.Mutex
	seq            uint64
	pending        map[uint64]*Call
	watchs         map[string]*Call
	bufferSize     int
	bufferPool     *buffer.Pool
	upgradeBuffer  []byte
	writeSched     scheduler.Scheduler
	readSched      scheduler.Scheduler
	writePipelines map[string]scheduler.Scheduler
	readPipelines  map[string]scheduler.Scheduler
	closing        bool
	shutdown       bool
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
		upgradeBuffer: make([]byte, 1),
		bufferSize:    bufferSize,
		bufferPool:    buffer.AssignPool(bufferSize),
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

func (conn *Conn) write(call *Call) {
	if conn.writePipelines != nil || conn.readSched != nil {
		conn.pipeMut.Lock()
		var sched scheduler.Scheduler
		if conn.writePipelines != nil {
			sched = conn.writePipelines[call.ServiceMethod]
			if sched == nil {
				sched = scheduler.New(1, &scheduler.Options{Threshold: 2})
				conn.writePipelines[call.ServiceMethod] = sched
			}
		} else {
			sched = conn.readSched
		}
		sched.Schedule(func() {
			conn.send(call)
		})
		conn.pipeMut.Unlock()
	} else {
		conn.send(call)
	}
}

func (conn *Conn) send(call *Call) {
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
	ctx := Context{}
	ctx.Seq = seq
	ctx.upgrade = call.upgrade
	if call.upgrade.Heartbeat == heartbeat ||
		call.upgrade.Watch == watch ||
		call.upgrade.Watch == stopWatch ||
		call.upgrade.NoRequest == noRequest ||
		call.upgrade.NoResponse == noResponse {
		ctx.Upgrade, _ = call.upgrade.Marshal(conn.upgradeBuffer)
	}
	ctx.ServiceMethod = call.ServiceMethod
	err := conn.codec.WriteRequest(&ctx, call.Args)
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
	var messages = conn.codec.Messages()
	wg := sync.WaitGroup{}
	var pipeline = scheduler.New(1, &scheduler.Options{Threshold: 2})
	var trans = transition.NewTransition(16, conn.codec.Concurrency)
	for err == nil {
		ctx := getContext()
		ctx.buffer = conn.bufferPool.GetBuffer(conn.bufferSize)
		ctx.data, err = messages.ReadMessage(ctx.buffer)
		if err != nil {
			break
		}
		trans.Smooth(func() {
			wg.Wait()
			conn.handle(ctx, false)
		}, func() {
			wg.Add(1)
			pipeline.Schedule(func() {
				conn.handle(ctx, true)
				wg.Done()
			})
		})
	}
	conn.mutex.Lock()
	conn.shutdown = true
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
	conn.mutex.Unlock()
	if conn.readPipelines != nil {
		for _, sched := range conn.readPipelines {
			sched.Close()
		}
	}
	if conn.readSched != nil {
		conn.readSched.Close()
	}
}

func (conn *Conn) handle(ctx *Context, async bool) {
	var err error
	err = conn.codec.ReadResponseHeader(ctx)
	if err != nil {
		return
	}
	seq := ctx.Seq
	conn.mutex.Lock()
	call := conn.pending[seq]
	if call != nil && call.upgrade.Watch != watch {
		delete(conn.pending, seq)
	}
	conn.mutex.Unlock()
	switch {
	case call == nil:
		err = conn.codec.ReadResponseBody(nil, nil)
		if err != nil {
			err = errors.New("reading error body: " + err.Error())
		}
		conn.bufferPool.PutBuffer(ctx.buffer)
		putContext(ctx)
	case len(ctx.Error) > 0:
		if ctx.Error == shutdownMsg {
			call.Error = ErrShutdown
		} else {
			call.Error = errors.New(ctx.Error)
		}
		err = conn.codec.ReadResponseBody(nil, nil)
		if err != nil {
			err = errors.New("reading error body: " + err.Error())
		}
		call.done()
		conn.bufferPool.PutBuffer(ctx.buffer)
		putContext(ctx)
	default:
		u := call.upgrade
		if u.NoResponse == noResponse {
			if u.Heartbeat == heartbeat {
				call.done()
				putUpgrade(u)
			} else if u.Watch == stopWatch {
				conn.mutex.Lock()
				if _, ok := conn.watchs[call.ServiceMethod]; ok {
					delete(conn.pending, seq)
				}
				delete(conn.watchs, call.ServiceMethod)
				conn.mutex.Unlock()
				call.done()
				putUpgrade(u)
			} else if u.Watch == watch {
				if len(ctx.value) > 0 {
					call.Value = make([]byte, len(ctx.value))
					copy(call.Value, ctx.value)
				}
				if len(call.Value) == 0 {
					call.done()
				} else {
					call.watch()
				}
			}
			conn.bufferPool.PutBuffer(ctx.buffer)
			putContext(ctx)
			return
		}
		putUpgrade(u)
		var sched scheduler.Scheduler
		if conn.readPipelines != nil {
			sched = conn.readPipelines[ctx.ServiceMethod]
			if sched == nil {
				sched = scheduler.New(1, &scheduler.Options{Threshold: 2})
				conn.readPipelines[ctx.ServiceMethod] = sched
			}
		} else if conn.readSched != nil {
			sched = conn.readSched
		}
		if sched != nil {
			sched.Schedule(func() {
				conn.finishCall(ctx, call, seq)
			})
		} else if async {
			scheduler.Schedule(func() {
				conn.finishCall(ctx, call, seq)
			})
		} else {
			conn.finishCall(ctx, call, seq)
		}
	}

}

func (conn *Conn) finishCall(ctx *Context, call *Call, seq uint64) {
	if len(ctx.value) > 0 {
		if cap(call.Buffer) >= len(ctx.value) {
			call.Value = call.Buffer[:len(ctx.value)]
		} else {
			call.Value = make([]byte, len(ctx.value))
		}
		copy(call.Value, ctx.value)
	}
	err := conn.codec.ReadResponseBody(call.Value, call.Reply)
	if err != nil {
		call.Error = errors.New("reading body " + err.Error())
	}
	call.done()
	buf := ctx.buffer
	conn.bufferPool.PutBuffer(buf)
	putContext(ctx)
}

// SetPipelining enables the client to use pipelining per connection.
func (conn *Conn) SetPipelining(enable bool) {
	conn.writeSched = scheduler.New(1, &scheduler.Options{Threshold: 2})
	conn.readSched = scheduler.New(1, &scheduler.Options{Threshold: 2})
}

// SetMethodPipelining enables the client to use pipelining per connection and method.
func (conn *Conn) SetMethodPipelining(enable bool) {
	conn.writePipelines = make(map[string]scheduler.Scheduler)
	conn.readPipelines = make(map[string]scheduler.Scheduler)
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
func (conn *Conn) Close() (err error) {
	conn.mutex.Lock()
	if conn.closing {
		conn.mutex.Unlock()
		return ErrShutdown
	}
	conn.closing = true
	conn.mutex.Unlock()
	err = conn.codec.Close()
	if conn.writePipelines != nil || conn.writeSched != nil {
		conn.pipeMut.Lock()
		if conn.writePipelines != nil {
			for _, sched := range conn.writePipelines {
				sched.Close()
			}
		}
		if conn.writeSched != nil {
			conn.writeSched.Close()
		}
		conn.pipeMut.Unlock()
	}

	return
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (conn *Conn) RoundTrip(call *Call) *Call {
	done := call.Done
	done = checkDone(done)
	call.upgrade = getUpgrade()
	call.Done = done
	conn.write(call)
	return call
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (conn *Conn) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := callPool.Get().(*Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = getUpgrade()
	done = checkDone(done)
	call.Done = done
	conn.write(call)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (conn *Conn) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := GetCall()
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = getUpgrade()
	conn.write(call)
	<-call.Done
	err := call.Error
	PutCall(call)
	return err
}

// CallWithContext acts like Call but takes a context.
func (conn *Conn) CallWithContext(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	call := GetCall()
	call.Buffer = GetContextBuffer(ctx)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = getUpgrade()
	conn.write(call)
	var err error
	select {
	case <-call.Done:
		err = call.Error
		PutCall(call)
	case <-ctx.Done():
		err = ctx.Err()
	}
	return err
}

// Watch returns the Watcher.
func (conn *Conn) Watch(key string) (Watcher, error) {
	watcher := &watcher{conn: conn, C: make(chan *event, 10), key: key, done: make(chan struct{}, 1)}
	upgrade := getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Watch = watch
	call := new(Call)
	call.ServiceMethod = key
	call.upgrade = upgrade
	call.Done = make(chan *Call, 1)
	call.watcher = watcher
	conn.write(call)
	var err error
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
	upgrade := getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Watch = stopWatch
	call := GetCall()
	call.ServiceMethod = key
	call.upgrade = upgrade
	conn.write(call)
	<-call.Done
	err := call.Error
	PutCall(call)
	return err
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (conn *Conn) Ping() error {
	upgrade := getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Heartbeat = heartbeat
	call := GetCall()
	call.upgrade = upgrade
	conn.write(call)
	runtime.Gosched()
	<-call.Done
	err := call.Error
	PutCall(call)
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
