package rpc

import (
	"fmt"
	"github.com/hslam/funcs"
	"github.com/hslam/protocol"
	"io"
	"sync"
	"time"
)

var (
	//DefaultServer defines the default server
	DefaultServer = NewServer()
)

//Server defines the struct of server
type Server struct {
	network      string
	listener     Listener
	Registry     *funcs.Funcs
	timeout      int64
	batching     bool
	pipelining   bool
	multiplexing bool
	asyncMax     int
	noDelay      bool
}

//NewServer creates a new server
func NewServer() *Server {
	return &Server{Registry: funcs.New(), timeout: DefaultServerTimeout, asyncMax: DefaultMaxAsyncPerConn, multiplexing: true, batching: false}
}

//SetBatching enables batching.
func SetBatching(enable bool) {
	DefaultServer.SetBatching(enable)
}

//SetBatching enables batching.
func (s *Server) SetBatching(enable bool) {
	s.batching = enable
}

//SetNoDelay enables no delay.
func SetNoDelay(enable bool) {
	DefaultServer.SetNoDelay(enable)
}

//SetNoDelay enables no delay.
func (s *Server) SetNoDelay(enable bool) {
	s.noDelay = enable
}

//SetPipelining enables pipelining.
func SetPipelining(enable bool) {
	DefaultServer.SetPipelining(enable)
}

//SetPipelining enables pipelining.
func (s *Server) SetPipelining(enable bool) {
	s.pipelining = enable
	s.multiplexing = !enable
}

//SetMultiplexing enables multiplexing.
func SetMultiplexing(enable bool) {
	DefaultServer.SetMultiplexing(enable)
}

//SetMultiplexing enables multiplexing.
func (s *Server) SetMultiplexing(enable bool) {
	s.multiplexing = enable
	s.pipelining = !enable
	s.SetSize(DefaultMaxMultiplexingPerConn)
}

//SetSize sets the size of max async requests.
func SetSize(size int) {
	DefaultServer.SetSize(size)
}

//SetSize sets the size of max async requests.
func (s *Server) SetSize(size int) {
	s.asyncMax = size
}

// Register publishes in the server the set of methods.
func Register(obj interface{}) error {
	return DefaultServer.Register(obj)
}

// Register publishes in the server the set of methods.
func (s *Server) Register(obj interface{}) error {
	return s.Registry.Register(obj)
}

// RegisterName is like Register but uses the provided name.
func RegisterName(name string, obj interface{}) error {
	return DefaultServer.RegisterName(name, obj)
}

// RegisterName is like Register but uses the provided name.
func (s *Server) RegisterName(name string, obj interface{}) error {
	return s.Registry.RegisterName(name, obj)
}

// ListenAndServe listens on the network address addr and then calls Serve.
func ListenAndServe(network, address string) error {
	return DefaultServer.ListenAndServe(network, address)
}

// ListenAndServe listens on the network address addr and then calls Serve.
func (s *Server) ListenAndServe(network, address string) error {
	s.network = network
	listener, err := Listen(network, address, s)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	s.listener = listener
	err = listener.Serve()
	if err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}

// Close closes listener
func Close() error {
	return DefaultServer.Close()
}

// Close closes listener
func (s *Server) Close() error {
	return s.listener.Close()
}

// ServeMessage serves message
func ServeMessage(ReadWriteCloser io.ReadWriteCloser) error {
	return DefaultServer.ServeMessage(ReadWriteCloser)
}

// ServeMessage serves message
func (s *Server) ServeMessage(ReadWriteCloser io.ReadWriteCloser) error {
	return s.serve(ReadWriteCloser, false)
}

// ServeConn serves conn
func ServeConn(ReadWriteCloser io.ReadWriteCloser) error {
	return DefaultServer.ServeConn(ReadWriteCloser)
}

// ServeConn serves conn
func (s *Server) ServeConn(ReadWriteCloser io.ReadWriteCloser) error {
	return s.serve(ReadWriteCloser, true)
}
func (s *Server) serve(ReadWriteCloser io.ReadWriteCloser, Stream bool) error {
	readChan := make(chan []byte, 1)
	writeChan := make(chan []byte, 1)
	finishChan := make(chan bool, 2)
	stopReadChan := make(chan bool, 1)
	stopWriteChan := make(chan bool, 1)
	stopChan := make(chan bool, 1)
	if Stream {
		go protocol.ReadStream(ReadWriteCloser, readChan, stopReadChan, finishChan)
		var noDelay bool = true
		if !s.batching && (s.pipelining || s.multiplexing) {
			noDelay = false
		}
		if s.noDelay {
			noDelay = true
		}
		go protocol.WriteStream(ReadWriteCloser, writeChan, stopWriteChan, finishChan, noDelay)
	} else {
		go protocol.ReadConn(ReadWriteCloser, readChan, stopReadChan, finishChan)
		go protocol.WriteConn(ReadWriteCloser, writeChan, stopWriteChan, finishChan)
	}
	if s.multiplexing {
		jobChan := make(chan bool, s.asyncMax)
		for {
			select {
			case data := <-readChan:
				go func(data []byte, writeChan chan []byte) {
					defer func() {
						if err := recover(); err != nil {
						}
						<-jobChan
					}()
					jobChan <- true
					priority, id, body, err := protocol.UnpackFrame(data)
					if err != nil {
						return
					}
					_, resBytes := s.Serve(body)
					if resBytes != nil {
						frameBytes := protocol.PacketFrame(priority, id, resBytes)
						writeChan <- frameBytes
					}
				}(data, writeChan)
			case stop := <-finishChan:
				if stop {
					stopReadChan <- true
					stopWriteChan <- true
					goto endfor
				}
			}
		}
	} else {
		for {
			select {
			case data := <-readChan:
				_, resBytes := s.Serve(data)
				if resBytes != nil {
					writeChan <- resBytes
				}
			case stop := <-finishChan:
				if stop {
					stopReadChan <- true
					stopWriteChan <- true
					goto endfor
				}
			}
		}
	}
endfor:
	defer ReadWriteCloser.Close()
	close(readChan)
	close(writeChan)
	close(finishChan)
	close(stopReadChan)
	close(stopWriteChan)
	close(stopChan)
	return ErrConnExit
}

