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
	msg:=&Msg{}
	err:=msg.Decode(b)
	if err!=nil{
		Warnf("CallService MrpcDecode error: %s", err)
		return s.RPCErrEncode(msg.batch,err)
	}
	if msg.version!=Version{
		return s.RPCErrEncode(msg.batch,fmt.Errorf("%d %d Version is not matched",Version,msg.version))
	}
	if msg.msgType==MsgType(MsgTypeHea){
		return true,b,err
	}
	var noResponse =false
	var responseBytes []byte
	if msg.batch{
		batchCodec:=&BatchCodec{}
		batchCodec.Decode(msg.data)
		res_bytes_s:=make([][]byte,len(batchCodec.data))
		NoResponseCnt:=&Count{}
		if batchCodec.async{
			waitGroup :=sync.WaitGroup{}
			for i,v:=range batchCodec.data{
				waitGroup.Add(1)
				go func(i int,v []byte,msg Msg,res_bytes_s[][]byte,NoResponseCnt *Count,waitGroup *sync.WaitGroup) {
					defer waitGroup.Done()
					msg.data=v
					res_bytes,nr:=s.Handler(&msg)
					if nr==true{
						NoResponseCnt.add(1)
						res_bytes_s[i]=[]byte("")
					}else {
						res_bytes_s[i]=res_bytes
					}
				}(i,v,*msg,res_bytes_s,NoResponseCnt,&waitGroup)
			}
			waitGroup.Wait()
		}else {
			for i,v:=range batchCodec.data{
				msg.data=v
				res_bytes,nr:=s.Handler(msg)
				if nr==true{
					NoResponseCnt.add(1)
					res_bytes_s[i]=[]byte("")
				}else {
					res_bytes_s[i]=res_bytes
				}
			}
		}
		if NoResponseCnt.load()==int64(len(batchCodec.data)){
			noResponse=true
		}else {
			batchCodec:=&BatchCodec{data:res_bytes_s}
			responseBytes,_=batchCodec.Encode()
		}
	}else{
		responseBytes,noResponse=s.Handler(msg)
	}
	if noResponse==true{
		return true,nil,nil
	}
	msg.data=responseBytes
	msg.msgType=MsgType(MsgTypeRes)
	msg_bytes,err:=msg.Encode()
	return true,msg_bytes,err
}

func (s *Server)Handler(msg *Msg) ([]byte,bool) {
	req:=&Request{}
	err:=req.Decode(msg.data)
	if err!=nil{
		Warnf("id:%d  req:%d RequestDecode error: %s ",msg.id, req.id,err)
		return s.ResponseErrEncode(req.id,err),req.noResponse
	}
	data,ok:=s.Middleware(req,msg.codecType)
	if ok{
		AllInfof("id:%d  req:%d CallService %s success",msg.id, req.id,req.method)
		if req.noResponse==true{
			return nil,true
		}
		res:=&Response{}
		res.id=req.id
		res.data=data
		res_bytes,_:=res.Encode()
		if err!=nil{
			Warnf("id:%d  req:%d ResponseEncode error: %s ",msg.id, req.id, err)
			return s.ResponseErrEncode(req.id,err),req.noResponse
		}
		return res_bytes,req.noResponse
	}
	return data,req.noResponse
}

func (s *Server)Middleware(req *Request,funcsCodecType CodecType) ([]byte,bool){
	if s.timeout>0{
		ch := make(chan int)
		var (
			data []byte
			ok bool
		)
		go func() {
			data,ok=s.CallService(req,funcsCodecType)
			ch<-1
		}()
		select {
		case <-ch:
			return data,ok
		case <-time.After(time.Millisecond * time.Duration(s.timeout)):
			return s.ResponseErrEncode(req.id,fmt.Errorf("method %s time out",req.method)),false
		}
	}
	return s.CallService(req,funcsCodecType)
}
func (s *Server)CallService(req *Request,funcsCodecType CodecType) ([]byte,bool) {
	if s.Funcs.GetFunc(req.method)==nil{
		AllInfof("CallService %s is not supposted",req.method)
		return s.ResponseErrEncode(req.id,fmt.Errorf("method %s is not supposted",req.method)),false
	}
	if req.noRequest && req.noResponse{
		if err := s.Funcs.Call(req.method); err != nil {
			return s.ResponseErrEncode(req.id,err),false
		}
		return nil,true
	}else if req.noRequest && !req.noResponse{
		reply :=s.Funcs.GetFuncIn(req.method,0)
		if err := s.Funcs.Call(req.method, reply); err != nil {
			return s.ResponseErrEncode(req.id,err),false
		}
		reply_bytes,err:=ReplyEncode(reply,funcsCodecType)
		if err!=nil{
			return s.ResponseErrEncode(req.id,err),false
		}
		return reply_bytes,true
	}else if !req.noRequest && req.noResponse{
		args := s.Funcs.GetFuncIn(req.method,0)
		err:=ArgsDecode(req.data,args,funcsCodecType)
		if err!=nil{
			return s.ResponseErrEncode(req.id,err),false
		}
		reply :=s.Funcs.GetFuncIn(req.method,1)
		if reply!=nil{
			if err := s.Funcs.Call(req.method, args, reply); err != nil {
				return s.ResponseErrEncode(req.id,err),false
			}
		}else {
			if err := s.Funcs.Call(req.method, args); err != nil {
				return s.ResponseErrEncode(req.id,err),false
			}
		}
		return nil,true

	}else {
		args := s.Funcs.GetFuncIn(req.method,0)
		err:=ArgsDecode(req.data,args,funcsCodecType)
		if err!=nil{
			return s.ResponseErrEncode(req.id,err),false
		}
		reply :=s.Funcs.GetFuncIn(req.method,1)
		if err := s.Funcs.Call(req.method, args, reply); err != nil {
			return s.ResponseErrEncode(req.id,err),false
		}
		var reply_bytes []byte
		reply_bytes,err=ReplyEncode(reply,funcsCodecType)
		if err!=nil{
			return s.ResponseErrEncode(req.id,err),false
		}
		return reply_bytes,true
	}
}

func (s *Server)RPCErrEncode(batch bool,errMsg error)(bool,[]byte,error)  {
	res_bytes:=s.ResponseErrEncode(0, errMsg)
	msg:=&Msg{}
	msg.data=res_bytes
	msg.batch=batch
	msg.msgType=MsgType(MsgTypeRes)
	msg.codecType=FUNCS_CODEC_INVALID
	msg_bytes,err:=msg.Encode()
	return false,msg_bytes,err
}


func (s *Server)ResponseErrEncode(id uint64,err error)[]byte  {
	res:=&Response{}
	res.id=id
	res.err=err
	res_bytes,_:=res.Encode()
	return res_bytes
}