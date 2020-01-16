package rpc

//Codec defines the interface of codec.
type Codec interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	Reset()
}

var rpcCodec = RPCCodecCode

//SETRPCCODEC sets rpc codec.
func SETRPCCODEC(t BasicCodecType) {
	rpcCodec = t
}
