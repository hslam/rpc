// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package rpc provides access to the exported methods of an object across a network or other I/O connection.
package rpc

import (
	"crypto/tls"
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/funcs"
	"github.com/hslam/netpoll"
	"github.com/hslam/socket"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var numCPU = runtime.NumCPU()

// Server represents an RPC Server.
type Server struct {
	Registry      *funcs.Funcs
	ctxPool       *sync.Pool
	upgradePool   *sync.Pool
	upgradeBuffer []byte
	pipelining    bool
	numWorkers    int
	poll          bool
	mut           sync.Mutex
	listeners     []socket.Listener
	mutex         sync.RWMutex
	codecs        map[ServerCodec]io.Closer
	watchs        map[string]map[ServerCodec]*Context
	watchFunc     WatchFunc
}

// NewServerCodecFunc is the function making a new ServerCodec by socket.Messages.
type NewServerCodecFunc func(messages socket.Messages) ServerCodec

// WatchFunc is the function getting value by key.
type WatchFunc func(key string) (value []byte, ok bool)

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		Registry:      funcs.New(),
		ctxPool:       &sync.Pool{New: func() interface{} { return &Context{} }},
		upgradePool:   &sync.Pool{New: func() interface{} { return &upgrade{} }},
		upgradeBuffer: make([]byte, 1024),
		numWorkers:    numCPU,
		codecs:        make(map[ServerCodec]io.Closer),
		watchs:        make(map[string]map[ServerCodec]*Context),
	}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//	- exported method of exported type
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (server *Server) Register(obj interface{}) error {
	return server.Registry.Register(obj)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (server *Server) RegisterName(name string, obj interface{}) error {
	return server.Registry.RegisterName(name, obj)
}

// SetPipelining enables the Server to use pipelining.
func (server *Server) SetPipelining(enable bool) {
	server.pipelining = enable
}

// SetPoll enables the Server to use netpoll based on epoll/kqueue.
func (server *Server) SetPoll(enable bool) {
	server.poll = enable
}

func (server *Server) getUpgrade() *upgrade {
	return server.upgradePool.Get().(*upgrade)
}

func (server *Server) putUpgrade(u *upgrade) {
	u.Reset()
	server.upgradePool.Put(u)
}

// Push triggers the waiting clients with the watch key value..
func (server *Server) Push(key string, value []byte) {
	server.mutex.RLock()
	watchs := server.watchs[key]
	server.mutex.RUnlock()
	for _, ctx := range watchs {
		server.push(ctx, value)
	}
}

func (server *Server) push(ctx *Context, value []byte) {
	ctx.value = value
	server.sendResponse(ctx)
	ctx.value = nil
}

// PushFunc sets a WatchFunc.
func (server *Server) PushFunc(watchFunc WatchFunc) {
	server.watchFunc = watchFunc
}

// ServeCodec uses the specified codec to decode requests and encode responses.
func (server *Server) ServeCodec(codec ServerCodec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	var ch chan *Context
	var done chan struct{}
	var workers chan struct{}
	var worker bool
	if !server.pipelining {
		ch = make(chan *Context)
		done = make(chan struct{}, 1)
		workers = make(chan struct{}, server.numWorkers)
		worker = true
	}
	for {
		err := server.ServeRequest(codec, nil, sending, wg, worker, ch, done, workers)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}
	wg.Wait()
	server.mutex.Lock()
	delete(server.codecs, codec)
	for k, events := range server.watchs {
		if _, ok := events[codec]; ok {
			delete(events, codec)
		}
		if len(events) == 0 {
			delete(server.watchs, k)
		}
	}
	server.mutex.Unlock()
	codec.Close()
	if done != nil {
		close(done)
	}
}

// ServeRequest is like ServeCodec but synchronously serves a single request.
// It does not close the codec upon completion.
func (server *Server) ServeRequest(codec ServerCodec, recving *sync.Mutex, sending *sync.Mutex, wg *sync.WaitGroup, worker bool, ch chan *Context, done chan struct{}, workers chan struct{}) error {
	if recving != nil {
		recving.Lock()
	}
	ctx, err := server.readRequest(codec)
	if recving != nil {
		recving.Unlock()
	}
	ctx.sending = sending
	ctx.codec = codec
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF && err != netpoll.EAGAIN {
			logger.Errorln(err)
		}
		if !ctx.keepReading {
			ctx.Reset()
			server.ctxPool.Put(ctx)
			return err
		}
		ctx.Error = err.Error()
		server.sendResponse(ctx)
		return err
	}
	if ctx.upgrade.Heartbeat == heartbeat {
		server.sendResponse(ctx)
		return nil
	} else if ctx.upgrade.Watch == stopWatch {
		server.mutex.Lock()
		if events, ok := server.watchs[ctx.ServiceMethod]; ok {
			delete(events, codec)
		}
		server.mutex.Unlock()
		server.sendResponse(ctx)
		return nil
	} else if ctx.upgrade.Watch == watch {
		server.mutex.Lock()
		var events map[ServerCodec]*Context
		var ok bool
		if events, ok = server.watchs[ctx.ServiceMethod]; !ok {
			events = make(map[ServerCodec]*Context)
		}
		events[codec] = ctx
		server.watchs[ctx.ServiceMethod] = events
		server.mutex.Unlock()
		if server.watchFunc == nil {
			return nil
		}
	}
	if server.pipelining {
		server.callService(nil, ctx)
		return nil
	}
	if wg != nil {
		wg.Add(1)
	}
	if worker {
		select {
		case ch <- ctx:
		case workers <- struct{}{}:
			go func(done chan struct{}, ch chan *Context, server *Server, wg *sync.WaitGroup, workers chan struct{}, ctx *Context) {
				defer func() { <-workers }()
				for {
					server.callService(wg, ctx)
					t := time.NewTimer(time.Second)
					runtime.Gosched()
					select {
					case ctx = <-ch:
						t.Stop()
					case <-t.C:
						return
					case <-done:
						return
					}
				}
			}(done, ch, server, wg, workers, ctx)
		default:
			go server.callService(wg, ctx)
		}
	} else {
		go server.callService(wg, ctx)
	}
	return nil
}

