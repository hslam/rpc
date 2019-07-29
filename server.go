package rpc

import (
	"hslam.com/mgit/Mort/funcs"
	"hslam.com/mgit/Mort/rpc/pb"
	"errors"
	"fmt"
)

var DefaultServer = NewServer()

type Server struct {
	network string
	listener Listener
	Funcs 	*funcs.Funcs
}
func NewServer() *Server {
	return &Server{Funcs:funcs.New()}
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

func ServeRPC(b []byte) (bool,[]byte,error) {
	return DefaultServer.ServeRPC(b)
}
func (s *Server)ServeRPC(b []byte) (bool,[]byte,error) {
	defer func() {
		if err := recover(); err != nil {
			Panicf("CallService error: %s", err)
		}
	}()

	msg:=&Msg{}
	err:=msg.Decode(b)
	if err!=nil{
		Warnf("CallService MrpcDecode error: %s", err)
		return s.ErrRPCEncode(msg.batch,err)
	}
	if msg.version!=Version{
		return s.ErrRPCEncode(msg.batch,errors.New("Version is not matched"))
	}
	var noResponse =false
	var responseBytes []byte
	if msg.batch{
		responseBytes,noResponse=s.BatchHandleRPC(msg)
	}else{
		responseBytes,noResponse=s.HandleRPC(msg)
	}
	if noResponse==true{
		return true,nil,nil
	}
	msg.data=responseBytes
	msg.msgType=MsgType(pb.MsgType_res)
	msg_bytes,err:=msg.Encode()
	return true,msg_bytes,err
}

func (s *Server)HandleRPC(msg *Msg) ([]byte,bool) {
	return s.Handler(msg)
}

func (s *Server)BatchHandleRPC(msg *Msg) ([]byte,bool) {
	var noResponse =false
	var responseBytes []byte
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
	return responseBytes,noResponse
}
func (s *Server)Handler(msg *Msg) ([]byte,bool) {
	req:=&Request{}
	err:=req.Decode(msg.data)
	if err!=nil{
		Warnln("id-%d  req-%d RequestDecode error: %s ",msg.id, req.id,err)
		return s.ErrResponseEncode(err),false
	}
	reply_bytes,ok:=s.CallService(req.method,req.data,req.noResponse,msg.codecType)
	if ok{
		AllInfof("id-%d  req-%d CallService %s success",msg.id, req.id,req.method)
		if req.noResponse==true{
			return nil,true
		}
		res:=&Response{}
		res.data=reply_bytes
		res_bytes,_:=res.Encode()
		if err!=nil{
			Warnln("id-%d  req-%d ResponseEncode error: %s ",msg.id, req.id, err)
			return s.ErrResponseEncode(err),req.noResponse
		}
		return res_bytes,req.noResponse
	}
	return reply_bytes,req.noResponse
}

func (s *Server)CallService(method string,args_bytes []byte, noResponse bool,funcsCodecType CodecType) ([]byte,bool) {
	if s.Funcs.GetFunc(method)==nil{
		AllInfof("%s CallService %s is not supposted",method)
		return s.ErrResponseEncode(errors.New(fmt.Sprintf("method %s is not supposted",method))),false
	}
	args := s.Funcs.GetFuncIn(method,0)
	err:=ArgsDecode(args_bytes,args,funcsCodecType)
	if err!=nil{
		return s.ErrResponseEncode(err),false
	}
	reply :=s.Funcs.GetFuncIn(method,1)
	if err := s.Funcs.Call(method, args, reply); err != nil {
		return s.ErrResponseEncode(err),false
	}
	if !noResponse{
		var reply_bytes []byte
		reply_bytes,err=ReplyEncode(reply,funcsCodecType)
		if err!=nil{
			return s.ErrResponseEncode(err),false
		}
		return reply_bytes,true
	}else {
		return nil,true
	}
}

func (s *Server)ErrRPCEncode(batch bool,errMsg error)(bool,[]byte,error)  {
	res_bytes:=s.ErrResponseEncode(errMsg)
	msg:=&Msg{}
	msg.data=res_bytes
	msg.batch=batch
	msg.msgType=MsgType(pb.MsgType_res)
	msg.codecType=FUNCS_CODEC_INVALID
	msg_bytes,err:=msg.Encode()
	return false,msg_bytes,err
}


func (s *Server)ErrResponseEncode(err error)[]byte  {
	res:=&Response{}
	res.err=err
	res_bytes,_:=res.Encode()
	return res_bytes
}
