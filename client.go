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

// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type Client struct {
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

// NewClient returns a new Client to handle requests to the
// set of services at the other end of the connection.
// It adds a buffer to the write side of the connection so
// the header and payload are sent as a unit.
//
// The read and write halves of the connection are serialized independently,
// so no interlocking is required. However each half may be accessed
// concurrently so the implementation of conn should protect against
// concurrent reads or concurrent writes.
func NewClient() *Client {
	return &Client{
		pending:       make(map[uint64]*Call),
		watchs:        make(map[string]*Call),
		callPool:      &sync.Pool{New: func() interface{} { return &Call{} }},
		donePool:      &sync.Pool{New: func() interface{} { return make(chan *Call, 10) }},
		upgradePool:   &sync.Pool{New: func() interface{} { return &upgrade{} }},
		upgradeBuffer: make([]byte, 1024),
	}
}

// Dial connects to an RPC server at the specified network address.
func (client *Client) Dial(s socket.Socket, address string, New NewClientCodecFunc) (*Client, error) {
	conn, err := s.Dial(address)
	if err != nil {
		return nil, err
	}
	client.codec = New(conn.Messages())
	go client.read()
	return client, nil
}

// NewClientWithCodec is like NewClient but uses the specified
// codec to encode requests and decode responses.
func NewClientWithCodec(codec ClientCodec) *Client {
	if codec == nil {
		return nil
	}
	c := NewClient()
	c.codec = codec
	go c.read()
	return c
}

func (client *Client) getUpgrade() *upgrade {
	return client.upgradePool.Get().(*upgrade)
}

func (client *Client) putUpgrade(u *upgrade) {
	u.Reset()
	client.upgradePool.Put(u)
}

