package rpc

import (
	"github.com/golang/protobuf/proto"
	"errors"
	"hslam.com/git/x/rpc/pb"
	"fmt"
)
type Response struct {
	id          uint64
	data		[]byte
	errMsg		string
	err			error
}
func(r *Response)Encode() ([]byte, error)  {
	switch rpc_codec {
	case RPC_CODEC_ME:
		return r.data,nil
	case RPC_CODEC_PROTOBUF:
		if r.err!=nil{
			r.errMsg=fmt.Sprint(r.err)
		}
		rpc_res:= pb.Response{Id:r.id,Data:r.data,ErrMsg:r.errMsg}
		if rpc_res_bytes, err := proto.Marshal(&rpc_res); err != nil {
			Errorln("ResponseEncode proto.Marshal error: ", err)
			return nil,err
		}else {
			return rpc_res_bytes,nil
		}
	default:
		return nil,errors.New("this mrpc_serialize is not supported")
	}
}

func(r *Response)Decode(b []byte) (error)  {
	switch rpc_codec {
	case RPC_CODEC_ME:
		r.data=b
		r.err=nil
		return nil
	case RPC_CODEC_PROTOBUF:
		var rpc_res_decode pb.Response
		if err := proto.Unmarshal(b, &rpc_res_decode); err != nil {
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
	default:
		return errors.New("mrpc_serialize is not supported")
	}
}
