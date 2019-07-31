package rpc

import (
	"github.com/golang/protobuf/proto"
	"hslam.com/mgit/Mort/rpc/log"
	"errors"
	"encoding/json"
	"encoding/xml"
)

func ArgsEncode(args interface{},funcsCodecType CodecType) ([]byte, error)  {
	var (
		req_bytes []byte
		err error
	)
	switch funcsCodecType {
	case FUNCS_CODEC_JSON:
		req_bytes,err =json.Marshal(args)
		if err!=nil{
			log.Errorln("ArgsEncode json.Marshal error: ", err)
			return nil,err
		}
	case FUNCS_CODEC_PROTOBUF:
		req_bytes, err = proto.Marshal(args.(proto.Message))
		if err != nil {
			log.Errorln("ArgsEncode proto.Marshal error: ", err)
			return nil,err
		}
	case FUNCS_CODEC_XML:
		req_bytes,err=xml.Marshal(args)
		if err!=nil{
			log.Errorln("ArgsEncode xml.Marshal error: ", err)
			return nil,err
		}
	default:
		return nil,errors.New("this serialize is not supported")
	}
	return req_bytes,err
}

func ArgsDecode(args_bytes []byte,args interface{},funcsCodecType CodecType) (error){
	switch funcsCodecType {
	case FUNCS_CODEC_JSON:
		err:=json.Unmarshal(args_bytes,args)
		if err!=nil{
			log.Errorln("ArgsDecode json.Unmarshal error: ", err)
			return err
		}
	case FUNCS_CODEC_PROTOBUF:
		err := proto.Unmarshal(args_bytes, args.(proto.Message))
		if err != nil {
			log.Errorln("ArgsDecode proto.Unmarshal error: ", err)
			return err
		}
	case FUNCS_CODEC_XML:
		err:=xml.Unmarshal(args_bytes,args)
		if err!=nil{
			log.Errorln("ArgsDecode xml.Unmarshal error: ", err)
			return err
		}
	default:
		return errors.New("this serialize is not supported")
	}
	return nil
}

func ReplyEncode(reply interface{},funcsCodecType CodecType) ([]byte, error)  {
	var(
		res_bytes []byte
		err error
	)
	switch funcsCodecType {
	case FUNCS_CODEC_JSON:
		res_bytes, err = json.Marshal(reply)
		if err != nil {
			log.Errorln("ReplyEncode json.Marshal error: ", err)
			return nil,err
		}
	case FUNCS_CODEC_PROTOBUF:
		res_bytes, err = proto.Marshal(reply.(proto.Message))
		if err != nil {
			log.Errorln("ReplyEncode proto.Marshal error: ", err)
			return nil,err
		}
	case FUNCS_CODEC_XML:
		res_bytes, err = xml.Marshal(reply)
		if err != nil {
			log.Errorln("ReplyEncode xml.Marshal error: ", err)
			return nil,err
		}
	default:
		return nil,errors.New("this serialize is not supported")
	}
	return res_bytes,nil
}

func ReplyDecode(reply_bytes []byte,reply interface{},funcsCodecType CodecType) (error){
	switch funcsCodecType {
	case FUNCS_CODEC_JSON:
		err := json.Unmarshal(reply_bytes, reply)
		if err != nil {
			log.Errorln("ReplyDecode json.Unmarshal error: ", err)
			return err
		}
		return nil
	case FUNCS_CODEC_PROTOBUF:
		err := proto.Unmarshal(reply_bytes, reply.(proto.Message))
		if err != nil {
			log.Errorln("ReplyDecode proto.Unmarshal error: ", err)
			return err
		}
		return nil
	case FUNCS_CODEC_XML:
		err := xml.Unmarshal(reply_bytes, reply)
		if err != nil {
			log.Errorln("ReplyDecode xml.Unmarshal error: ", err)
			return err
		}
		return nil
	default:
		return errors.New("this serialize is not supported")
	}
	return errors.New("failed")
}


func FuncsCodecType(codec string)  (CodecType, error)  {
	switch codec {
	case JSON:
		return FUNCS_CODEC_JSON,nil
	case PROTOBUF:
		return FUNCS_CODEC_PROTOBUF,nil
	case XML:
		return FUNCS_CODEC_XML,nil
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
	default:
		return ""
	}
}