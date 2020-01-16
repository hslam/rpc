package rpc

import (
	"errors"
	"fmt"
	"github.com/hslam/code"
	"github.com/hslam/rpc/pb"
)

type response struct {
	id     uint64
	data   []byte
	errMsg string
	err    error
}

func (r *response) Encode() ([]byte, error) {
	switch rpcCodec {
	case RPCCodecCode:
		return r.Marshal(nil)
	case RPCCodecProtobuf:
		if r.err != nil {
			r.errMsg = fmt.Sprint(r.err)
		}
		rpcRes := pb.Response{Id: r.id, Data: r.data, ErrMsg: r.errMsg}
		var rpcResBytes []byte
		var err error
		if rpcResBytes, err = rpcRes.Marshal(); err != nil {
			logger.Errorln("ResponseEncode proto.Marshal error: ", err)
			return nil, err
		}
		return rpcResBytes, nil
	default:
		return nil, errors.New("this rpc_serialize is not supported")
	}
}

func (r *response) Decode(b []byte) error {
	switch rpcCodec {
	case RPCCodecCode:
		return r.Unmarshal(b)
	case RPCCodecProtobuf:
		var rpcResDecode = &pb.Response{}
		if err := rpcResDecode.Unmarshal(b); err != nil {
			logger.Errorln("ResponseDecode proto.Unmarshal error: ", err)
			return err
		}
		r.id = rpcResDecode.Id
		if rpcResDecode.ErrMsg == "" {
			r.data = rpcResDecode.Data
			return nil
		}
		r.err = errors.New(rpcResDecode.ErrMsg)
		return nil
	default:
		return errors.New("rpc_serialize is not supported")
	}
}

func (r *response) Marshal(buf []byte) ([]byte, error) {
	if r.err != nil {
		r.errMsg = fmt.Sprint(r.err)
	}
	var size uint64
	size += 8
	size += code.SizeofBytes(r.data)
	size += code.SizeofString(r.errMsg)
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	n = code.EncodeUint64(buf[offset:], r.id)
	offset += n
	n = code.EncodeBytes(buf[offset:], r.data)
	offset += n

	n = code.EncodeString(buf[offset:], r.errMsg)
	offset += n
	return buf, nil
}
func (r *response) Unmarshal(b []byte) error {
	var offset uint64
	var n uint64
	n = code.DecodeUint64(b[offset:], &r.id)
	offset += n
	n = code.DecodeBytes(b[offset:], &r.data)
	offset += n
	n = code.DecodeString(b[offset:], &r.errMsg)
	offset += n
	if r.errMsg != "" {
		r.err = errors.New(r.errMsg)
	}
	return nil
}
func (r *response) Reset() {
}
