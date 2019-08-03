package rpc

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"errors"
	"hslam.com/mgit/Mort/rpc/pb"
	"hslam.com/mgit/Mort/rpc/log"
)

const (
	HeaderBits		=	16
	VersionBits		=	5
	VersionLeftShift	=	MethodBits

	MethodBits		=	HeaderBits-VersionBits
	MethodMask		=	-1 ^ (-1 << MethodBits)
)


var (
	//msg Msg
	rpc_codec = RPC_CODEC_PROTOBUF
)

func init() {
	//msg = Msg{}
}

//type Msg struct {}

func (m Msg)Serialize(version int32,method string,data []byte)[]byte  {
	method_len:=uint16(len(method))
	headeruint16:=uint16(version)<<VersionLeftShift+method_len
	header:=uint16ToBytes(headeruint16)
	var Bytes=append(append(header, []byte(method)...), data...)
	return Bytes
}

func (m Msg)Deserialize(msg []byte)(int,string,[]byte)  {
	var buf =bytes.NewBuffer(msg)
	header:=make([]byte , 2)
	buf.Read(header)
	headeruint16:=bytesToUint16(header)
	v:=headeruint16>>VersionLeftShift
	method_len:=headeruint16 & MethodMask
	method_bytes_len:=int(method_len)
	method_bytes:=make([]byte,method_bytes_len)
	buf.Read(method_bytes)
	method:=string(method_bytes)
	data:=make([]byte,len(msg)-method_bytes_len-2)
	buf.Read(data)
	return int(v),method,data
}

type Msg struct {
	version				int32
	id                  int64
	msgType				MsgType
	batch				bool
	codecType			CodecType
	data				[]byte
}
func (m *Msg)Encode() ([]byte, error) {
	switch rpc_codec {
	case RPC_CODEC_ME:
		return m.data,nil
	case RPC_CODEC_PROTOBUF:
		var msg pb.Msg
		if m.msgType==MsgType(pb.MsgType_req)||m.msgType==MsgType(pb.MsgType_res){
			msg=pb.Msg{Version:Version,Id:m.id,MsgType:pb.MsgType(m.msgType),Batch:m.batch,Data:m.data,CodecType:int32(m.codecType)}
		}else if m.msgType==MsgType(pb.MsgType_hea) {
			msg=pb.Msg{Version:Version,Id:m.id,MsgType:pb.MsgType(m.msgType)}
		}
		if data,err:=proto.Marshal(&msg);err!=nil{
			log.Errorln("MsgEncode proto.Unmarshal error: ", err)
			return nil,err
		}else {
			return data,nil
		}
	}
	return nil,errors.New("this mrpc_serialize is not supported")
}
func (m *Msg)Decode(b []byte)(error) {
	m.version=-1
	m.id=-1
	m.data=nil
	m.batch=false
	m.codecType=FUNCS_CODEC_INVALID
	switch rpc_codec {
	case RPC_CODEC_ME:
		m.version=Version
		m.data=b
		return nil
	case RPC_CODEC_PROTOBUF:
		var msg pb.Msg
		if err := proto.Unmarshal(b, &msg); err != nil {
			log.Errorln("MsgDecode proto.Unmarshal error: ", err)
			return err
		}
		m.version=msg.Version
		m.id=msg.Id
		m.msgType=MsgType(msg.MsgType)
		if msg.MsgType==pb.MsgType_req||msg.MsgType==pb.MsgType_res{
			m.data=msg.Data
			m.batch=msg.Batch
			m.codecType=CodecType(msg.CodecType)
		}
		return nil
	default:
		return errors.New("this mrpc_serialize is not supported")
	}
}
