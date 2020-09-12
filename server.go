// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package rpc provides access to the exported methods of an object across a network or other I/O connection.
package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/funcs"
	"github.com/hslam/poll"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/socket"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var numCPU = runtime.NumCPU()

type Server struct {
	Registry      *funcs.Funcs
	ctxPool       *sync.Pool
	upgradePool   *sync.Pool
	upgradeBuffer []byte
	pipelining    bool
	numWorkers    int
	poll          bool
}
type NewServerCodecFunc func(messages socket.Messages) ServerCodec

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		Registry:      funcs.New(),
		ctxPool:       &sync.Pool{New: func() interface{} { return &Context{} }},
		upgradePool:   &sync.Pool{New: func() interface{} { return &upgrade{} }},
		upgradeBuffer: make([]byte, 1024),
		numWorkers:    numCPU,
	}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

func (server *Server) Register(obj interface{}) error {
	return server.Registry.Register(obj)
}

func (server *Server) RegisterName(name string, obj interface{}) error {
	return server.Registry.RegisterName(name, obj)
}

func (server *Server) SetPipelining(enable bool) {
	server.pipelining = enable
}

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
		err := server.ServeRequest(codec, sending, wg, worker, ch, done, workers)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}
	wg.Wait()
	codec.Close()
	if done != nil {
		close(done)
	}
}

func (server *Server) ServeRequest(codec ServerCodec, sending *sync.Mutex, wg *sync.WaitGroup, worker bool, ch chan *Context, done chan struct{}, workers chan struct{}) error {
	ctx, err := server.readRequest(codec)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF && err != poll.EAGAIN {
			logger.Errorln("rpc:", err)
		}
		if !ctx.keepReading {
			ctx.Reset()
			server.ctxPool.Put(ctx)
			return err
		}
		ctx.Error = err.Error()
		server.sendResponse(sending, ctx, codec)
		return err
	}
	if ctx.heartbeat {
		server.sendResponse(sending, ctx, codec)
		return nil
	}
	if server.pipelining {
		server.callService(sending, nil, ctx, codec)
		return nil
	}
	if wg != nil {
		wg.Add(1)
	}
	if worker {
		select {
		case ch <- ctx:
		case workers <- struct{}{}:
			go func(done chan struct{}, ch chan *Context, server *Server, sending *sync.Mutex, wg *sync.WaitGroup, codec ServerCodec, workers chan struct{}, ctx *Context) {
				defer func() { <-workers }()
				for {
					server.callService(sending, wg, ctx, codec)
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
			}(done, ch, server, sending, wg, codec, workers, ctx)
		default:
			go server.callService(sending, wg, ctx, codec)
		}
	} else {
		go server.callService(sending, wg, ctx, codec)
	}
	return nil
}

func (server *Server) readRequest(codec ServerCodec) (ctx *Context, err error) {
	ctx = server.ctxPool.Get().(*Context)
	err = codec.ReadRequestHeader(ctx)
	if err != nil {
		return
	}
	ctx.keepReading = true
	if len(ctx.Upgrade) > 0 {
		u := server.getUpgrade()
		u.Unmarshal(ctx.Upgrade)
		if u.Heartbeat == Heartbeat {
			ctx.heartbeat = true
		}
		if u.NoRequest == NoRequest {
			ctx.noRequest = true
		}
		if u.NoResponse == NoResponse {
			ctx.noResponse = true
		}
		server.putUpgrade(u)
	}
	if !ctx.heartbeat {
		ctx.f = server.Registry.GetFunc(ctx.ServiceMethod)
		if ctx.f == nil {
			err = errors.New("rpc: can't find service " + ctx.ServiceMethod)
			codec.ReadRequestBody(nil)
			return
		}
	}
	if !ctx.noRequest {
		ctx.args = ctx.f.GetValueIn(0)
		if ctx.args == funcs.ZeroValue {
			err = errors.New("rpc: can't find args")
			codec.ReadRequestBody(nil)
			return
		}
		if err = codec.ReadRequestBody(ctx.args.Interface()); err != nil {
			return
		}
	}
	if !ctx.noResponse {
		ctx.reply = ctx.f.GetValueIn(1)
		if ctx.reply == funcs.ZeroValue {
			err = errors.New("rpc: can't find reply")
		}
	}
	return
}

