package rpc

import (
	"errors"
	"github.com/hslam/code"
	"github.com/hslam/rpc/pb"
)

type batchCodec struct {
	async bool
	data  [][]byte
}

func (c *batchCodec) Encode() ([]byte, error) {
	switch rpcCodec {
	case RPCCodecCode:
		return c.Marshal(nil)
	case RPCCodecProtobuf:
		batch := pb.Batch{Async: c.async, Data: c.data}
		batchBytes, err := batch.Marshal()
		if err != nil {
			logger.Errorln("BatchEncode proto.Marshal error: ", err)
			return nil, err
		}
		return batchBytes, nil
	default:
		return nil, errors.New("this rpc_serialize is not supported")
	}
}
func (c *batchCodec) Decode(b []byte) error {
	switch rpcCodec {
	case RPCCodecCode:
		return c.Unmarshal(b)
	case RPCCodecProtobuf:
		var batch = &pb.Batch{}
		if err := batch.Unmarshal(b); err != nil {
			logger.Errorln("BatchDecode proto.Unmarshal error: ", err)
			return err
		}
		c.async = batch.Async
		c.data = batch.Data
		return nil
	default:
		return errors.New("rpc_serialize is not supported")
	}
}

func (c *batchCodec) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size++
	size += code.SizeofBytesSlice(c.data)
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	n = code.EncodeBool(buf[offset:], c.async)
	offset += n
	n = code.EncodeBytesSlice(buf[offset:], c.data)
	offset += n
	return buf, nil
}
func (c *batchCodec) Unmarshal(b []byte) error {
	var offset uint64
	var n uint64
	n = code.DecodeBool(b[offset:], &c.async)
	offset += n
	n = code.DecodeBytesSlice(b[offset:], &c.data)
	offset += n
	return nil
}
func (c *batchCodec) Reset() {
	for i, v := range c.data {
		c.data[i] = v[len(v):]
	}
}
