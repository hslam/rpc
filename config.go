package rpc

type CodecType int32
type MsgType int32
type CompressLevel int
type CompressType int32
type ReliabilityType int32
const (
	Version			float32	= 1.0

	IPC   					= "ipc"
	TCP   					= "tcp"
	UDP						= "udp"
	QUIC  					= "quic"
	WS   					= "ws"
	HTTP   					= "http"
	HTTP1   				= "http1"
	HTTP2   				= "http2"

	RPC_CODEC_RAW			= 0
	RPC_CODEC_PROTOBUF		= 1
	RPC_CODEC_GENCODE		= 2

	JSON   					= "json"
	PROTOBUF   				= "pb"
	XML						= "xml"
	BYTES					= "bytes"
	GENCODE					= "gen"
	GOB						= "gob"

	FUNCS_CODEC_INVALID 	CodecType= 0
	FUNCS_CODEC_JSON 		CodecType= 1
	FUNCS_CODEC_PROTOBUF   	CodecType= 2
	FUNCS_CODEC_XML   		CodecType= 3
	FUNCS_CODEC_BYTES   	CodecType= 4
	FUNCS_CODEC_GENCODE   	CodecType= 5
	FUNCS_CODEC_GOB   		CodecType= 9

	DefaultMaxCacheRequest	= 10240
	DefaultMaxBatchRequest	= 1
	DefaultMaxDelayNanoSecond= 1000
	DefaultMaxRequests=32

	DefaultServerTimeout	=-1


	DefaultClientTimeout	=-1
	DefaultClientHearbeatTimeout	=1000
	DefaultClientHearbeatTicker=1000
	DefaultClientRetryTicker=10000

	DefaultClientMaxErrPerSecond=10000
	DefaultClientMaxErrHearbeat=60

	MsgTypeReq MsgType = 0
	MsgTypeRes MsgType = 1
	MsgTypeHea MsgType = 2

	NoCompression      CompressLevel= 0
	BestSpeed          CompressLevel= 1
	BestCompression    CompressLevel= 9
	DefaultCompression CompressLevel= -1


	NC			= "no"
	SPEED		= "speed"
	COMPRESSION	= "compression"
	DC			= "default"

	CompressTypeNo CompressType = 0
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

 	HttpConnected = "200 Connected to RPC"
	HttpPath = "/"
)

