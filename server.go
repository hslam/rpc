package rpc

import (
	"hslam.com/mgit/Mort/funcs"
	"hslam.com/mgit/Mort/rpc/log"
	"errors"
	"fmt"
	"time"
)

var DefaultServer = NewServer()

type Server struct {
	network string
	listener Listener
	Funcs 	*funcs.Funcs
	timeout int64
}
func NewServer() *Server {
	return &Server{Funcs:funcs.New(),timeout:DefaultServerTimeout}
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
	listener,err:=Listen(network,address)
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
		return s.ErrRPCEncode(msg.batch,errors.New("Version is not matched"))
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
	reply_bytes,ok:=s.CallServiceTimeOut(req.id,req.method,req.data,req.noResponse,msg.codecType)
	if ok{
		log.AllInfof("id-%d  req-%d CallService %s success",msg.id, req.id,req.method)
		if req.noResponse==true{
			return nil,true
		}
		res:=&Response{}
		res.id=req.id
		res.data=reply_bytes
		res_bytes,_:=res.Encode()
		if err!=nil{
			log.Warnln("id-%d  req-%d ResponseEncode error: %s ",msg.id, req.id, err)
			return s.ErrResponseEncode(req.id,err),req.noResponse
		}
		return res_bytes,req.noResponse
	}
	return reply_bytes,req.noResponse
}

func (s *Server)CallServiceTimeOut(id uint64,method string,args_bytes []byte, noResponse bool,funcsCodecType CodecType) ([]byte,bool){
	if s.timeout>0{
		ch := make(chan int)
		var (
			data []byte
			ok bool
		)
		go func() {
			data,ok=s.CallService(id,method,args_bytes, noResponse,funcsCodecType)
			ch<-1
		}()
		select {
		case <-ch:
			return data,ok
		case <-time.After(time.Millisecond * time.Duration(s.timeout)):
			return s.ErrResponseEncode(id,errors.New(fmt.Sprintf("method %s time out",method))),false
		}
	}
	return s.CallService(id,method,args_bytes, noResponse,funcsCodecType)
}
func (s *Server)CallService(id uint64,method string,args_bytes []byte, noResponse bool,funcsCodecType CodecType) ([]byte,bool) {
	if s.Funcs.GetFunc(method)==nil{
		log.AllInfof("CallService %s is not supposted",method)
		return s.ErrResponseEncode(id,errors.New(fmt.Sprintf("method %s is not supposted",method))),false
	}
	args := s.Funcs.GetFuncIn(method,0)
	err:=ArgsDecode(args_bytes,args,funcsCodecType)
	if err!=nil{
		return s.ErrResponseEncode(id,err),false
	}
	reply :=s.Funcs.GetFuncIn(method,1)
	if err := s.Funcs.Call(method, args, reply); err != nil {
		return s.ErrResponseEncode(id,err),false
	}
	if !noResponse{
		var reply_bytes []byte
		reply_bytes,err=ReplyEncode(reply,funcsCodecType)
		if err!=nil{
			return s.ErrResponseEncode(id,err),false
		}
		return reply_bytes,true
	}else {
		return nil,true
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
