package rpc

import (
	"errors"
	"github.com/hslam/funcs"
	"github.com/hslam/transport"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

var numCPU = runtime.NumCPU()

type Server struct {
	Registry   *funcs.Funcs
	ctxPool    *sync.Pool
	pipelining bool
	numWorkers int
}
type NewServerCodecFunc func(conn io.ReadWriteCloser) ServerCodec

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		Registry:   funcs.New(),
		ctxPool:    &sync.Pool{New: func() interface{} { return &Context{} }},
		numWorkers: numCPU,
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

func (server *Server) ServeCodec(codec ServerCodec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	Count := int64(0)
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
		ctx.Count = &Count
		if server.pipelining {
			server.callService(sending, nil, ctx, codec)
			continue
		}
		if atomic.AddInt64(ctx.Count, 1) < int64(server.numWorkers+1) && server.numWorkers > 0 {
			wg.Add(1)
			ch <- ctx
			continue
		}
		wg.Add(1)
		go server.callService(sending, wg, ctx, codec)
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
	ctx.f = server.Registry.GetFunc(ctx.ServiceMethod)
	if ctx.f == nil {
		err = errors.New("rpc: can't find service " + ctx.ServiceMethod)
		codec.ReadRequestBody(nil)
		return
	}
	ctx.args = ctx.f.GetValueIn(0)
	if ctx.args == funcs.ZeroValue {
		err = errors.New("rpc: can't find args")
		codec.ReadRequestBody(nil)
		return
	}
	if err = codec.ReadRequestBody(ctx.args.Interface()); err != nil {
		return
	}
	ctx.reply = ctx.f.GetValueIn(1)
	if ctx.reply == funcs.ZeroValue {
		err = errors.New("rpc: can't find reply")
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
	err := codec.WriteResponse(ctx, ctx.reply.Interface())
	if err != nil {
		logger.Allln("rpc: writing response:", err)
	}
	sending.Unlock()
	atomic.AddInt64(ctx.Count, -1)
	ctx.Reset()
	server.ctxPool.Put(ctx)
}

func (server *Server) ServeConn(conn io.ReadWriteCloser, New NewServerCodecFunc) {
	server.ServeCodec(New(conn))
}

func (server *Server) Listen(tran transport.Transport, address string, New NewServerCodecFunc) error {
	logger.Noticef("pid - %d", os.Getpid())
	logger.Noticef("network - %s", tran.Scheme())
	logger.Noticef("listening on %s", address)
	lis, err := tran.Listen(address)
	if err != nil {
		return err
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go server.ServeConn(conn, New)
	}
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

func ServeConn(conn io.ReadWriteCloser, New NewServerCodecFunc) {
	DefaultServer.ServeConn(conn, New)
}
