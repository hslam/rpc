package rpc

//BasicCodecType defines the rpc codec type.
type BasicCodecType int32

//CodecType defines the codec type.
type CodecType int32

//MsgType defines the msg type.
type MsgType int32

//CompressLevel defines the compress level.
type CompressLevel int

//CompressType defines the compress type.
type CompressType int32

const (
	//Version defines the version.
	Version float32 = 1.0

	//IPC defines the ipc.
	IPC = "ipc"
	//TCP defines the tcp.
	TCP = "tcp"
	//UDP defines the udp.
	UDP = "udp"
	//QUIC defines the quic.
	QUIC = "quic"
	//WS defines the ws.
	WS = "ws"
	//HTTP defines the http.
	HTTP = "http"
	//HTTP1 defines the http1.
	HTTP1 = "http1"

	//RPCCodecCode defines the RPCCodecType of code.
	RPCCodecCode BasicCodecType = 0
	//RPCCodecProtobuf defines the RPCCodecType of protobuf.
	RPCCodecProtobuf BasicCodecType = 1

	//JSON defines the json.
	JSON = "json"
	//PROTOBUF defines the pb.
	PROTOBUF = "pb"
	//XML defines the xml.
	XML = "xml"
	//BYTES defines the bytes.
	BYTES = "bytes"
	//CODE defines the code.
	CODE = "code"
	//GOB defines the gob.
	GOB = "gob"

	//FuncsCodecINVALID defines the CodecType of invalid.
	FuncsCodecINVALID CodecType = 0
	//FuncsCodecJSON defines the CodecType of json.
	FuncsCodecJSON CodecType = 1
	//FuncsCodecPROTOBUF defines the CodecType of protobuf.
	FuncsCodecPROTOBUF CodecType = 2
	//FuncsCodecXML defines the CodecType of xml.
	FuncsCodecXML CodecType = 3
	//FuncsCodecBYTES defines the CodecType of bytes.
	FuncsCodecBYTES CodecType = 4
	//FuncsCodecCODE defines the CodecType of code.
	FuncsCodecCODE CodecType = 5
	//FuncsCodecGOB defines the CodecType of gob.
	FuncsCodecGOB CodecType = 9

	//DefaultMaxCacheRequest defines the number of max cache requests.
	DefaultMaxCacheRequest = 10240
	//DefaultMaxBatchRequest defines the number of max batch requests.
	DefaultMaxBatchRequest = 1
	//DefaultMaxDelayNanoSecond defines the max delay.
	DefaultMaxDelayNanoSecond = 1000
	//DefaultMaxRequests defines the number of max requests.
	DefaultMaxRequests = 64

	//DefaultServerTimeout defines the server timeout.
	DefaultServerTimeout = -1

	//DefaultClientTimeout defines the client timeout.
	DefaultClientTimeout = -1
	//DefaultClientHearbeatTimeout defines the client hearbeat timeout.
	DefaultClientHearbeatTimeout = 1000
	//DefaultClientHearbeatTicker defines the client hearbeat ticker.
	DefaultClientHearbeatTicker = 1000
	//DefaultClientRetryTicker defines the client retry ticker.
	DefaultClientRetryTicker = 10000
	//DefaultClientMaxErrPerSecond defines the number of client max error per second.
	DefaultClientMaxErrPerSecond = 10000
	//DefaultClientMaxErrHeartbeat defines the number of client max heartbeat error.
	DefaultClientMaxErrHeartbeat = 60

	//MsgTypeReq defines the MsgType of request.
	MsgTypeReq MsgType = 0
	//MsgTypeRes defines the MsgType of response.
	MsgTypeRes MsgType = 1
	//MsgTypeHea defines the MsgType of heartbeat.
	MsgTypeHea MsgType = 2

	//NoCompression defines the compress level of no compression.
	NoCompression CompressLevel = 0
	//BestSpeed defines the compress level of best speed.
	BestSpeed CompressLevel = 1
	//BestCompression defines the compress level of best compression.
	BestCompression CompressLevel = 9
	//DefaultCompression defines the compress level of default compression.
	DefaultCompression CompressLevel = -1

	//NC defines the no compression.
	NC = "no"
	//SPEED defines the speed.
	SPEED = "speed"
	//COMPRESSION defines the compression.
	COMPRESSION = "compression"
	//DC defines the default compression.
	DC = "default"

	//CompressTypeNo defines the compress type of no compression.
	CompressTypeNo CompressType = 0
	//CompressTypeFlate defines the compress type of flate.
	CompressTypeFlate CompressType = 1
	//CompressTypeZlib defines the compress type of zlib.
	CompressTypeZlib CompressType = 2
	//CompressTypeGzip defines the compress type of gzip.
	CompressTypeGzip CompressType = 3

	//NOCOM defines the no compression.
	NOCOM = "no"
	//FLATE defines the flate.
	FLATE = "flate"
	//ZLIB defines the zlib.
	ZLIB = "zlib"
	//GZIP defines the gzip.
	GZIP = "gzip"

	//DefaultMaxConnNum defines the number of max connection.
	DefaultMaxConnNum = 1024 * 64
	//DefaultMaxMultiplexingPerConn defines the number of max multiplexing requests per connection.
	DefaultMaxMultiplexingPerConn = 64
	//DefaultMaxAsyncPerConn defines the number of max async requests per connection.
	DefaultMaxAsyncPerConn = 64
	//HTTPConnected defines the http connected.
	HTTPConnected = "200 Connected to RPC"
	//HTTPPath defines the http path.
	HTTPPath = "/"
)
