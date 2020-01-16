package rpc

import (
	"errors"
	"github.com/hslam/codec"
)

func argsEncode(args interface{}, funcsCodecType CodecType) ([]byte, error) {
	codec := funcsCodec(funcsCodecType)
	reqBytes, err := codec.Encode(args)
	if err != nil {
		logger.Errorln("ArgsEncode error: ", err)
		return nil, err
	}
	return reqBytes, nil
}

func argsDecode(argsBytes []byte, args interface{}, funcsCodecType CodecType) error {
	codec := funcsCodec(funcsCodecType)
	err := codec.Decode(argsBytes, args)
	if err != nil {
		logger.Errorln("ArgsDecode error: ", err)
		return err
	}
	return nil
}

func replyEncode(reply interface{}, funcsCodecType CodecType) ([]byte, error) {
	codec := funcsCodec(funcsCodecType)
	resBytes, err := codec.Encode(reply)
	if err != nil {
		logger.Errorln("ReplyEncode error: ", err)
		return nil, err
	}
	return resBytes, nil
}

func replyDecode(replyBytes []byte, reply interface{}, funcsCodecType CodecType) error {
	codec := funcsCodec(funcsCodecType)
	err := codec.Decode(replyBytes, reply)
	if err != nil {
		logger.Errorln("ArgsDecode error: ", err)
		return err
	}
	return nil
}

func funcsCodecType(codec string) (CodecType, error) {
	switch codec {
	case JSON:
		return FuncsCodecJSON, nil
	case PROTOBUF:
		return FuncsCodecPROTOBUF, nil
	case XML:
		return FuncsCodecXML, nil
	case GOB:
		return FuncsCodecGOB, nil
	case BYTES:
		return FuncsCodecBYTES, nil
	case CODE:
		return FuncsCodecCODE, nil
	default:
		return FuncsCodecINVALID, errors.New("this codec is not supported")
	}
}

func funcsCodecName(funcsCodecType CodecType) string {
	switch funcsCodecType {
	case FuncsCodecJSON:
		return JSON
	case FuncsCodecPROTOBUF:
		return PROTOBUF
	case FuncsCodecXML:
		return XML
	case FuncsCodecGOB:
		return GOB
	case FuncsCodecBYTES:
		return BYTES
	case FuncsCodecCODE:
		return CODE
	default:
		return ""
	}
}
func funcsCodec(funcsCodecType CodecType) codec.Codec {
	switch funcsCodecType {
	case FuncsCodecJSON:
		return &codec.JSONCodec{}
	case FuncsCodecPROTOBUF:
		return &codec.GOGOPBCodec{}
	case FuncsCodecXML:
		return &codec.XMLCodec{}
	case FuncsCodecGOB:
		return &codec.GOBCodec{}
	case FuncsCodecBYTES:
		return &codec.BYTESCodec{}
	case FuncsCodecCODE:
		return &codec.CODECodec{}
	default:
		return nil
	}
}
