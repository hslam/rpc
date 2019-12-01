package rpc

import (
	"bytes"
	"errors"
	"hslam.com/git/x/rpc/pb"
	"hslam.com/git/x/compress"
	"hslam.com/git/x/rpc/gen"
)

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
		return m.Marshal(nil)
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
		if data,err:=msg.Marshal();err!=nil{
			Errorln("MsgEncode proto.Unmarshal error: ", err)
			return nil,err
		}else {
			return data,nil
		}
	case RPC_CODEC_GENCODE:
		var msg gen.Msg
		if m.msgType==MsgType(pb.MsgType_req)||m.msgType==MsgType(pb.MsgType_res){
			msg= gen.Msg{
				Version:Version,
				Id:m.id,
				MsgType:int32(m.msgType),
				Batch:m.batch,
				Data:m.data,
				CodecType:int32(m.codecType),
				CompressType:int32(m.compressType),
				CompressLevel:int32(m.compressLevel),
			}
		}else if m.msgType==MsgTypeHea {
			msg= gen.Msg{Version:Version,Id:m.id,MsgType:int32(m.msgType)}
		}
		if data,err:=msg.Marshal(nil);err!=nil{
			Errorln("MsgEncode proto.Unmarshal error: ", err)
			return nil,err
		}else {
			return data,nil
		}
	}
	return nil,errors.New("this rpc_serialize is not supported")
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
		return m.Unmarshal(b)
	case RPC_CODEC_PROTOBUF:
		var msg =&pb.Msg{}
		if err := msg.Unmarshal(b); err != nil {
			Errorln("MsgDecode proto.Unmarshal error: ", err)
			return err
		}
		m.version=msg.Version
		m.id=msg.Id
		m.msgType=MsgType(msg.MsgType)
		if m.msgType==MsgTypeReq||m.msgType==MsgTypeRes{
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
	case RPC_CODEC_GENCODE:
		var msg =&gen.Msg{}
		if _,err := msg.Unmarshal(b); err != nil {
			Errorln("MsgDecode gencode.Unmarshal error: ", err)
			return err
		}
		m.version=msg.Version
		m.id=msg.Id
		m.msgType=MsgType(msg.MsgType)
		if m.msgType==MsgTypeReq||m.msgType==MsgTypeRes{
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
		return errors.New("this rpc_serialize is not supported")
	}
}

func(m *Msg)Marshal(buf []byte)([]byte,error)  {
	return nil,nil
}
func(m *Msg)Unmarshal(b []byte)(error)  {
	return nil
}
func(m *Msg)Reset()()  {
}
func getCompressor(compressType CompressType,level CompressLevel)  (compress.Compressor)  {
	if level==NoCompression{
		return nil
	}
	switch compressType {
	case CompressTypeFlate:
		return &compress.FlateCompressor{Level:compress.Level(level)}
	case CompressTypeZlib:
		return &compress.ZlibCompressor{Level:compress.Level(level)}
	case CompressTypeGzip:
		return &compress.GzipCompressor{Level:compress.Level(level)}
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

