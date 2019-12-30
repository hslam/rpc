package rpc

import (
	"errors"
	"github.com/hslam/rpc/pb"
	"github.com/hslam/code"
)

type Request struct {
	id                  	uint64
	method					string
	noRequest				bool
	noResponse				bool
	data					[]byte
}
func(r *Request)Encode() ([]byte, error)  {
	switch rpc_codec {
	case RPC_CODEC_CODE:
		return r.Marshal(nil)
	case RPC_CODEC_PROTOBUF:
		req:=pb.Request{
			Id:r.id,
			Method:r.method,
			NoRequest:r.noRequest,
			NoResponse:r.noResponse,
			Data:r.data,
		}
		if data,err:=req.Marshal();err!=nil{
			Errorln("RequestEncode proto.Unmarshal error: ", err)
			return nil,err
		}else {
			return data,nil
		}
	default:
		return nil,errors.New("this rpc_serialize is not supported")
	}
}

func(r *Request)Decode(b []byte) (error)  {
	r.noResponse=false
	switch rpc_codec {
	case RPC_CODEC_CODE:
		return r.Unmarshal(b)
	case RPC_CODEC_PROTOBUF:
		var rpc_req_decode =&pb.Request{}
		if err := rpc_req_decode.Unmarshal(b); err != nil {
			Errorln("RequestDecode proto.Unmarshal error: ", err)
			return err
		}
		r.id=rpc_req_decode.Id
		r.method=rpc_req_decode.Method
		r.noRequest=rpc_req_decode.NoRequest
		r.noResponse=rpc_req_decode.NoResponse
		r.data=rpc_req_decode.Data
	default:
		return errors.New("this rpc_serialize is not supported")
	}
	return nil
}

func(r *Request)Marshal(buf []byte)([]byte,error)  {
	var size uint64
	size+=8
	size+=code.SizeofString(r.method)
	size+=1
	size+=1
	size+=code.SizeofBytes(r.data)
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n  uint64
	n = code.EncodeUint64(buf[offset:],r.id)
	offset+=n
	n = code.EncodeString(buf[offset:],r.method)
	offset+=n
	n = code.EncodeBool(buf[offset:],r.noRequest)
	offset+=n
	n = code.EncodeBool(buf[offset:],r.noResponse)
	offset+=n
	n = code.EncodeBytes(buf[offset:],r.data)
	offset+=n
	return buf,nil
}
func(r *Request)Unmarshal(b []byte)(error)  {
	var offset uint64
	var n uint64
	n=code.DecodeUint64(b[offset:],&r.id)
	offset+=n
	n=code.DecodeString(b[offset:],&r.method)
	offset+=n
	n=code.DecodeBool(b[offset:],&r.noRequest)
	offset+=n
	n=code.DecodeBool(b[offset:],&r.noResponse)
	offset+=n
	n=code.DecodeBytes(b[offset:],&r.data)
	offset+=n
	return nil
}
func(r *Request)Reset()()  {
}
