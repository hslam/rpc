package rpc

import (
	"errors"
	"hslam.com/git/x/rpc/pb"
	"hslam.com/git/x/rpc/gen"
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
	case RPC_CODEC_RAW:
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
	case RPC_CODEC_GENCODE:
		req:= gen.Request{
			Id:r.id,
			Method:r.method,
			NoRequest:r.noRequest,
			NoResponse:r.noResponse,
			Data:r.data,
		}
		if data,err:=req.Marshal(nil);err!=nil{
			Errorln("RequestEncode gencode.Unmarshal error: ", err)
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
	case RPC_CODEC_RAW:
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
	case RPC_CODEC_GENCODE:
		var rpc_req_decode =&gen.Request{}
		if _,err := rpc_req_decode.Unmarshal(b); err != nil {
			Errorln("RequestDecode gencode.Unmarshal error: ", err)
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
	return nil,nil
}
func(r *Request)Unmarshal(b []byte)(error)  {
	return nil
}
func(r *Request)Reset()()  {
}
