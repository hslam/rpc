package rpc

import (
	"errors"
	"github.com/hslam/rpc/pb"
	"fmt"
	"github.com/hslam/rpc/gen"
)
type Response struct {
	id          uint64
	data		[]byte
	errMsg		string
	err			error
}
func(r *Response)Encode() ([]byte, error)  {
	switch rpc_codec {
	case RPC_CODEC_RAW:
		return r.Marshal(nil)
	case RPC_CODEC_PROTOBUF:
		if r.err!=nil{
			r.errMsg=fmt.Sprint(r.err)
		}
		rpc_res:= pb.Response{Id:r.id,Data:r.data,ErrMsg:r.errMsg}
		if rpc_res_bytes, err := rpc_res.Marshal(); err != nil {
			Errorln("ResponseEncode proto.Marshal error: ", err)
			return nil,err
		}else {
			return rpc_res_bytes,nil
		}
	case RPC_CODEC_GENCODE:
		if r.err!=nil{
			r.errMsg=fmt.Sprint(r.err)
		}
		rpc_res:= gen.Response{Id:r.id,Data:r.data,ErrMsg:r.errMsg}
		if rpc_res_bytes, err := rpc_res.Marshal(nil); err != nil {
			Errorln("ResponseEncode gencode.Marshal error: ", err)
			return nil,err
		}else {
			return rpc_res_bytes,nil
		}
	default:
		return nil,errors.New("this rpc_serialize is not supported")
	}
}

func(r *Response)Decode(b []byte) (error)  {
	switch rpc_codec {
	case RPC_CODEC_RAW:
		return r.Unmarshal(b)
	case RPC_CODEC_PROTOBUF:
		var rpc_res_decode =&pb.Response{}
		if err := rpc_res_decode.Unmarshal(b); err != nil {
			Errorln("ResponseDecode proto.Unmarshal error: ", err)
			return err
		}
		r.id=rpc_res_decode.Id
		if rpc_res_decode.ErrMsg==""{
			r.data=rpc_res_decode.Data
			return nil
		}else {
			r.err=errors.New(rpc_res_decode.ErrMsg)
			return nil
		}
	case RPC_CODEC_GENCODE:
		var rpc_res_decode =&gen.Response{}
		if _,err := rpc_res_decode.Unmarshal(b); err != nil {
			Errorln("ResponseDecode gencode.Unmarshal error: ", err)
			return err
		}
		r.id=rpc_res_decode.Id
		if rpc_res_decode.ErrMsg==""{
			r.data=rpc_res_decode.Data
			return nil
		}else {
			r.err=errors.New(rpc_res_decode.ErrMsg)
			return nil
		}
	default:
		return errors.New("rpc_serialize is not supported")
	}
}

func(r *Response)Marshal(buf []byte)([]byte,error)  {
	return nil,nil
}
func(r *Response)Unmarshal(b []byte)(error)  {
	return nil
}
func(r *Response)Reset()()  {
}