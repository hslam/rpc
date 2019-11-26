package rpc

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"errors"
	"hslam.com/git/x/rpc/pb"
)

var (
	rpc_codec = RPC_CODEC_PROTOBUF
)

func init() {
}

func (m Msg)Serialize(version float32,method string,data []byte)[]byte  {
	header:=append(float32ToByte(version),uint16ToBytes(uint16(len(method)))...)
	var Bytes=append(append(header, []byte(method)...), data...)
	return Bytes
}

func (m Msg)Deserialize(msg []byte)(float32,string,[]byte)  {
	var buf =bytes.NewBuffer(msg)
	version_bytes:=make([]byte , 4)
	buf.Read(version_bytes)
	version:=byteToFloat32(version_bytes)
	method_len_bytes:=make([]byte , 2)
	method_len:=bytesToUint16(method_len_bytes)
	method_bytes_len:=int(method_len)
	method_bytes:=make([]byte,method_bytes_len)
	buf.Read(method_bytes)
	method:=string(method_bytes)
	data:=make([]byte,len(msg)-method_bytes_len-6)
	buf.Read(data)
	return version,method,data
}

type Msg struct {
	version				float32
	id                  int64
	msgType				MsgType
	batch				bool
	codecType			CodecType
	compressType		CompressType
	compressLevel 		CompressLevel
	data				[]byte

}
func (m *Msg)Encode() ([]byte, error) {
	compressor:=getCompressor(m.compressType,m.compressLevel)
	if compressor!=nil{
		m.data,_=compressor.Compress(m.data)
	}
	switch rpc_codec {
	case RPC_CODEC_ME:
		return m.data,nil
	case RPC_CODEC_PROTOBUF:
		var msg pb.Msg
		if m.msgType==MsgType(pb.MsgType_req)||m.msgType==MsgType(pb.MsgType_res){
			msg=pb.Msg{
				Version:Version,
				Id:m.id,
				MsgType:pb.MsgType(m.msgType),
				Batch:m.batch,
				Data:m.data,
				CodecType:pb.CodecType(m.codecType),
				CompressType:pb.CompressType(m.compressType),
				CompressLevel:pb.CompressLevel(m.compressLevel),
			}
		}else if m.msgType==MsgType(pb.MsgType_hea) {
			msg=pb.Msg{Version:Version,Id:m.id,MsgType:pb.MsgType(m.msgType)}
		}
		if data,err:=proto.Marshal(&msg);err!=nil{
			Errorln("MsgEncode proto.Unmarshal error: ", err)
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
	m.compressType=CompressTypeNo
	m.compressLevel=NoCompression
	switch rpc_codec {
	case RPC_CODEC_ME:
		m.version=Version
		m.data=b
		return nil
	case RPC_CODEC_PROTOBUF:
		var msg pb.Msg
		if err := proto.Unmarshal(b, &msg); err != nil {
			Errorln("MsgDecode proto.Unmarshal error: ", err)
			return err
		}
		m.version=msg.Version
		m.id=msg.Id
		m.msgType=MsgType(msg.MsgType)
		if msg.MsgType==pb.MsgType_req||msg.MsgType==pb.MsgType_res{
			m.data=msg.Data
			m.batch=msg.Batch
			m.codecType=CodecType(msg.CodecType)
			m.compressType=CompressType(msg.CompressType)
			m.compressLevel=CompressLevel(msg.CompressLevel)
			compressor:=getCompressor(m.compressType,m.compressLevel)
			if compressor!=nil{
				m.data,_=compressor.Uncompress(m.data)
			}
		}
		return nil
	default:
		return errors.New("this mrpc_serialize is not supported")
	}
}

func getCompressor(compressType CompressType,level CompressLevel)  (Compressor)  {
	if level==NoCompression{
		return nil
	}
	switch compressType {
	case CompressTypeFlate:
		return &FlateCompressor{Level:level}
	case CompressTypeZlib:
		return &ZlibCompressor{Level:level}
	case CompressTypeGzip:
		return &GzipCompressor{Level:level}
	default:
		return nil
	}
}

func getCompressType(name string)  (CompressType)  {
	switch name {
	case FLATE:
		return CompressTypeFlate
	case ZLIB:
		return CompressTypeZlib
	case GZIP:
		return CompressTypeGzip
	default:
		return CompressTypeNo
	}
}
func getCompressLevel(name string)  (CompressLevel)  {
	switch name {
	case SPEED:
		return BestSpeed
	case COMPRESSION:
		return BestCompression
	case DC:
		return DefaultCompression
	default:
		return NoCompression
	}
}