// Serve serves bytes
func Serve(b []byte) (bool, []byte) {
	return DefaultServer.Serve(b)
}

// Serve serves bytes
func (s *Server) Serve(b []byte) (ok bool, body []byte) {
	ctx := &serverCodec{}
	ctx.Decode(b)
	if ctx.msg.msgType == MsgTypeHea {
		return true, b
	}
	ctx.msg.msgType = MsgTypeRes
	var noResponse = false
	var responseBytes []byte
	if ctx.msg.batch {
		NoResponseCnt := &count{}
		if ctx.batchCodec.async {
			waitGroup := sync.WaitGroup{}
			for i, v := range ctx.requests {
				waitGroup.Add(1)
				go func(i int, req *request, res *response, NoResponseCnt *count, waitGroup *sync.WaitGroup) {
					defer waitGroup.Done()
					s.handle(ctx, req, res)
					if req.noResponse == true {
						NoResponseCnt.add(1)
						ctx.responses[i].data = []byte("")
					} else {
						body, _ := ctx.responses[i].Encode()
						ctx.batchCodec.data[i] = body
					}
				}(i, v, ctx.responses[i], NoResponseCnt, &waitGroup)
			}
			waitGroup.Wait()
		} else {
			for i, v := range ctx.requests {
				s.handle(ctx, v, ctx.responses[i])
				if v.noResponse == true {
					NoResponseCnt.add(1)
					ctx.responses[i].data = []byte("")
				} else {
					body, _ := ctx.responses[i].Encode()
					ctx.batchCodec.data[i] = body
				}
			}
		}
		if NoResponseCnt.load() == int64(len(ctx.requests)) {
			noResponse = true
		} else {
			responseBytes, _ = ctx.batchCodec.Encode()
		}
	} else {
		s.handle(ctx, ctx.request, ctx.response)
		noResponse = ctx.request.noResponse
		responseBytes, _ = ctx.response.Encode()
	}
	if noResponse == true {
		return true, nil
	}
	ctx.msg.data = responseBytes
	body, _ = ctx.Encode()
	return true, body
}

func (s *Server) handle(ctx *serverCodec, req *request, res *response) {
	res.id = req.id
	if s.timeout > 0 {
		ch := make(chan int)
		go func() {
			s.callService(ctx, req, res)
			ch <- 1
		}()
		select {
		case <-ch:
			return
		case <-time.After(time.Millisecond * time.Duration(s.timeout)):
			res.err = fmt.Errorf("method %s time out", req.method)
			return
		}
	}
	s.callService(ctx, req, res)
}
func (s *Server) callService(ctx *serverCodec, req *request, res *response) {
	if s.Registry.GetFunc(req.method) == nil {
		logger.Noticef("Server.CallService method %s is not supposted", req.method)
		res.err = fmt.Errorf("Server.CallService method %s is not supposted", req.method)
		return
	}
	if req.noRequest && req.noResponse {
		if err := s.Registry.Call(req.method); err != nil {
			logger.Noticef("Server.CallService OnlyCall err %s", err)
			res.err = fmt.Errorf("Server.CallService OnlyCall err %s", err)
			return
		}
		return
	} else if req.noRequest && !req.noResponse {
		reply := s.Registry.GetFuncIn(req.method, 0)
		if err := s.Registry.Call(req.method, reply); err != nil {
			logger.Noticef("Server.CallService CallNoRequest err %s", err)
			res.err = fmt.Errorf("Server.CallService CallNoRequest err %s", err)
			return
		}
		replyBytes, err := replyEncode(reply, ctx.msg.codecType)
		if err != nil {
			res.err = fmt.Errorf("Server.CallService ReplyEncode err %s", err)
			return
		}
		res.data = replyBytes
		return
	} else if !req.noRequest && req.noResponse {
		args := s.Registry.GetFuncIn(req.method, 0)
		err := argsDecode(req.data, args, ctx.msg.codecType)
		if err != nil {
			res.err = fmt.Errorf("Server.CallService ArgsDecode err %s", err)
			return
		}
		reply := s.Registry.GetFuncIn(req.method, 1)
		if reply != nil {
			if err := s.Registry.Call(req.method, args, reply); err != nil {
				res.err = fmt.Errorf("Server.CallService CallNoResponseWithReply err %s", err)
				return
			}
		} else {
			if err := s.Registry.Call(req.method, args); err != nil {
				res.err = fmt.Errorf("Server.CallService CallNoResponseWithoutReply err %s", err)
				return
			}
		}
		return

	} else {
		args := s.Registry.GetFuncIn(req.method, 0)
		err := argsDecode(req.data, args, ctx.msg.codecType)
		if err != nil {
			res.err = fmt.Errorf("Server.CallService ArgsDecode err %s", err)
			return
		}
		reply := s.Registry.GetFuncIn(req.method, 1)
		if err := s.Registry.Call(req.method, args, reply); err != nil {
			res.err = fmt.Errorf("Server.CallService Call err %s", err)
			return
		}
		var replyBytes []byte
		replyBytes, err = replyEncode(reply, ctx.msg.codecType)
		if err != nil {
			res.err = fmt.Errorf("Server.CallService ReplyEncode err %s", err)
			return
		}
		res.data = replyBytes
		return
	}
}
