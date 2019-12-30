package rpc

import (
	"errors"
	"github.com/hslam/rpc/pb"
	"github.com/hslam/compress"
	"github.com/hslam/code"
)

type Msg struct {
	version				float32
	id                  uint64
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
	m.version=Version
	switch rpc_codec {
	case RPC_CODEC_CODE:
		return m.Marshal(nil)
	case RPC_CODEC_PROTOBUF:
		var msg pb.Msg
		if m.msgType==MsgType(pb.MsgType_req)||m.msgType==MsgType(pb.MsgType_res){
			msg=pb.Msg{
				Version:m.version,
				Id:m.id,
				MsgType:pb.MsgType(m.msgType),
				Batch:m.batch,
				Data:m.data,
				CodecType:pb.CodecType(m.codecType),
				CompressType:pb.CompressType(m.compressType),
				CompressLevel:pb.CompressLevel(m.compressLevel),
			}
		}else if m.msgType==MsgType(pb.MsgType_hea) {
			msg=pb.Msg{Version:m.version,Id:m.id,MsgType:pb.MsgType(m.msgType)}
		}
		if data,err:=msg.Marshal();err!=nil{
			Errorln("MsgEncode proto.Unmarshal error: ", err)
			return nil,err
		}else {
			return data,nil
		}
	}
	return nil,errors.New("this rpc_serialize is not supported")
}
func (m *Msg)Decode(b []byte)(error) {
	m.version=0
	m.id=0
	m.data=nil
	m.batch=false
	m.codecType=FUNCS_CODEC_INVALID
	m.compressType=CompressTypeNo
	m.compressLevel=NoCompression
	switch rpc_codec {
	case RPC_CODEC_CODE:
		err:= m.Unmarshal(b)
		compressor:=getCompressor(m.compressType,m.compressLevel)
		if compressor!=nil{
			m.data,_=compressor.Uncompress(m.data)
		}
		return err
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
	default:
		return errors.New("this rpc_serialize is not supported")
	}
}

func(m *Msg)Marshal(buf []byte)([]byte,error)  {
	var size uint64
	size+=4
	size+=8
	size+=4
	size+=1
	size+=4
	size+=4
	size+=4
	size+=code.SizeofBytes(m.data)
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n  uint64
	n = code.EncodeFloat32(buf[offset:],m.version)
	offset+=n
	n = code.EncodeUint64(buf[offset:],m.id)
	offset+=n
	n = code.EncodeUint32(buf[offset:],uint32(m.msgType))
	offset+=n
	n = code.EncodeBool(buf[offset:],m.batch)
	offset+=n
	n = code.EncodeUint32(buf[offset:],uint32(m.codecType))
	offset+=n
	n = code.EncodeUint32(buf[offset:],uint32(m.compressType))
	offset+=n
	n = code.EncodeUint32(buf[offset:],uint32(m.compressLevel))
	offset+=n
	n = code.EncodeBytes(buf[offset:],m.data)
	offset+=n
	return buf,nil
}
func(m *Msg)Unmarshal(b []byte)(error)  {
	var offset uint64
	var n uint64
	n=code.DecodeFloat32(b[offset:],&m.version)
	offset+=n
	n=code.DecodeUint64(b[offset:],&m.id)
	offset+=n
	var msgType uint32
	n=code.DecodeUint32(b[offset:],&msgType)
	m.msgType=MsgType(msgType)
	offset+=n
	n=code.DecodeBool(b[offset:],&m.batch)
	offset+=n
	var codecType uint32
	n=code.DecodeUint32(b[offset:],&codecType)
	m.codecType=CodecType(codecType)
	offset+=n
	var compressType uint32
	n=code.DecodeUint32(b[offset:],&compressType)
	m.compressType=CompressType(compressType)
	offset+=n
	var compressLevel uint32
	n=code.DecodeUint32(b[offset:],&compressLevel)
	m.compressLevel=CompressLevel(compressLevel)
	offset+=n
	n=code.DecodeBytes(b[offset:],&m.data)
	offset+=n
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

