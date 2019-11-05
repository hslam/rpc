package rpc

import (
	"hslam.com/git/x/rpc/log"
)


type CodecType int32
type MsgType int32
type CompressLevel int32
type CompressType int32
type ReliabilityType int32
const (
	Version			float32	= 1.0

	TCP   					= "tcp"
	UDP						= "udp"
	QUIC  					= "quic"
	WS   					= "ws"
	HTTP   					= "http"
	HTTP2   				= "http2"

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
	DefaultMaxRequests=32

	DefaultServerTimeout	=-1


	DefaultClientTimeout	=-1
	DefaultClientHearbeatTimeout	=1000
	DefaultClientHearbeatTicker=1000
	DefaultClientRetryTicker=10000

	DefaultClientMaxErrPerSecond=10000
	DefaultClientMaxErrHearbeat=3

	MsgTypeReq MsgType = 0
	MsgTypeRes MsgType = 1
	MsgTypeHea MsgType = 2

	NoCompression      CompressLevel= 0
	BestSpeed          CompressLevel= 1
	BestCompression    CompressLevel= 2
	DefaultCompression CompressLevel= 3

	NC			= "no"
	SPEED		= "speed"
	COMPRESSION	= "compression"
	DC			= "default"

	CompressTypeNocom CompressType = 0
	CompressTypeFlate CompressType = 1
	CompressTypeZlib  CompressType = 2
	CompressTypeGzip  CompressType = 3

	NOCOM  = "no"
	FLATE  = "flate"
	ZLIB   = "zlib"
	GZIP   = "gzip"

	DefaultMaxConnNum=1024*64
	DefaultMaxMultiplexingPerConn=64
	DefaultMaxAsyncPerConn =64
)

func SetLogLevel(level log.Level) {
	log.SetLogLevel(level)
}
