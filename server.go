// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package rpc implements a remote procedure call over TCP, UNIX, HTTP and WS.
package rpc

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/hslam/buffer"
	"github.com/hslam/funcs"
	"github.com/hslam/log"
	"github.com/hslam/netpoll"
	"github.com/hslam/scheduler"
	"github.com/hslam/socket"
	"github.com/hslam/transition"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

var numCPU = runtime.NumCPU()

// Server represents an RPC Server.
type Server struct {
	Funcs       *funcs.Funcs
	logger      *log.Logger
	ctxPool     *sync.Pool
	upgradePool *sync.Pool
	bufferSize  int
	bufferPool  *buffer.Pool
	pipelining  bool
	noBatch     bool
	poll        bool
	shared      bool
	noCopy      bool
	mut         sync.Mutex
	listeners   []socket.Listener
	mutex       sync.RWMutex
	codecs      map[ServerCodec]io.Closer
	watchs      map[string]map[ServerCodec]*Context
	watchFunc   WatchFunc
}

// NewServerCodecFunc is the function making a new ServerCodec by socket.Messages.
type NewServerCodecFunc func(messages socket.Messages) ServerCodec

// WatchFunc is the function getting value by key.
type WatchFunc func(key string) (value []byte, ok bool)

// NewServer returns a new Server.
func NewServer() *Server {
	var logger = log.New()
	logger.SetPrefix(logPrefix)
	logger.SetLevel(log.Level(InfoLogLevel))
	return &Server{
		Funcs:       funcs.New(),
		logger:      logger,
		ctxPool:     &sync.Pool{New: func() interface{} { return &Context{} }},
		upgradePool: &sync.Pool{New: func() interface{} { return &upgrade{} }},
		bufferSize:  bufferSize,
		bufferPool:  buffer.AssignPool(bufferSize),
		codecs:      make(map[ServerCodec]io.Closer),
		watchs:      make(map[string]map[ServerCodec]*Context),
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
	return server.Funcs.Register(obj)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (server *Server) RegisterName(name string, obj interface{}) error {
	return server.Funcs.RegisterName(name, obj)
}

// Services returns registered services.
func (server *Server) Services() []string {
	return server.Funcs.Services()
}

//SetBufferSize sets buffer size.
func (server *Server) SetBufferSize(size int) {
	if size > 0 {
		server.bufferPool = buffer.AssignPool(size)
	} else {
		server.bufferPool = nil
	}
}

//SetContextBuffer sets shared buffer.
func (server *Server) SetContextBuffer(shared bool) {
	server.shared = shared
}

// SetNoCopy reuses a buffer from the pool for minimizing memory allocations.
// The RPC handler takes ownership of buffer, and the handler should not use buffer after this handle.
// The default noCopy is false to make a copy of data for every RPC handler.
func (server *Server) SetNoCopy(noCopy bool) {
	server.noCopy = noCopy
}

//SetLogLevel sets log's level.
func (server *Server) SetLogLevel(level LogLevel) {
	server.logger.SetLevel(log.Level(level))
}

//GetLogLevel returns log's level.
func (server *Server) GetLogLevel() LogLevel {
	return LogLevel(server.logger.GetLevel())
}

// SetPipelining enables the Server to use pipelining per connection.
func (server *Server) SetPipelining(enable bool) {
	server.pipelining = enable
}

// SetNoBatch disables the Server to use batch writer.
func (server *Server) SetNoBatch(noBatch bool) {
	server.noBatch = noBatch
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
	wg := new(sync.WaitGroup)
	var sched scheduler.Scheduler
	if server.pipelining {
		sched = scheduler.New(1, &scheduler.Options{Threshold: 2})
	}
	messages := codec.Messages()
	var trans = transition.NewTransition(16, codec.Concurrency)
	for {
		ctx := server.ctxPool.Get().(*Context)
		ctx.upgrade = server.getUpgrade()
		if server.bufferPool != nil {
			ctx.buffer = server.bufferPool.GetBuffer(server.bufferSize)
		}
		ctx.codec = codec
		data, err := messages.ReadMessage(ctx.buffer)
		if err != nil {
			break
		}
		ctx.data = data
		trans.Smooth(func() {
			server.ServeRequest(ctx, nil, wg, sched)
		}, func() {
			server.ServeRequest(ctx, nil, wg, sched)
		})
	}
	wg.Wait()
	server.mutex.Lock()
	server.deleteCodec(codec)
	server.mutex.Unlock()
	codec.Close()
	if sched != nil {
		sched.Close()
	}
	if trans != nil {
		trans.Close()
	}
}

// deleteCodec closes the specified codec.
func (server *Server) deleteCodec(codec ServerCodec) {
	delete(server.codecs, codec)
	for k, events := range server.watchs {
		delete(events, codec)
		if len(events) == 0 {
			delete(server.watchs, k)
		}
	}
}

// ServeRequest is like ServeCodec but synchronously serves a single request.
// It does not close the codec upon completion.
func (server *Server) ServeRequest(ctx *Context, recving *sync.Mutex, wg *sync.WaitGroup, sched scheduler.Scheduler) error {
	err := server.readRequestHeader(ctx)
	if err != nil {
		server.putUpgrade(ctx.upgrade)
		if server.bufferPool != nil && cap(ctx.buffer) > 0 {
			server.bufferPool.PutBuffer(ctx.buffer)
		}
		ctx.Reset()
		server.ctxPool.Put(ctx)
		return err
	}
	var codec = ctx.codec
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
	}
	wg.Add(1)
	if sched != nil {
		sched.Schedule(func() {
			server.handleRequest(wg, ctx)
		})
	} else {
		scheduler.Schedule(func() {
			server.handleRequest(wg, ctx)
		})
	}
	return nil
}

func (server *Server) readRequestHeader(ctx *Context) (err error) {
	err = ctx.codec.ReadRequestHeader(ctx)
	if err == nil && len(ctx.Upgrade) > 0 {
		ctx.upgrade.Unmarshal(ctx.Upgrade)
		ctx.Upgrade = nil
	}
	return err
}

func (server *Server) handleRequest(wg *sync.WaitGroup, ctx *Context) {
	if wg != nil {
		defer wg.Done()
	}
	err := server.readRequestBody(ctx)
	if err != nil {
		server.logger.Errorln(err)
		ctx.Error = err.Error()
		server.sendResponse(ctx)
	} else {
		server.callService(ctx)
	}
}

func (server *Server) readRequestBody(ctx *Context) (err error) {
	var codec = ctx.codec
	if ctx.upgrade.NoRequest != noRequest {
		ctx.f = server.Funcs.GetFunc(ctx.ServiceMethod)
		if ctx.f == nil {
			err = errors.New("can't find service " + ctx.ServiceMethod)
			codec.ReadRequestBody(nil, nil)
			return
		}
		ctx.args = ctx.f.GetValueIn(0)
		if ctx.args == funcs.ZeroValue {
			err = errors.New("can't find args")
			codec.ReadRequestBody(nil, nil)
			return
		}
		var value []byte
		if server.noCopy {
			value = ctx.value
		} else {
			if ctx.f.WithContext() && server.shared {
				value = GetBuffer(len(ctx.value))[:len(ctx.value)]
			} else {
				value = make([]byte, len(ctx.value))
			}
			copy(value, ctx.value)
		}
		if err = codec.ReadRequestBody(value, ctx.args.Interface()); err != nil {
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

func (server *Server) callService(ctx *Context) {
	if ctx.upgrade.Watch == watch {
		server.push(ctx, nil)
		if server.watchFunc != nil {
			if value, ok := server.watchFunc(ctx.ServiceMethod); ok {
				server.push(ctx, value)
			}
		}
		return
	}
	var err error
	if ctx.f.WithContext() {
		var c context.Context
		if !server.noCopy && server.shared {
			c = context.WithValue(context.Background(), BufferContextKey, ctx.value)
		} else {
			c = context.Background()
		}
		err = ctx.f.ValueCall(funcs.ValueOf(c), ctx.args, ctx.reply)
	} else {
		err = ctx.f.ValueCall(ctx.args, ctx.reply)
	}
	if err != nil {
		ctx.Error = err.Error()
	}
	server.sendResponse(ctx)
}

func (server *Server) sendResponse(ctx *Context) {
	var reply interface{}
	if len(ctx.Error) == 0 && ctx.upgrade.NoResponse != noResponse {
		reply = ctx.reply.Interface()
	}
	err := ctx.codec.WriteResponse(ctx, reply)
	if err != nil {
		server.logger.Errorln("writing response:", err)
	}
	if ctx.upgrade.Watch != watch {
		server.putUpgrade(ctx.upgrade)
		if server.bufferPool != nil && cap(ctx.buffer) > 0 {
			server.bufferPool.PutBuffer(ctx.buffer)
		}
		ctx.Reset()
		server.ctxPool.Put(ctx)
	}
}

func (server *Server) listen(sock socket.Socket, address string, New NewServerCodecFunc) error {
	server.logger.Noticef("pid - %d", os.Getpid())
	if server.poll {
		server.logger.Noticef("poll - %s", netpoll.Tag)
	} else {
		server.logger.Noticef("poll - %s", "disabled")
	}
	server.logger.Noticef("network - %s", sock.Scheme())
	if server.pipelining {
		server.logger.Noticef("io - pipelining")
	} else {
		server.logger.Noticef("io - multiplexing")
	}
	lis, err := sock.Listen(address)
	if err != nil {
		server.logger.Errorf("%s", err.Error())
		return err
	}
	server.logger.Noticef("listening on %s", address)
	codecs := make(map[ServerCodec]io.Closer)
	defer func() {
		server.mutex.Lock()
		for codec, closer := range codecs {
			delete(codecs, codec)
			server.deleteCodec(codec)
			closer.Close()
		}
		server.mutex.Unlock()
	}()
	server.mut.Lock()
	server.listeners = append(server.listeners, lis)
	server.mut.Unlock()
	type ServerContext struct {
		codec    ServerCodec
		recving  *sync.Mutex
		wg       *sync.WaitGroup
		messages socket.Messages
		trans    *transition.Transition
		sched    scheduler.Scheduler
		pipe     int32
		closed   int32
	}
	if server.poll {
		return lis.ServeMessages(func(messages socket.Messages) (socket.Context, error) {
			codec := New(messages)
			server.mutex.Lock()
			codecs[codec] = messages
			server.codecs[codec] = messages
			server.mutex.Unlock()
			var sched scheduler.Scheduler
			if server.pipelining {
				sched = scheduler.New(1, &scheduler.Options{Threshold: 2})
			}
			return &ServerContext{
				codec:    codec,
				recving:  new(sync.Mutex),
				wg:       new(sync.WaitGroup),
				messages: messages,
				trans:    transition.NewTransition(16, codec.Concurrency),
				sched:    sched,
			}, nil
		}, func(context socket.Context) error {
			svrctx := context.(*ServerContext)
			ctx := server.ctxPool.Get().(*Context)
			ctx.upgrade = server.getUpgrade()
			if server.bufferPool != nil {
				ctx.buffer = server.bufferPool.GetBuffer(server.bufferSize)
			}
			ctx.codec = svrctx.codec
			svrctx.recving.Lock()
			data, err := svrctx.messages.ReadMessage(ctx.buffer)
			if len(data) > 0 {
				ctx.data = data
				svrctx.trans.Smooth(func() {
					server.ServeRequest(ctx, svrctx.recving, svrctx.wg, svrctx.sched)
				}, func() {
					server.ServeRequest(ctx, svrctx.recving, svrctx.wg, svrctx.sched)
				})
			}
			svrctx.recving.Unlock()
			if len(data) == 0 {
				server.putUpgrade(ctx.upgrade)
				if server.bufferPool != nil && cap(ctx.buffer) > 0 {
					server.bufferPool.PutBuffer(ctx.buffer)
				}
				ctx.Reset()
				server.ctxPool.Put(ctx)
			}
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if atomic.CompareAndSwapInt32(&svrctx.closed, 0, 1) {
					svrctx.wg.Wait()
					server.mutex.Lock()
					delete(codecs, svrctx.codec)
					server.deleteCodec(svrctx.codec)
					server.mutex.Unlock()
					svrctx.codec.Close()
					if svrctx.sched != nil {
						svrctx.sched.Close()
					}
					if svrctx.trans != nil {
						svrctx.trans.Close()
					}
				}
			}
			return err
		})
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
			codecs[codec] = messages
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
				return NewServerCodec(newCodec(), nil, messages, server.noBatch)
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
				return NewServerCodec(newCodec(), nil, messages, server.noBatch)
			})
		}
		return errors.New("unsupported codec: " + codec)
	}
	return errors.New("unsupported protocol scheme: " + network)
}

