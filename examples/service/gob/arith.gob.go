package gob

//ArithRequest defines the request of arith.
type ArithRequest struct {
	A int32
	B int32
}

//ArithResponse defines the response of arith.
type ArithResponse struct {
	Pro int32
	Quo int32
	Rem int32
}
