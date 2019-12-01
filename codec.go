package rpc


var (
	rpc_codec = RPC_CODEC_PROTOBUF
)

type Codec interface {
	Marshal()([]byte,error)
	Unmarshal(b []byte)(error)
	Reset()
}
