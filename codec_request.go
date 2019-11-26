package rpc

import (
	"github.com/golang/protobuf/proto"
	"errors"
	"hslam.com/git/x/rpc/pb"
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
	case RPC_CODEC_ME:
		var msg = Msg{}
		return msg.Serialize(Version,r.method,r.data),nil
	case RPC_CODEC_PROTOBUF:
		req:=pb.Request{
			Id:r.id,
			Method:r.method,
			NoRequest:r.noRequest,
			NoResponse:r.noResponse,
			Data:r.data,
		}
		if data,err:=proto.Marshal(&req);err!=nil{
			Errorln("RequestEncode proto.Unmarshal error: ", err)
			return nil,err
		}else {
			return data,nil
		}
	}
	return nil,errors.New("this mrpc_serialize is not supported")
}

func(r *Request)Decode(b []byte) (error)  {
	r.noResponse=false
	switch rpc_codec {
	case RPC_CODEC_ME:
		var msg = Msg{}
		_,r.method,r.data=msg.Deserialize(b)
	case RPC_CODEC_PROTOBUF:
		var rpc_req_decode pb.Request
		if err := proto.Unmarshal(b, &rpc_req_decode); err != nil {
			Errorln("RequestDecode proto.Unmarshal error: ", err)
			return err
		}
		r.id=rpc_req_decode.Id
		r.method=rpc_req_decode.Method
		r.noRequest=rpc_req_decode.NoRequest
		r.noResponse=rpc_req_decode.NoResponse
		r.data=rpc_req_decode.Data
	default:
		return errors.New("this mrpc_serialize is not supported")
	}
	return nil
}
