package rpc

import (
	"hslam.com/mgit/Mort/funcs"
	"hslam.com/mgit/Mort/rpc/log"
	"fmt"
	"time"
	"sync"
)

var (

	DefaultServer = NewServer()
)

type Server struct {
	network 					string
	listener 					Listener
	Funcs 						*funcs.Funcs
	timeout 					int64
	multiplexing				bool
	async						bool
	asyncMax					int
}
func NewServer() *Server {
	return &Server{Funcs:funcs.New(),timeout:DefaultServerTimeout,asyncMax:DefaultMaxAsyncPerConn}
}

func EnableMultiplexing()  {
	DefaultServer.EnableMultiplexing()
}
func (s *Server) EnableMultiplexing() {
	s.EnableMultiplexingWithSize(DefaultMaxMultiplexingPerConn)
}
func EnableMultiplexingWithSize(size  int)  {
	DefaultServer.EnableMultiplexingWithSize(size)
}
func (s *Server) EnableMultiplexingWithSize(size  int) {
	s.multiplexing=true
	s.asyncMax=size
}
func Async() bool {
	return DefaultServer.Async()
}
func (s *Server) Async() bool {
	return s.async
}
func EnableAsyncHandle()  {
	DefaultServer.EnableAsyncHandle()
}
func (s *Server) EnableAsyncHandle() {
	s.EnableAsyncHandleWithSize(DefaultMaxAsyncPerConn)
}

func EnableAsyncHandleWithSize(size  int) {
	DefaultServer.EnableAsyncHandleWithSize(size )
}

func (s *Server) EnableAsyncHandleWithSize(size int) {
	s.async=true
	s.asyncMax=size
}
func Register(obj interface{}) error {
	return DefaultServer.Register(obj)
}
func (s *Server) Register(obj interface{}) error {
	return s.Funcs.Register(obj)
}

func RegisterName(name string, obj interface{}) error {
	return DefaultServer.RegisterName(name, obj)
}

func (s *Server) RegisterName(name string, obj interface{}) error {
	return s.Funcs.RegisterName(name, obj)
}

func ListenAndServe(network,address string) error {
	return DefaultServer.ListenAndServe(network,address)
}

