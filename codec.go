package rpc

type Codec interface {
	Marshal()([]byte,error)
	Unmarshal(b []byte)(error)
	Reset()
}
