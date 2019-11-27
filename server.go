package rpc

import (
	"hslam.com/git/x/funcs"
	"hslam.com/git/x/protocol"
	"fmt"
	"time"
	"sync"
	"io"
)

var (
	DefaultServer = NewServer()
)

type registerObject struct {
	name string
	obj	interface{}
}
type Server struct {
	network 					string
	listener 					Listener
	Funcs 						*funcs.Funcs
	timeout 					int64
	batch 						bool
	pipelining 					bool
	multiplexing				bool
	pipeliningAsync				bool
	asyncMax					int
	lowDelay 					bool
	objs 						[]*registerObject
}
func NewServer() *Server {
	return &Server{Funcs:funcs.New(),timeout:DefaultServerTimeout,asyncMax:DefaultMaxAsyncPerConn,multiplexing:true,batch:true}
}
func SetBatch(enable bool)  {
	DefaultServer.SetBatch(enable)
}
func (s *Server) SetBatch(enable bool)  {
	s.batch=enable
}
func SetLowDelay(enable bool)  {
	DefaultServer.SetLowDelay(enable)
}
func (s *Server) SetLowDelay(enable bool)  {
	s.lowDelay=enable
}
func SetPipelining(enable bool) {
	DefaultServer.SetPipelining(enable)
}
func (s *Server) SetPipelining(enable bool) {
	s.pipelining=enable
}
func SetMultiplexing(enable bool)  {
	DefaultServer.SetMultiplexing(enable)
}
func (s *Server) SetMultiplexing(enable bool) {
	if enable{
		s.EnableMultiplexingWithSize(DefaultMaxMultiplexingPerConn)
	}else {
		s.multiplexing=false
	}
}
func EnableMultiplexingWithSize(size  int)  {
	DefaultServer.EnableMultiplexingWithSize(size)
}
func (s *Server) EnableMultiplexingWithSize(size  int) {
	s.multiplexing=true
	s.asyncMax=size
}
func PipeliningAsync() bool {
	return DefaultServer.PipeliningAsync()
}
func (s *Server) PipeliningAsync() bool {
	return s.pipeliningAsync
}
func SetPipeliningAsync(enable bool)  {
	DefaultServer.SetPipeliningAsync(enable)
}
func (s *Server) SetPipeliningAsync(enable bool) {
	if enable{
		s.EnablePipeliningAsyncWithSize(DefaultMaxAsyncPerConn)
	}else {
		s.pipeliningAsync=false
	}
}

func EnablePipeliningAsyncWithSize(size  int) {
	DefaultServer.EnablePipeliningAsyncWithSize(size )
}

func (s *Server) EnablePipeliningAsyncWithSize(size int) {
	s.pipeliningAsync=true
	s.asyncMax=size
}
func Register(obj interface{}) error {
	return DefaultServer.Register(obj)
}
func (s *Server) Register(obj interface{}) error {
	s.objs=append(s.objs, &registerObject{"",obj})
	return s.Funcs.Register(obj)
}

func RegisterName(name string, obj interface{}) error {
	return DefaultServer.RegisterName(name, obj)
}

func (s *Server) RegisterName(name string, obj interface{}) error {
	s.objs=append(s.objs, &registerObject{name,obj})
	return s.Funcs.RegisterName(name, obj)
}

func ListenAndServe(network,address string) error {
	return DefaultServer.ListenAndServe(network,address)
}

