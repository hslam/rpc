package rpc

import (
	"errors"
	"hslam.com/mgit/Mort/rpc/log"
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
	GOB						= "gob"
	BYTES					= "bytes"

	FUNCS_CODEC_INVALID 	CodecType= 0
	FUNCS_CODEC_JSON 		CodecType= 1
	FUNCS_CODEC_PROTOBUF   	CodecType= 2
	FUNCS_CODEC_XML   		CodecType= 3
	FUNCS_CODEC_GOB   		CodecType= 4
	FUNCS_CODEC_BYTES   	CodecType= 5

	DefaultMaxCacheRequest	= 1024
	DefaultMaxBatchRequest	= 8
	DefaultMaxDelayNanoSecond= 1000
	DefaultMaxConcurrentRequest=32

	DefaultServerTimeout	=-1


	DefaultClientTimeout	=60000
	DefaultClientHearbeatTimeout	=10000
	DefaultClientHearbeatTicker=10000
	DefaultClientRetryTicker=10000

	DefaultClientMaxErrPerSecond=1000
	DefaultClientMaxErrHearbeat=3

	MsgTypeReq MsgType = 0
	MsgTypeRes MsgType = 1
	MsgTypeHea MsgType = 2

)


var (
	ErrConnExit=errors.New("exit")
	RPCConnNoResponse=errors.New("RPC NoResponse")
	ErrHystrix=errors.New("Hystrix")
	ErrSetClientID=errors.New("0<=ClientID<=1023")
	ErrSetMaxErrPerSecond=errors.New("maxErrPerSecond>0")
	ErrSetMaxErrHeartbeat=errors.New("maxErrHeartbeat>0")
	ErrSetMaxBatchRequest=errors.New("maxBatchRequest>0")
	ErrSetTimeout=errors.New("timeout>0")
	ErrRemoteCall=errors.New("RemoteCall cbChan is close")
	ErrTimeOut=errors.New("time out")
	ErrReqId=errors.New("req_id err")
	ErrClientId=errors.New("client_id err")

)

func SetLogLevel(level log.Level) {
	log.SetLogLevel(level)
}
