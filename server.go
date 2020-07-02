package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/funcs"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/transport"
	"io"
	"os"
	"sync"
)

type Server struct {
	Registry   *funcs.Funcs
	ctxPool    *sync.Pool
	pipelining bool
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		Registry: funcs.New(),
		ctxPool:  &sync.Pool{New: func() interface{} { return &Context{} }},
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
func (server *Server) ServeCodec(codec ServerCodec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
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
		if server.pipelining {
			server.callService(sending, nil, ctx, codec)
			continue
		}
		wg.Add(1)
		go server.callService(sending, wg, ctx, codec)
	}
	wg.Wait()
	codec.Close()
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
	ctx.Reset()
	server.ctxPool.Put(ctx)
}

func (server *Server) ServeConn(conn io.ReadWriteCloser, bodyCodec codec.Codec, headerEncoder *encoder.Encoder) {
	server.ServeCodec(NewServerCodec(conn, bodyCodec, headerEncoder))
}

func (server *Server) ListenAndServe(tran transport.Transport, address string, bodyCodec codec.Codec, headerEncoder *encoder.Encoder) {
	logger.Noticef("pid %d", os.Getpid())
	logger.Noticef("listening on %s", address)
	lis, err := tran.Listen(address)
	if err != nil {
		logger.Fatalln("fatal error: ", err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go server.ServeConn(conn, bodyCodec, headerEncoder)
	}
}

func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

func RegisterName(name string, rcvr interface{}) error {
	return DefaultServer.RegisterName(name, rcvr)
}

func SetPipelining(enable bool) {
	DefaultServer.SetPipelining(enable)
}

func ServeCodec(codec ServerCodec) {
	DefaultServer.ServeCodec(codec)
}

func ServeConn(conn io.ReadWriteCloser, bodyCodec codec.Codec, headerEncoder *encoder.Encoder) {
	DefaultServer.ServeConn(conn, bodyCodec, headerEncoder)
}