func (server *Server) callService(sending *sync.Mutex, wg *sync.WaitGroup, ctx *Context, codec ServerCodec) {
	if wg != nil {
		defer wg.Done()
	}
	if err := ctx.f.ValueCall(ctx.args, ctx.reply); err != nil {
		ctx.Error = err.Error()
	}
	server.sendResponse(sending, ctx, codec)
}

func (server *Server) sendResponse(sending *sync.Mutex, ctx *Context, codec ServerCodec) {
	sending.Lock()
	var reply interface{}
	if len(ctx.Error) == 0 && !ctx.noResponse {
		reply = ctx.reply.Interface()
	}
	err := codec.WriteResponse(ctx, reply)
	if err != nil {
		logger.Allln("rpc: writing response:", err)
	}
	sending.Unlock()
	ctx.Reset()
	server.ctxPool.Put(ctx)
}

func (server *Server) listen(sock socket.Socket, address string, New NewServerCodecFunc) error {
	logger.Noticef("pid - %d", os.Getpid())
	if server.poll {
		logger.Noticef("poll - %s", poll.Tag)
	} else {
		logger.Noticef("poll - %s", "disabled")
	}
	logger.Noticef("network - %s", sock.Scheme())
	logger.Noticef("listening on %s", address)
	lis, err := sock.Listen(address)
	if err != nil {
		return err
	}
	type ServerContext struct {
		codec   ServerCodec
		sending *sync.Mutex
		worker  bool
		wg      *sync.WaitGroup
		ch      chan *Context
		done    chan struct{}
		workers chan struct{}
		closed  int32
	}
	if server.poll {
		lis.ServeMessages(func(messages socket.Messages) (socket.Context, error) {
			return &ServerContext{
				codec:   New(messages),
				sending: new(sync.Mutex),
				wg:      new(sync.WaitGroup),
			}, nil
		}, func(context socket.Context) error {
			ctx := context.(*ServerContext)
			err := server.ServeRequest(ctx.codec, ctx.sending, ctx.wg, ctx.worker, ctx.ch, ctx.done, ctx.workers)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if atomic.CompareAndSwapInt32(&ctx.closed, 0, 1) {
					ctx.wg.Wait()
					ctx.codec.Close()
					if ctx.done != nil {
						close(ctx.done)
					}
				}
			}
			return err
		})
		return nil
	} else {
		for {
			conn, err := lis.Accept()
			if err != nil {
				continue
			}
			go server.ServeCodec(New(conn.Messages()))
		}
	}
}

func (server *Server) Listen(network, address string, codec string) error {
	if newSocket := NewSocket(network); newSocket != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return server.listen(newSocket(), address, func(messages socket.Messages) ServerCodec {
				return NewServerCodec(newCodec(), nil, messages)
			})
		}
		return errors.New("unsupported codec: " + codec)
	}
	return errors.New("unsupported protocol scheme: " + network)
}

func (server *Server) ListenWithOptions(address string, opts *Options) error {
	if opts.NewCodec == nil && opts.NewEncoder == nil && opts.Codec == "" {
		return errors.New("need opts.NewCodec, opts.NewEncoder or opts.Codec")
	}
	if opts.NewSocket == nil && opts.Network == "" {
		return errors.New("need opts.NewSocket, opts.NewMessages or opts.Network")
	}
	var sock socket.Socket
	if newSocket := NewSocket(opts.Network); newSocket != nil {
		sock = newSocket()
	} else if opts.NewSocket != nil {
		sock = opts.NewSocket()
	}
	return server.listen(sock, address, func(messages socket.Messages) ServerCodec {
		var bodyCodec codec.Codec
		if newCodec := NewCodec(opts.Codec); newCodec != nil {
			bodyCodec = newCodec()
		} else if opts.NewCodec != nil {
			bodyCodec = opts.NewCodec()
		}
		var headerEncoder *encoder.Encoder
		if newEncoder := NewEncoder(opts.Encoder); newEncoder != nil {
			headerEncoder = newEncoder()
		} else if opts.NewEncoder != nil {
			headerEncoder = opts.NewEncoder()
		}
		return NewServerCodec(bodyCodec, headerEncoder, messages)
	})
}

func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

func RegisterName(name string, rcvr interface{}) error {
	return DefaultServer.RegisterName(name, rcvr)
}

func SetPipelining(enable bool) {
	DefaultServer.SetPipelining(enable)
}

func SetPoll(enable bool) {
	DefaultServer.poll = enable
}

func ServeCodec(codec ServerCodec) {
	DefaultServer.ServeCodec(codec)
}
