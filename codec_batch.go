package rpc

import (
	"hslam.com/git/x/rpc/pb"
	"hslam.com/git/x/rpc/gen"
	"errors"
)
type BatchCodec struct{
	async bool
	data [][]byte
}

func(c *BatchCodec)Encode()([]byte,error)  {
	switch rpc_codec {
	case RPC_CODEC_ME:
		return c.Marshal(nil)
	case RPC_CODEC_PROTOBUF:
		batch:=pb.Batch{Async:c.async,Data:c.data}
		batch_bytes,err:= batch.Marshal()
		if err != nil {
			Errorln("BatchEncode proto.Marshal error: ", err)
			return nil, err
		}
		return batch_bytes,nil
	case RPC_CODEC_GENCODE:
		batch:= gen.Batch{Async:c.async,Data:c.data}
		batch_bytes,err:= batch.Marshal(nil)
		if err != nil {
			Errorln("BatchEncode gencode.Marshal error: ", err)
			return nil, err
		}
		return batch_bytes,nil
	default:
		return nil,errors.New("this rpc_serialize is not supported")
	}
}
func(c *BatchCodec)Decode(b []byte)(error)  {
	switch rpc_codec {
	case RPC_CODEC_ME:
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
	case RPC_CODEC_GENCODE:
		var batch =&gen.Batch{}
		if _,err := batch.Unmarshal(b); err != nil {
			Errorln("BatchDecode gencode.Unmarshal error: ", err)
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
	return nil,nil
}
func(c *BatchCodec)Unmarshal(b []byte)(error)  {
	return nil
}