// ListenWithOptions announces on the local network address with Options.
func (server *Server) ListenWithOptions(address string, opts *Options) error {
	if opts.NewCodec == nil && opts.NewHeaderEncoder == nil && opts.Codec == "" {
		return errors.New("need opts.NewCodec, opts.NewHeaderEncoder or opts.Codec")
	}
	if opts.NewSocket == nil && opts.Network == "" {
		return errors.New("need opts.NewSocket or opts.Network")
	}
	var sock socket.Socket
	if newSocket := NewSocket(opts.Network); newSocket != nil {
		sock = newSocket(opts.TLSConfig)
	} else if opts.NewSocket != nil {
		sock = opts.NewSocket(opts.TLSConfig)
	}
	return server.listen(sock, address, func(messages socket.Messages) ServerCodec {
		var bodyCodec Codec
		if newCodec := NewCodec(opts.Codec); newCodec != nil {
			bodyCodec = newCodec()
		} else if opts.NewCodec != nil {
			bodyCodec = opts.NewCodec()
		}
		var headerEncoder Encoder
		if newEncoder := NewHeaderEncoder(opts.HeaderEncoder); newEncoder != nil {
			headerEncoder = newEncoder()
		} else if opts.NewHeaderEncoder != nil {
			headerEncoder = opts.NewHeaderEncoder()
		}
		return NewServerCodec(bodyCodec, headerEncoder, messages, server.noBatch)
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
	return nil
}

// Register publishes the receiver's methods in the DefaultServer.
func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func RegisterName(name string, rcvr interface{}) error {
	return DefaultServer.RegisterName(name, rcvr)
}

// Services returns registered services.
func Services() []string {
	return DefaultServer.Services()
}

// SetPipelining enables the Server to use pipelining per connection.
func SetPipelining(enable bool) {
	DefaultServer.SetPipelining(enable)
}

// SetNoBatch disables the Server to use batch writer.
func SetNoBatch(noBatch bool) {
	DefaultServer.SetNoBatch(noBatch)
}

// SetPoll enables the Server to use netpoll based on epoll/kqueue.
func SetPoll(enable bool) {
	DefaultServer.SetPoll(enable)
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

//SetLogLevel sets log's level
func SetLogLevel(level LogLevel) {
	DefaultServer.SetLogLevel(level)
}

//GetLogLevel returns log's level
func GetLogLevel() LogLevel {
	return DefaultServer.GetLogLevel()
}

//SetBufferSize sets buffer size.
func SetBufferSize(size int) {
	DefaultServer.SetBufferSize(size)
}

//SetContextBuffer sets shared buffer.
func SetContextBuffer(shared bool) {
	DefaultServer.SetContextBuffer(shared)
}

// SetNoCopy reuses a buffer from the pool for minimizing memory allocations.
// The RPC handler takes ownership of buffer, and the handler should not use buffer after this handle.
// The default option is to make a copy of data for every RPC handler.
func SetNoCopy(noCopy bool) {
	DefaultServer.SetNoCopy(noCopy)
}
