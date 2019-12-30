package rpc

type Codec interface {
	Marshal()([]byte,error)
	Unmarshal(b []byte)(error)
	Reset()
}
var rpc_codec = RPC_CODEC_CODE

func SETRPCCODEC(t RPC_CODEC_TYPE)  {
	rpc_codec = t
}