func (server *Server) readRequest(codec ServerCodec) (ctx *Context, err error) {
	ctx = server.ctxPool.Get().(*Context)
	ctx.upgrade = server.getUpgrade()
	err = codec.ReadRequestHeader(ctx)
	if err != nil {
		return
	}
	ctx.keepReading = true
	if len(ctx.Upgrade) > 0 {
		ctx.upgrade.Unmarshal(ctx.Upgrade)
		ctx.Upgrade = nil
	}
	if ctx.upgrade.NoRequest != noRequest {
		ctx.f = server.Registry.GetFunc(ctx.ServiceMethod)
		if ctx.f == nil {
			err = errors.New("can't find service " + ctx.ServiceMethod)
			codec.ReadRequestBody(nil)
			return
		}
		ctx.args = ctx.f.GetValueIn(0)
		if ctx.args == funcs.ZeroValue {
			err = errors.New("can't find args")
			codec.ReadRequestBody(nil)
			return
		}
		if err = codec.ReadRequestBody(ctx.args.Interface()); err != nil {
			return
		}
	}
	if ctx.upgrade.NoResponse != noResponse {
		ctx.reply = ctx.f.GetValueIn(1)
		if ctx.reply == funcs.ZeroValue {
			err = errors.New("can't find reply")
		}
	}
	return
}

func (server *Server) callService(wg *sync.WaitGroup, ctx *Context) {
	if wg != nil {
		defer wg.Done()
	}
	if ctx.upgrade.Watch == watch {
		if server.watchFunc != nil {
			if value, ok := server.watchFunc(ctx.ServiceMethod); ok {
				server.push(ctx, value)
			}
		}
		return
	}
	if err := ctx.f.ValueCall(ctx.args, ctx.reply); err != nil {
		ctx.CallError = true
		ctx.Error = err.Error()
	}
	server.sendResponse(ctx)
}

func (server *Server) sendResponse(ctx *Context) {
	ctx.sending.Lock()
	var reply interface{}
	if len(ctx.Error) == 0 && ctx.upgrade.NoResponse != noResponse {
		reply = ctx.reply.Interface()
	}
	err := ctx.codec.WriteResponse(ctx, reply)
	if err != nil {
		logger.Errorln("writing response:", err)
	}
	ctx.sending.Unlock()
	if ctx.upgrade.Watch != watch {
		server.putUpgrade(ctx.upgrade)
		ctx.Reset()
		server.ctxPool.Put(ctx)
	}
}

