package xml

//ArithRequest defines the request of arith.
type ArithRequest struct {
	A int32 `xml:"a,omitempty"`
	B int32 `xml:"b,omitempty"`
}

//ArithResponse defines the response of arith.
type ArithResponse struct {
	Pro int32 `xml:"pro,omitempty"`
	Quo int32 `xml:"quo,omitempty"`
	Rem int32 `xml:"rem,omitempty"`
}
