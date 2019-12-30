package rpc

import (
	"github.com/hslam/rpc/pb"
	"errors"
)
type ClientCodec struct{
	client_id uint64
	req_id   uint64
	name string
	args interface{}
	funcsCodecType CodecType
	noRequest bool
	noResponse bool
	reply interface{}
	res	*Response
	batch				bool
	batchAsync			bool
	compressType		CompressType
	compressLevel 		CompressLevel
	msg					*Msg
	batchCodec			*BatchCodec
	requests			[]*BatchRequest
	responses			[]*Response
}
func(c *ClientCodec)Encode() ([]byte, error)  {
	var err error
	c.msg=&Msg{}
	c.msg.version=Version
	c.msg.id=c.client_id
	c.msg.msgType=MsgTypeReq
	c.msg.batch=c.batch
	c.msg.codecType=c.funcsCodecType
	c.msg.compressType=c.compressType
	c.msg.compressLevel=c.compressLevel
	if c.batch==false{
		var req *Request
		if c.noRequest==false{
			args_bytes,err:=ArgsEncode(c.args,c.funcsCodecType)
			if err!=nil{
				Errorln("ArgsEncode error: ", err)
				return nil,err
			}
			req=&Request{c.req_id,c.name,c.noRequest,c.noResponse,args_bytes,}
		}else {
			req=&Request{c.req_id,c.name,c.noRequest,c.noResponse,nil,}
		}
		c.msg.data,err=req.Encode()
		if err!=nil{
			Errorln("RequestEncode error: ", err)
			return nil,err
		}
	}else {
		req_bytes_s:=make([][]byte,len(c.requests))
		c.responses=make([]*Response,len(c.requests))
		for i,v :=range c.requests{
			req:=&Request{v.id,v.name,v.noRequest,v.noResponse,v.args_bytes}
			req_bytes,_:=req.Encode()
			req_bytes_s[i]=req_bytes
			c.responses[i]=&Response{}
		}
		batchCodec:=&BatchCodec{async:c.batchAsync,data:req_bytes_s}
		c.batchCodec=batchCodec
		c.msg.data,_=batchCodec.Encode()
	}
	return c.msg.Encode()
}

func (c *ClientCodec)Decode(b []byte) error  {
	var req_id =c.req_id
	if c.msg==nil{
		c.msg=&Msg{}
	}else {
		c.msg.Reset()
	}
	err:=c.msg.Decode(b)
	if c.msg.msgType!=MsgType(pb.MsgType_res){
		return errors.New("not be MsgType_res")
	}
	if c.msg.batch==false {
		res:=&Response{}
		err=res.Decode(c.msg.data)
		if err!=nil{
			Errorln("ResponseDecode error: ", err)
			return err
		}
		c.res=res
		if req_id==c.res.id{
			if res.err!=nil{
				return res.err
			}
			return ReplyDecode(res.data,c.reply,c.msg.codecType)
		}else {
			return ErrReqId
		}
	}else {
		c.batchCodec.Reset()
		err=c.batchCodec.Decode(c.msg.data)
		if err!=nil{
			return err
		}
		if len(c.batchCodec.data)==len(c.requests){
			for i,res:=range c.responses{
				res.Decode(c.batchCodec.data[i])
			}
		}
		return nil
	}
}