func (s *Server)ListenAndServe(network,address string) error {
	s.network=network
	listener,err:=Listen(network,address,s)
	if err != nil {
		log.Errorln(err)
		return err
	}
	s.listener=listener
	err=s.listener.Serve()
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func ServeRPC(b []byte) (bool,[]byte,error) {
	return DefaultServer.ServeRPC(b)
}
func (s *Server)ServeRPC(b []byte) (bool,[]byte,error) {
	msg:=&Msg{}
	err:=msg.Decode(b)
	if err!=nil{
		log.Warnf("CallService MrpcDecode error: %s", err)
		return s.ErrRPCEncode(msg.batch,err)
	}
	if msg.version!=Version{
		return s.ErrRPCEncode(msg.batch,fmt.Errorf("%d %d Version is not matched",Version,msg.version))
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
		NoResponseCnt:=0
		if s.async||batchCodec.async{
			waitGroup :=sync.WaitGroup{}
			for i,v:=range batchCodec.data{
				waitGroup.Add(1)
				go func(i int,v []byte,msg Msg,res_bytes_s[][]byte,NoResponseCnt *int,waitGroup *sync.WaitGroup) {
					msg.data=v
					res_bytes,nr:=s.Handler(&msg)
					if nr==true{
						*NoResponseCnt++
						res_bytes_s[i]=[]byte("")
					}else {
						res_bytes_s[i]=res_bytes
					}
					waitGroup.Done()
				}(i,v,*msg,res_bytes_s,&NoResponseCnt,&waitGroup)
			}
			waitGroup.Wait()
		}else {
			for i,v:=range batchCodec.data{
				msg.data=v
				res_bytes,nr:=s.Handler(msg)
				if nr==true{
					NoResponseCnt++
					res_bytes_s[i]=[]byte("")
				}else {
					res_bytes_s[i]=res_bytes
				}
			}
		}
		if NoResponseCnt==len(batchCodec.data){
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
		log.Warnln("id-%d  req-%d RequestDecode error: %s ",msg.id, req.id,err)
		return s.ErrResponseEncode(req.id,err),false
	}
	data,ok:=s.Middleware(req,msg.codecType)
	if ok{
		log.AllInfof("id-%d  req-%d CallService %s success",msg.id, req.id,req.method)
		if req.noResponse==true{
			return nil,true
		}
		res:=&Response{}
		res.id=req.id
		res.data=data
		res_bytes,_:=res.Encode()
		if err!=nil{
			log.Warnln("id-%d  req-%d ResponseEncode error: %s ",msg.id, req.id, err)
			return s.ErrResponseEncode(req.id,err),req.noResponse
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
			return s.ErrResponseEncode(req.id,fmt.Errorf("method %s time out",req.method)),false
		}
	}
	return s.CallService(req,funcsCodecType)
}
func (s *Server)CallService(req *Request,funcsCodecType CodecType) ([]byte,bool) {
	if s.Funcs.GetFunc(req.method)==nil{
		log.AllInfof("CallService %s is not supposted",req.method)
		return s.ErrResponseEncode(req.id,fmt.Errorf("method %s is not supposted",req.method)),false
	}
	if req.noRequest && req.noResponse{
		if err := s.Funcs.Call(req.method); err != nil {
			return s.ErrResponseEncode(req.id,err),false
		}
		return nil,true
	}else if req.noRequest && !req.noResponse{
		reply :=s.Funcs.GetFuncIn(req.method,0)
		if err := s.Funcs.Call(req.method, reply); err != nil {
			return s.ErrResponseEncode(req.id,err),false
		}
		reply_bytes,err:=ReplyEncode(reply,funcsCodecType)
		if err!=nil{
			return s.ErrResponseEncode(req.id,err),false
		}
		return reply_bytes,true
	}else if !req.noRequest && req.noResponse{
		args := s.Funcs.GetFuncIn(req.method,0)
		err:=ArgsDecode(req.data,args,funcsCodecType)
		if err!=nil{
			return s.ErrResponseEncode(req.id,err),false
		}
		reply :=s.Funcs.GetFuncIn(req.method,1)
		if reply!=nil{
			if err := s.Funcs.Call(req.method, args, reply); err != nil {
				return s.ErrResponseEncode(req.id,err),false
			}
		}else {
			if err := s.Funcs.Call(req.method, args); err != nil {
				return s.ErrResponseEncode(req.id,err),false
			}
		}
		return nil,true

	}else {
		args := s.Funcs.GetFuncIn(req.method,0)
		err:=ArgsDecode(req.data,args,funcsCodecType)
		if err!=nil{
			return s.ErrResponseEncode(req.id,err),false
		}
		reply :=s.Funcs.GetFuncIn(req.method,1)
		if err := s.Funcs.Call(req.method, args, reply); err != nil {
			return s.ErrResponseEncode(req.id,err),false
		}
		var reply_bytes []byte
		reply_bytes,err=ReplyEncode(reply,funcsCodecType)
		if err!=nil{
			return s.ErrResponseEncode(req.id,err),false
		}
		return reply_bytes,true
	}
}

func (s *Server)ErrRPCEncode(batch bool,errMsg error)(bool,[]byte,error)  {
	res_bytes:=s.ErrResponseEncode(0, errMsg)
	msg:=&Msg{}
	msg.data=res_bytes
	msg.batch=batch
	msg.msgType=MsgType(MsgTypeRes)
	msg.codecType=FUNCS_CODEC_INVALID
	msg_bytes,err:=msg.Encode()
	return false,msg_bytes,err
}


func (s *Server)ErrResponseEncode(id uint64,err error)[]byte  {
	res:=&Response{}
	res.id=id
	res.err=err
	res_bytes,_:=res.Encode()
	return res_bytes
}
