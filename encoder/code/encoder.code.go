package code

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
)

func NewEncoder() *encoder.Encoder {
	return encoder.NewEncoder(NewRequest(), NewResponse(), NewCodec())
}

//NewCodec returns the instance of Codec.
func NewCodec() codec.Codec {
	return &codec.CODECodec{}
}

//NewRequest returns the instance of Request.
func NewRequest() *Request {
	return &Request{}
}

//Reset resets the Request.
func (req *Request) Reset() {
	*req = Request{}
}

//SetSeq sets the value of Seq.
func (req *Request) SetSeq(seq uint64) {
	req.Seq = seq
}

//GetSeq returns the value of Seq.
func (req *Request) GetSeq() uint64 {
	return req.Seq
}

//SetServiceMethod sets the value of ServiceMethod.
func (req *Request) SetServiceMethod(serviceMethod string) {
	req.ServiceMethod = serviceMethod
}

//GetServiceMethod returns the value of ServiceMethod.
func (req *Request) GetServiceMethod() string {
	return req.ServiceMethod
}

//SetArgs sets the value of Args.
func (req *Request) SetArgs(args []byte) {
	req.Args = args
}

//GetArgs returns the value of Args.
func (req *Request) GetArgs() []byte {
	return req.Args
}

//NewResponse returns the instance of Response.
func NewResponse() *Response {
	return &Response{}
}

//Reset resets the Response.
func (res *Response) Reset() {
	*res = Response{}
}

//SetSeq sets the value of Seq.
func (res *Response) SetSeq(seq uint64) {
	res.Seq = seq
}

//GetSeq returns the value of Seq.
func (res *Response) GetSeq() uint64 {
	return res.Seq
}

//SetError sets the value of Error.
func (res *Response) SetError(errorMsg string) {
	res.Error = errorMsg
}

//GetError returns the value of Error.
func (res *Response) GetError() string {
	return res.Error
}

//SetReply sets the value of Reply.
func (res *Response) SetReply(reply []byte) {
	res.Reply = reply
}

//GetReply returns the value of Reply.
func (res *Response) GetReply() []byte {
	return res.Reply
}
