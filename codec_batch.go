package rpc

import (
	"github.com/golang/protobuf/proto"
	"hslam.com/git/x/rpc/pb"
)
type BatchCodec struct{
	async bool
	data [][]byte
}

func(c *BatchCodec)Encode()([]byte,error)  {
	batch:=pb.Batch{Async:c.async,Data:c.data}
	batch_bytes,err:= proto.Marshal(&batch)
	if err != nil {
		Errorln("BatchEncode proto.Marshal error: ", err)
		return nil, err
	}
	return batch_bytes,nil
}
func(c *BatchCodec)Decode(b []byte)(error)  {
	var batch pb.Batch
	if err := proto.Unmarshal(b, &batch); err != nil {
		Errorln("BatchDecode proto.Unmarshal error: ", err)
		return  err
	}
	c.async=batch.Async
	c.data=batch.Data
	return nil
}

