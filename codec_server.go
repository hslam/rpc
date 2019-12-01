package rpc

import (
	"fmt"
)

type ServerCodec struct{
	msg					*Msg
	batchCodec			*BatchCodec
	request				*Request
	response			*Response
	requests			[]*Request
	responses			[]*Response
}
func(c *ServerCodec)Encode() ([]byte, error)  {
	return c.msg.Encode()
}

func (c *ServerCodec)Decode(b []byte) error  {
	msg:=&Msg{}
	err:=msg.Decode(b)
	c.msg=msg
	if err!=nil{
		Warnf("ServerCodec.Decode msg error: %s", err)
		return fmt.Errorf("ServerCodec.Decode msg error: %s",err)
	}
	if msg.version!=Version{
		Warnf("%d %d Version is not matched",Version,msg.version)
		return fmt.Errorf("%d %d Version is not matched",Version,msg.version)
	}
	if msg.msgType==MsgTypeHea{
		return nil
	}
	if msg.batch {
		batchCodec:=&BatchCodec{}
		batchCodec.Decode(msg.data)
		c.batchCodec=batchCodec
		length:=len(batchCodec.data)
		c.requests=make([]*Request,length)
		c.responses=make([]*Response,length)
		for i,v:=range batchCodec.data{
			req:=&Request{}
			err:=req.Decode(v)
			if err!=nil{
				c.requests[i]=nil
			}
			c.requests[i]=req
			c.responses[i]=&Response{}
		}
		batchCodec.data=make([][]byte,length)
	}else {
		req:=&Request{}
		err:=req.Decode(msg.data)
		if err!=nil{
			Warnf("ServerCodec.Decode id:%d req:%d error:%s ",msg.id, req.id,err)
			return fmt.Errorf("ServerCodec.Decode id:%d req:%d error:%s ",msg.id, req.id,err)
		}
		c.request=req
		c.response=&Response{}
	}
	return nil
}