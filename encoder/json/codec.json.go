package json

type Request struct {
	Seq           uint64 `json:"i"`
	ServiceMethod string `json:"m"`
	Args          []byte `json:"p"`
}

type Response struct {
	Seq   uint64 `json:"i"`
	Error string `json:"e"`
	Reply []byte `json:"r"`
}