func (server *Server) listen(sock socket.Socket, address string, New NewServerCodecFunc) error {
	logger.Noticef("pid - %d", os.Getpid())
	if server.poll {
		logger.Noticef("poll - %s", netpoll.Tag)
	} else {
		logger.Noticef("poll - %s", "disabled")
	}
	logger.Noticef("network - %s", sock.Scheme())
	logger.Noticef("listening on %s", address)
	lis, err := sock.Listen(address)
	if err != nil {
		return err
	}
	server.mut.Lock()
	server.listeners = append(server.listeners, lis)
	server.mut.Unlock()
	type ServerContext struct {
		codec   ServerCodec
		recving *sync.Mutex
		sending *sync.Mutex
		worker  bool
		wg      *sync.WaitGroup
		ch      chan *Context
		done    chan struct{}
		workers chan struct{}
		pipe    int32
		closed  int32
	}
	if server.poll {
		lis.ServeMessages(func(messages socket.Messages) (socket.Context, error) {
			codec := New(messages)
			server.mutex.Lock()
			server.codecs[codec] = messages
			server.mutex.Unlock()
			return &ServerContext{
				codec:   codec,
				recving: new(sync.Mutex),
				sending: new(sync.Mutex),
				wg:      new(sync.WaitGroup),
			}, nil
		}, func(context socket.Context) error {
			ctx := context.(*ServerContext)
			if server.pipelining {
				if !atomic.CompareAndSwapInt32(&ctx.pipe, 0, 1) {
					return nil
				}
			}
			err := server.ServeRequest(ctx.codec, ctx.recving, ctx.sending, ctx.wg, ctx.worker, ctx.ch, ctx.done, ctx.workers)
			if server.pipelining {
				atomic.StoreInt32(&ctx.pipe, 0)
			}
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if atomic.CompareAndSwapInt32(&ctx.closed, 0, 1) {
					ctx.wg.Wait()
					server.mutex.Lock()
					delete(server.codecs, ctx.codec)
					for _, events := range server.watchs {
						delete(events, ctx.codec)
					}
					server.mutex.Unlock()
					ctx.codec.Close()
				}
			}
			return err
		})
		return nil
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go func() {
			messages := conn.Messages()
			codec := New(messages)
			server.mutex.Lock()
			server.codecs[codec] = messages
			server.mutex.Unlock()
			server.ServeCodec(codec)
		}()
	}
}

// Listen announces on the local network address.
func (server *Server) Listen(network, address string, codec string) error {
	if newSocket := NewSocket(network); newSocket != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return server.listen(newSocket(nil), address, func(messages socket.Messages) ServerCodec {
				return NewServerCodec(newCodec(), nil, messages)
			})
		}
		return errors.New("unsupported codec: " + codec)
	}
	return errors.New("unsupported protocol scheme: " + network)
}

// ListenTLS announces on the local network address with tls.Config.
func (server *Server) ListenTLS(network, address string, codec string, config *tls.Config) error {
	if newSocket := NewSocket(network); newSocket != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return server.listen(newSocket(config), address, func(messages socket.Messages) ServerCodec {
				return NewServerCodec(newCodec(), nil, messages)
			})
		}
		return errors.New("unsupported codec: " + codec)
	}
	return errors.New("unsupported protocol scheme: " + network)
}

// ListenWithOptions announces on the local network address with Options.
func (server *Server) ListenWithOptions(address string, opts *Options) error {
	if opts.NewCodec == nil && opts.NewHeaderEncoder == nil && opts.Codec == "" {
		return errors.New("need opts.NewCodec, opts.NewEncoder or opts.Codec")
	}
	if opts.NewSocket == nil && opts.Network == "" {
		return errors.New("need opts.NewSocket, opts.NewMessages or opts.Network")
	}
	var sock socket.Socket
	if newSocket := NewSocket(opts.Network); newSocket != nil {
		sock = newSocket(opts.TLSConfig)
	} else if opts.NewSocket != nil {
		sock = opts.NewSocket(opts.TLSConfig)
	}
	return server.listen(sock, address, func(messages socket.Messages) ServerCodec {
		var bodyCodec codec.Codec
		if newCodec := NewCodec(opts.Codec); newCodec != nil {
			bodyCodec = newCodec()
		} else if opts.NewCodec != nil {
			bodyCodec = opts.NewCodec()
		}
		var headerEncoder *Encoder
		if newEncoder := NewHeaderEncoder(opts.HeaderEncoder); newEncoder != nil {
			headerEncoder = newEncoder()
		} else if opts.NewHeaderEncoder != nil {
			headerEncoder = opts.NewHeaderEncoder()
		}
		return NewServerCodec(bodyCodec, headerEncoder, messages)
	})
}

// Close closes the server.
func (server *Server) Close() error {
	server.mut.Lock()
	for _, lis := range server.listeners {
		lis.Close()
	}
	server.listeners = []socket.Listener{}
	server.mut.Unlock()

	server.mutex.Lock()
	for codec, closer := range server.codecs {
		delete(server.codecs, codec)
		closer.Close()
	}
	server.mutex.Unlock()
	return nil
}

// Register publishes the receiver's methods in the DefaultServer.
func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func RegisterName(name string, rcvr interface{}) error {
	return DefaultServer.RegisterName(name, rcvr)
}

// SetPipelining enables the Server to use pipelining.
func SetPipelining(enable bool) {
	DefaultServer.SetPipelining(enable)
}

// SetPoll enables the Server to use netpoll based on epoll/kqueue.
func SetPoll(enable bool) {
	DefaultServer.poll = enable
}

// Push triggers the waiting clients with the watch key value.
func Push(key string, value []byte) {
	DefaultServer.Push(key, value)
}

// PushFunc sets a WatchFunc.
func PushFunc(watchFunc WatchFunc) {
	DefaultServer.PushFunc(watchFunc)
}

// ServeCodec uses the specified codec to decode requests and encode responses.
func ServeCodec(codec ServerCodec) {
	DefaultServer.ServeCodec(codec)
}
