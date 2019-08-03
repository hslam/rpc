package rpc

import (
	"github.com/golang/protobuf/proto"
	"errors"
	"hslam.com/mgit/Mort/rpc/pb"
	"hslam.com/mgit/Mort/rpc/log"
)

type Request struct {
	id                  	uint64
	method					string
	data					[]byte
	noResponse				bool
}
func(r *Request)Encode() ([]byte, error)  {
	switch rpc_codec {
	case RPC_CODEC_ME:
		var msg = Msg{}
		return msg.Serialize(Version,r.method,r.data),nil
	case RPC_CODEC_PROTOBUF:
		req:=pb.Request{Id:r.id,Method:r.method,Data:r.data,NoResponse:r.noResponse}
		if data,err:=proto.Marshal(&req);err!=nil{
			log.Errorln("RequestEncode proto.Unmarshal error: ", err)
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
			log.Errorln("RequestDecode proto.Unmarshal error: ", err)
			return err
		}
		r.id=rpc_req_decode.Id
		r.method=rpc_req_decode.Method
		r.data=rpc_req_decode.Data
		r.noResponse=rpc_req_decode.NoResponse
	default:
		return errors.New("this mrpc_serialize is not supported")
	}
	return nil
}
