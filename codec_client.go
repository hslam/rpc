package rpc

import (
	"hslam.com/mgit/Mort/rpc/log"
)
type ClientCodec struct{
	client_id int64
	req_id   uint64
	name string
	args interface{}
	funcsCodecType CodecType
	noRequest bool
	noResponse bool
	reply interface{}
	res	*Response
}
func(c *ClientCodec)Encode() ([]byte, error)  {
	var req *Request
	if c.noRequest==false{
		args_bytes,err:=ArgsEncode(c.args,c.funcsCodecType)
		if err!=nil{
			log.Errorln("ArgsEncode error: ", err)
			return nil,err
		}
		req=&Request{c.req_id,c.name,c.noRequest,c.noResponse,args_bytes,}
	}else {
		req=&Request{c.req_id,c.name,c.noRequest,c.noResponse,nil,}
	}
	req_bytes,err:=req.Encode()
	if err!=nil{
		log.Errorln("RequestEncode error: ", err)
		return nil,err
	}
	msg:=&Msg{}
	msg.id=c.client_id
	msg.data=req_bytes
	msg.batch=false
	msg.msgType=MsgType(MsgTypeReq)
	msg.codecType=c.funcsCodecType
	return msg.Encode()
}

func (c *ClientCodec)Decode(b []byte) error  {
	var req_id =c.req_id
	msg:=&Msg{}
	err:=msg.Decode(b)
	if msg.msgType==MsgType(MsgTypeRes){
		res:=&Response{}
		err=res.Decode(msg.data)
		if err!=nil{
			log.Errorln("ResponseDecode error: ", err)
			return err
		}
		c.res=res
		if req_id==c.res.id{
			if res.err!=nil{
				return res.err
			}
			return ReplyDecode(res.data,c.reply,msg.codecType)
		}else {
			return ErrReqId
		}

	}else if msg.msgType==MsgType(MsgTypeHea){
		return nil
	}
	return nil
}