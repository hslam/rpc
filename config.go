package rpc

import (
	"errors"
)

type	CodecType	int32
type	MsgType		int32

const (
	Version			int32	= 0

	TCP   					= "tcp"
	UDP						= "udp"
	QUIC  					= "quic"
	WS   					= "ws"
	FASTHTTP   				= "fast"

	RPC_CODEC_ME			= 0
	RPC_CODEC_PROTOBUF		= 1

	JSON   					= "json"
	PROTOBUF   				= "pb"
	XML						= "xml"

	FUNCS_CODEC_INVALID 	CodecType= 0
	FUNCS_CODEC_JSON 		CodecType= 1
	FUNCS_CODEC_PROTOBUF   	CodecType= 2
	FUNCS_CODEC_XML   		CodecType= 3

	DefultMaxCacheRequest	= 1024
	DefultMaxBatchRequest	= 8
	DefultMaxDelayNanoSecond= 1000
	DefultMaxConcurrentRequest=32


	MsgTypeReq MsgType = 0
	MsgTypeRes MsgType = 1
)


var (
	ErrConnExit=errors.New("exit")
	RPCConnNoResponse=errors.New("RPC NoResponse")
)