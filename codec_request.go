package rpc

import (
	"errors"
	"github.com/hslam/code"
	"github.com/hslam/rpc/pb"
)

type request struct {
	id         uint64
	method     string
	noRequest  bool
	noResponse bool
	data       []byte
}

func (r *request) Encode() ([]byte, error) {
	switch rpcCodec {
	case RPCCodecCode:
		return r.Marshal(nil)
	case RPCCodecProtobuf:
		req := pb.Request{
			Id:         r.id,
			Method:     r.method,
			NoRequest:  r.noRequest,
			NoResponse: r.noResponse,
			Data:       r.data,
		}
		var data []byte
		var err error
		if data, err = req.Marshal(); err != nil {
			logger.Errorln("RequestEncode proto.Unmarshal error: ", err)
			return nil, err
		}
		return data, nil
	default:
		return nil, errors.New("this rpc_serialize is not supported")
	}
}

func (r *request) Decode(b []byte) error {
	r.noResponse = false
	switch rpcCodec {
	case RPCCodecCode:
		return r.Unmarshal(b)
	case RPCCodecProtobuf:
		var rpcReqDecode = &pb.Request{}
		if err := rpcReqDecode.Unmarshal(b); err != nil {
			logger.Errorln("RequestDecode proto.Unmarshal error: ", err)
			return err
		}
		r.id = rpcReqDecode.Id
		r.method = rpcReqDecode.Method
		r.noRequest = rpcReqDecode.NoRequest
		r.noResponse = rpcReqDecode.NoResponse
		r.data = rpcReqDecode.Data
	default:
		return errors.New("this rpc_serialize is not supported")
	}
	return nil
}

func (r *request) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 8
	size += code.SizeofString(r.method)
	size++
	size++
	size += code.SizeofBytes(r.data)
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	n = code.EncodeUint64(buf[offset:], r.id)
	offset += n
	n = code.EncodeString(buf[offset:], r.method)
	offset += n
	n = code.EncodeBool(buf[offset:], r.noRequest)
	offset += n
	n = code.EncodeBool(buf[offset:], r.noResponse)
	offset += n
	n = code.EncodeBytes(buf[offset:], r.data)
	offset += n
	return buf, nil
}
func (r *request) Unmarshal(b []byte) error {
	var offset uint64
	var n uint64
	n = code.DecodeUint64(b[offset:], &r.id)
	offset += n
	n = code.DecodeString(b[offset:], &r.method)
	offset += n
	n = code.DecodeBool(b[offset:], &r.noRequest)
	offset += n
	n = code.DecodeBool(b[offset:], &r.noResponse)
	offset += n
	n = code.DecodeBytes(b[offset:], &r.data)
	offset += n
	return nil
}
func (r *request) Reset() {
}