func (client *Client) write(call *Call) {
	client.reqMutex.Lock()
	defer client.reqMutex.Unlock()
	client.mutex.Lock()
	if client.shutdown || client.closing {
		client.mutex.Unlock()
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := client.seq
	if call.upgrade.Watch == watch {
		if _, ok := client.watchs[call.ServiceMethod]; ok {
			client.mutex.Unlock()
			call.Error = ErrWatch
			call.done()
			return
		}
		client.watchs[call.ServiceMethod] = call
	}
	client.seq++
	client.pending[seq] = call
	client.mutex.Unlock()
	client.ctx = Context{}
	client.ctx.Seq = seq
	client.ctx.upgrade = call.upgrade
	if call.upgrade.Heartbeat == heartbeat ||
		call.upgrade.Watch == watch ||
		call.upgrade.Watch == stopWatch ||
		call.upgrade.NoRequest == noRequest ||
		call.upgrade.NoResponse == noResponse {
		client.ctx.Upgrade, _ = call.upgrade.Marshal(client.upgradeBuffer)
	}
	client.ctx.ServiceMethod = call.ServiceMethod
	err := client.codec.WriteRequest(&client.ctx, call.Args)
	if err != nil {
		client.mutex.Lock()
		delete(client.pending, seq)
		if call.upgrade.Watch == watch {
			delete(client.watchs, call.ServiceMethod)
		}
		client.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (client *Client) read() {
	var err error
	var ctx Context
	for err == nil {
		ctx = Context{}
		err = client.codec.ReadResponseHeader(&ctx)
		if err != nil {
			break
		}
		seq := ctx.Seq
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		switch {
		case call == nil:
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		case ctx.Error != "":
			call.Error = errors.New(ctx.Error)
			err = client.codec.ReadResponseBody(nil)
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
					client.putUpgrade(u)
				} else if u.Watch == stopWatch {
					client.mutex.Lock()
					if _, ok := client.watchs[call.ServiceMethod]; ok {
						delete(client.pending, seq)
					}
					delete(client.watchs, call.ServiceMethod)
					client.mutex.Unlock()
					call.done()
					client.putUpgrade(u)
				} else if u.Watch == watch {
					call.Value = ctx.value
					client.mutex.Lock()
					if _, ok := client.watchs[call.ServiceMethod]; ok {
						client.pending[seq] = call
					}
					client.mutex.Unlock()
					if len(call.Value) == 0 {
						call.done()
					} else {
						call.watch()
					}
				}
				continue
			}
			client.putUpgrade(u)
			err = client.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	client.reqMutex.Lock()
	client.mutex.Lock()
	client.shutdown = true
	closing := client.closing
	if err == io.EOF {
		err = ErrShutdown
	}
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
	for _, call := range client.watchs {
		if call.watcher != nil {
			call.watcher.stop()
		}
	}
	client.pending = make(map[uint64]*Call)
	client.watchs = make(map[string]*Call)
	client.mutex.Unlock()
	client.reqMutex.Unlock()
	if err != io.EOF && !closing {
		logger.Allln("rpc: client protocol error:", err)
	}
}

// NumCalls returns the number of calls.
func (client *Client) NumCalls() (n uint64) {
	client.mutex.Lock()
	n = uint64(len(client.pending))
	w := uint64(len(client.watchs))
	if w > n {
		n = w
	}
	client.mutex.Unlock()
	return
}

// Close calls the underlying codec's Close method. If the connection is already
// shutting down, ErrShutdown is returned.
func (client *Client) Close() error {
	client.mutex.Lock()
	if client.closing {
		client.mutex.Unlock()
		return nil
	}
	client.closing = true
	client.mutex.Unlock()
	return client.codec.Close()
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (client *Client) RoundTrip(call *Call) *Call {
	done := call.Done
	done = checkDone(done)
	call.upgrade = client.getUpgrade()
	call.Done = done
	client.write(call)
	return call
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = client.getUpgrade()
	done = checkDone(done)
	call.Done = done
	client.write(call)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	done := client.donePool.Get().(chan *Call)
	upgrade := client.getUpgrade()
	call := client.callPool.Get().(*Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = upgrade
	call.Done = done
	client.write(call)
	runtime.Gosched()
	<-call.Done
	ResetDone(done)
	client.donePool.Put(done)
	err := call.Error
	*call = Call{}
	client.callPool.Put(call)
	return err
}

// CallTimeout acts like Call but takes a timeout.
func (client *Client) CallTimeout(serviceMethod string, args interface{}, reply interface{}, timeout time.Duration) error {
	if timeout <= 0 {
		return client.Call(serviceMethod, args, reply)
	}
	done := client.donePool.Get().(chan *Call)
	upgrade := client.getUpgrade()
	call := client.callPool.Get().(*Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.upgrade = upgrade
	call.Done = done
	client.write(call)
	timer := time.NewTimer(timeout)
	var err error
	runtime.Gosched()
	select {
	case <-call.Done:
		timer.Stop()
		ResetDone(done)
		client.donePool.Put(done)
		err = call.Error
		*call = Call{}
		client.callPool.Put(call)
	case <-timer.C:
		err = ErrTimeout
	}
	return err
}

// Watch returns the Watcher.
func (client *Client) Watch(key string) (Watcher, error) {
	watcher := &watcher{client: client, C: make(chan *event, 10), key: key, done: make(chan struct{}, 1)}
	upgrade := client.getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Watch = watch
	done := client.donePool.Get().(chan *Call)
	call := new(Call)
	call.ServiceMethod = key
	call.upgrade = upgrade
	call.Done = done
	call.watcher = watcher
	client.write(call)
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
func (client *Client) stopWatch(key string) error {
	upgrade := client.getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Watch = stopWatch
	done := client.donePool.Get().(chan *Call)
	call := client.callPool.Get().(*Call)
	call.ServiceMethod = key
	call.upgrade = upgrade
	call.Done = done
	client.write(call)
	runtime.Gosched()
	<-call.Done
	ResetDone(done)
	client.donePool.Put(done)
	err := call.Error
	*call = Call{}
	client.callPool.Put(call)
	return err
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (client *Client) Ping() error {
	upgrade := client.getUpgrade()
	upgrade.NoRequest = noRequest
	upgrade.NoResponse = noResponse
	upgrade.Heartbeat = heartbeat
	done := client.donePool.Get().(chan *Call)
	call := client.callPool.Get().(*Call)
	call.upgrade = upgrade
	call.Done = done
	client.write(call)
	runtime.Gosched()
	<-call.Done
	ResetDone(done)
	client.donePool.Put(done)
	err := call.Error
	*call = Call{}
	client.callPool.Put(call)
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