func (s *Server)ListenAndServe(network,address string) error {
	s.network=network
	listener,err:=Listen(network,address,s)
	if err != nil {
		Errorln(err)
		return err
	}
	s.listener=listener
	err=s.listener.Serve()
	if err != nil {
		Errorln(err)
		return err
	}
	return nil
}
func ServeMessage(ReadWriteCloser io.ReadWriteCloser) error {
	return DefaultServer.ServeMessage(ReadWriteCloser)
}
func (s *Server) ServeMessage(ReadWriteCloser io.ReadWriteCloser) error {
	return s.serve(ReadWriteCloser,false)
}
func ServeConn(ReadWriteCloser io.ReadWriteCloser) error {
	return DefaultServer.ServeConn(ReadWriteCloser)
}
func (s *Server) ServeConn(ReadWriteCloser io.ReadWriteCloser) error {
	return s.serve(ReadWriteCloser,true)
}
func (s *Server) serve(ReadWriteCloser io.ReadWriteCloser,Stream bool) error {
	readChan := make(chan []byte,1)
	writeChan := make(chan []byte,1)
	finishChan:= make(chan bool,2)
	stopReadChan := make(chan bool,1)
	stopWriteChan := make(chan bool,1)
	stopChan := make(chan bool,1)
	if Stream{
		go protocol.ReadStream(ReadWriteCloser, readChan, stopReadChan,finishChan)
		var useBuffer bool
		if !s.batch&&!s.lowDelay&&(s.pipelining||s.multiplexing){
			useBuffer=true
		}
		go protocol.WriteStream(ReadWriteCloser, writeChan, stopWriteChan,finishChan,useBuffer)
	}else {
		go protocol.ReadConn(ReadWriteCloser, readChan, stopReadChan,finishChan)
		go protocol.WriteConn(ReadWriteCloser, writeChan, stopWriteChan,finishChan)
	}
	if s.multiplexing{
		jobChan := make(chan bool,s.asyncMax)
		for {
			select {
			case data := <-readChan:
				go func(data []byte ,writeChan chan []byte) {
					defer func() {
						if err := recover(); err != nil {
						}
						<-jobChan
					}()
					jobChan<-true
					priority,id,body,err:=protocol.UnpackFrame(data)
					if err!=nil{
						return
					}
					_,res_bytes, _ := s.Serve(body)
					if res_bytes!=nil{
						frameBytes:=protocol.PacketFrame(priority,id,res_bytes)
						writeChan <- frameBytes
					}
				}(data,writeChan)
			case stop := <-finishChan:
				if stop {
					stopReadChan<-true
					stopWriteChan<-true
					goto endfor
				}
			}
		}
	}else if s.pipeliningAsync{
		syncConn:=newSyncConn(s)
		go protocol.HandleSyncConn(syncConn, writeChan,readChan,stopChan,s.asyncMax)
		select {
		case stop := <-finishChan:
			if stop {
				stopReadChan<-true
				stopWriteChan<-true
				stopChan<-true
				goto endfor
			}
		}
	}else {
		for {
			select {
			case data := <-readChan:
				_,res_bytes, _ := s.Serve(data)
				if res_bytes!=nil{
					writeChan <- res_bytes
				}
			case stop := <-finishChan:
				if stop {
					stopReadChan<-true
					stopWriteChan<-true
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

func Serve(b []byte) (bool,[]byte,error) {
	return DefaultServer.Serve(b)
}
func (s *Server)Serve(b []byte) (bool,[]byte,error) {
	ctx:=&ServerCodec{}
	ctx.Decode(b)
	if ctx.msg.msgType==MsgTypeHea{
		return true,b,nil
	}
	ctx.msg.msgType=MsgTypeRes
	var noResponse =false
	var responseBytes []byte
	if ctx.msg.batch{
		NoResponseCnt:=&Count{}
		if ctx.batchCodec.async{
			waitGroup :=sync.WaitGroup{}
			for i,v:=range ctx.requests{
				waitGroup.Add(1)
				go func(i int,req *Request,res *Response,NoResponseCnt *Count,waitGroup *sync.WaitGroup) {
					defer waitGroup.Done()
					s.Handler(ctx,req,res)
					if req.noResponse==true{
						NoResponseCnt.add(1)
						ctx.responses[i].data=[]byte("")
					}else{
						body,_:=ctx.responses[i].Encode()
						ctx.batchCodec.data[i]=body
					}
				}(i,v,ctx.responses[i],NoResponseCnt,&waitGroup)
			}
			waitGroup.Wait()
		}else {
			for i,v:=range ctx.requests{
				s.Handler(ctx,v,ctx.responses[i])
				if v.noResponse==true{
					NoResponseCnt.add(1)
					ctx.responses[i].data=[]byte("")
				}else{
					body,_:=ctx.responses[i].Encode()
					ctx.batchCodec.data[i]=body
				}
			}
		}
		if NoResponseCnt.load()==int64(len(ctx.requests)){
			noResponse=true
		}else {
			responseBytes,_=ctx.batchCodec.Encode()
		}
	}else{
		s.Handler(ctx,ctx.request,ctx.response)
		noResponse=ctx.request.noResponse
		responseBytes,_=ctx.response.Encode()
	}
	if noResponse==true{
		return true,nil,nil
	}
	ctx.msg.data=responseBytes
	body,err:=ctx.Encode()
	return true,body,err
}

func (s *Server)Handler(ctx *ServerCodec,req *Request,res *Response){
	res.id=req.id
	if s.timeout>0{
		ch := make(chan int)
		go func() {
			s.CallService(ctx,req,res)
			ch<-1
		}()
		select {
		case <-ch:
			return
		case <-time.After(time.Millisecond * time.Duration(s.timeout)):
			res.err=fmt.Errorf("method %s time out",req.method)
			return
		}
	}
	s.CallService(ctx,req,res)
}
func (s *Server)CallService(ctx *ServerCodec,req *Request,res *Response) {
	if s.Funcs.GetFunc(req.method)==nil{
		AllInfof("Server.CallService method %s is not supposted",req.method)
		res.err=fmt.Errorf("Server.CallService method %s is not supposted",req.method)
		return
	}
	if req.noRequest && req.noResponse{
		if err := s.Funcs.Call(req.method); err != nil {
			AllInfof("Server.CallService OnlyCall err %s",err)
			res.err=fmt.Errorf("Server.CallService OnlyCall err %s",err)
			return
		}
		return
	}else if req.noRequest && !req.noResponse{
		reply :=s.Funcs.GetFuncIn(req.method,0)
		if err := s.Funcs.Call(req.method, reply); err != nil {
			AllInfof("Server.CallService CallNoRequest err %s",err)
			res.err= fmt.Errorf("Server.CallService CallNoRequest err %s",err)
			return
		}
		reply_bytes,err:=ReplyEncode(reply,ctx.msg.codecType)
		if err!=nil{
			res.err= fmt.Errorf("Server.CallService ReplyEncode err %s",err)
			return
		}
		res.data=reply_bytes
		return
	}else if !req.noRequest && req.noResponse{
		args := s.Funcs.GetFuncIn(req.method,0)
		err:=ArgsDecode(req.data,args,ctx.msg.codecType)
		if err!=nil{
			res.err= fmt.Errorf("Server.CallService ArgsDecode err %s",err)
			return
		}
		reply :=s.Funcs.GetFuncIn(req.method,1)
		if reply!=nil{
			if err := s.Funcs.Call(req.method, args, reply); err != nil {
				res.err= fmt.Errorf("Server.CallService CallNoResponseWithReply err %s",err)
				return
			}
		}else {
			if err := s.Funcs.Call(req.method, args); err != nil {
				res.err= fmt.Errorf("Server.CallService CallNoResponseWithoutReply err %s",err)
				return
			}
		}
		return

	}else {
		args := s.Funcs.GetFuncIn(req.method,0)
		err:=ArgsDecode(req.data,args,ctx.msg.codecType)
		if err!=nil{
			res.err= fmt.Errorf("Server.CallService ArgsDecode err %s",err)
			return
		}
		reply :=s.Funcs.GetFuncIn(req.method,1)
		if err := s.Funcs.Call(req.method, args, reply); err != nil {
			res.err= fmt.Errorf("Server.CallService Call err %s",err)
			return
		}
		var reply_bytes []byte
		reply_bytes,err=ReplyEncode(reply,ctx.msg.codecType)
		if err!=nil{
			res.err= fmt.Errorf("Server.CallService ReplyEncode err %s",err)
			return
		}
		res.data=reply_bytes
		return
	}
}

