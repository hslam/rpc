// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package pb

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
)

func NewEncoder() *encoder.Encoder {
	return encoder.NewEncoder(NewRequest(), NewResponse(), NewCodec())
}

//NewCodec returns the instance of Codec.
func NewCodec() codec.Codec {
	return &codec.GOGOPBCodec{}
}

//NewRequest returns the instance of Request.
func NewRequest() *Request {
	return &Request{}
}

//SetSeq sets the value of Seq.
func (req *Request) SetSeq(seq uint64) {
	req.Seq = seq
}

//SetUpgrade sets the value of Upgrade.
func (req *Request) SetUpgrade(upgrade []byte) {
	req.Upgrade = upgrade
}

//SetServiceMethod sets the value of ServiceMethod.
func (req *Request) SetServiceMethod(serviceMethod string) {
	req.ServiceMethod = serviceMethod
}

//SetArgs sets the value of Args.
func (req *Request) SetArgs(args []byte) {
	req.Args = args
}

//NewResponse returns the instance of Response.
func NewResponse() *Response {
	return &Response{}
}

//SetSeq sets the value of Seq.
func (res *Response) SetSeq(seq uint64) {
	res.Seq = seq
}

//SetUpgrade sets the value of Upgrade.
func (res *Response) SetUpgrade(upgrade []byte) {
	res.Upgrade = upgrade
}

//SetError sets the value of Error.
func (res *Response) SetError(errorMsg string) {
	res.Error = errorMsg
}

//SetReply sets the value of Reply.
func (res *Response) SetReply(reply []byte) {
	res.Reply = reply
}
