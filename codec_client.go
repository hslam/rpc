package rpc

import (
	"hslam.com/mgit/Mort/rpc/pb"
)
type ClientCodec struct{
	client_id int64
	req_id   uint64
	name string
	args interface{}
	funcsCodecType CodecType
	noResponse bool
	reply interface{}
}
func(c *ClientCodec)Encode() ([]byte, error)  {
	args_bytes,err:=ArgsEncode(c.args,c.funcsCodecType)
	if err!=nil{
		Errorln("ArgsEncode error: ", err)
		return nil,err
	}
	req:=&Request{c.req_id,c.name,args_bytes,c.noResponse}
	req_bytes,err:=req.Encode()
	if err!=nil{
		Errorln("RequestEncode error: ", err)
		return nil,err
	}
	msg:=&Msg{}
	msg.id=c.client_id
	msg.data=req_bytes
	msg.batch=false
	msg.msgType=MsgType(pb.MsgType_req)
	msg.codecType=c.funcsCodecType
	return msg.Encode()
}

func (c *ClientCodec)Decode(b []byte) error  {
	msg:=&Msg{}
	err:=msg.Decode(b)
	res:=&Response{}
	err=res.Decode(msg.data)
	if err!=nil{
		Errorln("ResponseDecode error: ", err)
		return err
	}
	if res.err!=nil{
		return res.err
	}
	return ReplyDecode(res.data,c.reply,msg.codecType)
}