package rpc

import (
	"github.com/golang/protobuf/proto"
	"hslam.com/mgit/Mort/rpc/pb"
	"hslam.com/mgit/Mort/rpc/log"
)
type BatchCodec struct{
	data [][]byte
}

func(c *BatchCodec)Encode()([]byte,error)  {
	batch:=pb.Batch{Data:c.data}
	batch_bytes,err:= proto.Marshal(&batch)
	if err != nil {
		log.Errorln("BatchEncode proto.Marshal error: ", err)
		return nil, err
	}
	return batch_bytes,nil
}
func(c *BatchCodec)Decode(b []byte)(error)  {
	var batch pb.Batch
	if err := proto.Unmarshal(b, &batch); err != nil {
		log.Errorln("BatchDecode proto.Unmarshal error: ", err)
		return  err
	}
	c.data=batch.Data
	return nil
}

