// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package rpc provides access to the exported methods of an object across a network or other I/O connection.
package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/funcs"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/socket"
	"io"
	"os"
	"runtime"
	"sync"
)

var numCPU = runtime.NumCPU()

type Server struct {
	Registry      *funcs.Funcs
	ctxPool       *sync.Pool
	upgradePool   *sync.Pool
	upgradeBuffer []byte
	pipelining    bool
	numWorkers    int
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

func (server *Server) SetNumWorkers(num int) {
	server.numWorkers = num
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
	var done chan bool
	if !server.pipelining && server.numWorkers > 0 {
		ch = make(chan *Context, server.numWorkers)
		done = make(chan bool, 1)
		for i := 0; i < server.numWorkers; i++ {
			go func(done chan bool, ch chan *Context, server *Server, sending *sync.Mutex, wg *sync.WaitGroup, codec ServerCodec) {
				for {
					select {
					case ctx := <-ch:
						server.callService(sending, wg, ctx, codec)
					case <-done:
						return
					}
				}
			}(done, ch, server, sending, wg, codec)
		}
	}
	for {
		ctx, err := server.readRequest(codec)
		if err != nil {
			if err != io.EOF {
				logger.Errorln("rpc:", err)
			}
			if !ctx.keepReading {
				break
			}
			if ctx.decodeHeader {
				ctx.Error = err.Error()
				server.sendResponse(sending, ctx, codec)
			} else {
				ctx.Reset()
				server.ctxPool.Put(ctx)
			}
			continue
		}
		if ctx.heartbeat {
			server.sendResponse(sending, ctx, codec)
			continue
		}
		if server.pipelining {
			server.callService(sending, nil, ctx, codec)
			continue
		}
		wg.Add(1)
		if server.numWorkers > 0 {
			select {
			case ch <- ctx:
			default:
				go server.callService(sending, wg, ctx, codec)
			}
		} else {
			go server.callService(sending, wg, ctx, codec)
		}
	}
	wg.Wait()
	codec.Close()
	if done != nil {
		close(done)
	}
}

func (server *Server) readRequest(codec ServerCodec) (ctx *Context, err error) {
	ctx = server.ctxPool.Get().(*Context)
	ctx.decodeHeader = true
	err = codec.ReadRequestHeader(ctx)
	if err != nil {
		ctx.decodeHeader = false
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		err = errors.New("rpc: server cannot decode request: " + err.Error())
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

func (server *Server) listen(socket socket.Socket, address string, New NewServerCodecFunc) error {
	logger.Noticef("pid - %d", os.Getpid())
	logger.Noticef("network - %s", socket.Scheme())
	logger.Noticef("listening on %s", address)
	lis, err := socket.Listen(address)
	if err != nil {
		return err
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go server.ServeCodec(New(conn.Messages()))
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
	if opts.NewSocket == nil && opts.NewMessages == nil && opts.Network == "" {
		return errors.New("need opts.NewSocket, opts.NewMessages or opts.Network")
	}
	if opts.NewMessages != nil {
		if messages := opts.NewMessages(); messages == nil {
			return errors.New("NewMessages failed")
		} else {
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
			if codec := NewServerCodec(bodyCodec, headerEncoder, messages); codec == nil {
				return errors.New("NewClientCodec failed")
			} else {
				server.ServeCodec(codec)
				return nil
			}
		}
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

func SetNumWorkers(num int) {
	DefaultServer.SetNumWorkers(num)
}

func ServeCodec(codec ServerCodec) {
	DefaultServer.ServeCodec(codec)
}
