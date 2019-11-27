package rpc

import (
	"errors"
	"hslam.com/git/x/codec"
)

func ArgsEncode(args interface{},funcsCodecType CodecType) ([]byte, error)  {
	codec:=FuncsCodec(funcsCodecType)
	req_bytes,err :=codec.Encode(args)
	if err!=nil{
		Errorln("ArgsEncode error: ", err)
		return nil,err
	}
	return req_bytes,nil
}

func ArgsDecode(args_bytes []byte,args interface{},funcsCodecType CodecType) (error){
	codec:=FuncsCodec(funcsCodecType)
	err:=codec.Decode(args_bytes,args)
	if err!=nil{
		Errorln("ArgsDecode error: ", err)
		return err
	}
	return nil
}

func ReplyEncode(reply interface{},funcsCodecType CodecType) ([]byte, error)  {
	codec:=FuncsCodec(funcsCodecType)
	res_bytes,err :=codec.Encode(reply)
	if err!=nil{
		Errorln("ReplyEncode error: ", err)
		return nil,err
	}
	return res_bytes,nil
}

func ReplyDecode(reply_bytes []byte,reply interface{},funcsCodecType CodecType) (error){
	codec:=FuncsCodec(funcsCodecType)
	err:=codec.Decode(reply_bytes,reply)
	if err!=nil{
		Errorln("ArgsDecode error: ", err)
		return err
	}
	return nil
}


func FuncsCodecType(codec string)  (CodecType, error)  {
	switch codec {
	case JSON:
		return FUNCS_CODEC_JSON,nil
	case PROTOBUF:
		return FUNCS_CODEC_PROTOBUF,nil
	case XML:
		return FUNCS_CODEC_XML,nil
	case GOB:
		return FUNCS_CODEC_GOB,nil
	case BYTES:
		return FUNCS_CODEC_BYTES,nil
	case GENCODE:
		return FUNCS_CODEC_GENCODE,nil
	default:
		return FUNCS_CODEC_INVALID,errors.New("this codec is not supported")
	}
}

func FuncsCodecName(funcsCodecType CodecType)string  {
	switch funcsCodecType {
	case FUNCS_CODEC_JSON:
		return JSON
	case FUNCS_CODEC_PROTOBUF:
		return PROTOBUF
	case FUNCS_CODEC_XML:
		return XML
	case FUNCS_CODEC_GOB:
		return GOB
	case FUNCS_CODEC_BYTES:
		return BYTES
	case FUNCS_CODEC_GENCODE:
		return GENCODE
	default:
		return ""
	}
}
func FuncsCodec(funcsCodecType CodecType)  (codec.Codec)  {
	switch funcsCodecType {
	case FUNCS_CODEC_JSON:
		return &codec.JsonCodec{}
	case FUNCS_CODEC_PROTOBUF:
		return &codec.FastProtoCodec{}
	case FUNCS_CODEC_XML:
		return &codec.XmlCodec{}
	case FUNCS_CODEC_GOB:
		return &codec.GobCodec{}
	case FUNCS_CODEC_BYTES:
		return &codec.BytesCodec{}
	case FUNCS_CODEC_GENCODE:
		return &codec.GencodeCodec{}
	default:
		return nil
	}
}