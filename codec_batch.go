package rpc

import (
	"github.com/hslam/rpc/pb"
	"github.com/hslam/code"
	"errors"
)
type BatchCodec struct{
	async bool
	data [][]byte
}

func(c *BatchCodec)Encode()([]byte,error)  {
	switch rpc_codec {
	case RPC_CODEC_CODE:
		return c.Marshal(nil)
	case RPC_CODEC_PROTOBUF:
		batch:=pb.Batch{Async:c.async,Data:c.data}
		batch_bytes,err:= batch.Marshal()
		if err != nil {
			Errorln("BatchEncode proto.Marshal error: ", err)
			return nil, err
		}
		return batch_bytes,nil
	default:
		return nil,errors.New("this rpc_serialize is not supported")
	}
}
func(c *BatchCodec)Decode(b []byte)(error)  {
	switch rpc_codec {
	case RPC_CODEC_CODE:
		return c.Unmarshal(b)
	case RPC_CODEC_PROTOBUF:
		var batch =&pb.Batch{}
		if err := batch.Unmarshal(b); err != nil {
			Errorln("BatchDecode proto.Unmarshal error: ", err)
			return  err
		}
		c.async=batch.Async
		c.data=batch.Data
		return nil
	default:
		return errors.New("rpc_serialize is not supported")
	}
}

func(c *BatchCodec)Marshal(buf []byte)([]byte,error)  {
	var size uint64
	size+=1
	size+=code.SizeofSliceBytes(c.data)
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n  uint64
	n = code.EncodeBool(buf[offset:],c.async)
	offset+=n
	n = code.EncodeSliceBytes(buf[offset:],c.data)
	offset+=n
	return buf,nil
}
func(c *BatchCodec)Unmarshal(b []byte)(error)  {
	var offset uint64
	var n uint64
	n=code.DecodeBool(b[offset:],&c.async)
	offset+=n
	n=code.DecodeSliceBytes(b[offset:],&c.data)
	offset+=n
	return nil
}
func(c *BatchCodec)Reset(){
	for i,v:=range c.data{
		c.data[i]=v[len(v):]
	}
}

