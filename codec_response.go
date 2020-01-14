package rpc

import (
	"errors"
	"fmt"
	"github.com/hslam/code"
	"github.com/hslam/rpc/pb"
)

type Response struct {
	id     uint64
	data   []byte
	errMsg string
	err    error
}

func (r *Response) Encode() ([]byte, error) {
	switch rpc_codec {
	case RPC_CODEC_CODE:
		return r.Marshal(nil)
	case RPC_CODEC_PROTOBUF:
		if r.err != nil {
			r.errMsg = fmt.Sprint(r.err)
		}
		rpc_res := pb.Response{Id: r.id, Data: r.data, ErrMsg: r.errMsg}
		if rpc_res_bytes, err := rpc_res.Marshal(); err != nil {
			Errorln("ResponseEncode proto.Marshal error: ", err)
			return nil, err
		} else {
			return rpc_res_bytes, nil
		}
	default:
		return nil, errors.New("this rpc_serialize is not supported")
	}
}

func (r *Response) Decode(b []byte) error {
	switch rpc_codec {
	case RPC_CODEC_CODE:
		return r.Unmarshal(b)
	case RPC_CODEC_PROTOBUF:
		var rpc_res_decode = &pb.Response{}
		if err := rpc_res_decode.Unmarshal(b); err != nil {
			Errorln("ResponseDecode proto.Unmarshal error: ", err)
			return err
		}
		r.id = rpc_res_decode.Id
		if rpc_res_decode.ErrMsg == "" {
			r.data = rpc_res_decode.Data
			return nil
		} else {
			r.err = errors.New(rpc_res_decode.ErrMsg)
			return nil
		}
	default:
		return errors.New("rpc_serialize is not supported")
	}
}

func (r *Response) Marshal(buf []byte) ([]byte, error) {
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
func (r *Response) Unmarshal(b []byte) error {
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
func (r *Response) Reset() {
}
