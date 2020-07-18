package json

type Request struct {
	Seq           uint64 `json:"i"`
	Upgrade       []byte `json:"u"`
	ServiceMethod string `json:"m"`
	Args          []byte `json:"p"`
}

type Response struct {
	Seq     uint64 `json:"i"`
	Upgrade []byte `json:"u"`
	Error   string `json:"e"`
	Reply   []byte `json:"r"`
}
