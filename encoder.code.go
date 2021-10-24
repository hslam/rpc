// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

// CODEEncoder implements a header Encoder.
type CODEEncoder struct{}

// NewCODEEncoder returns a header Encoder.
func NewCODEEncoder() Encoder {
	return &CODEEncoder{}
}

// NewRequest returns the instance of Request.
func (e *CODEEncoder) NewRequest() Request {
	return NewCODERequest()
}

// NewResponse returns the instance of Response.
func (e *CODEEncoder) NewResponse() Response {
	return NewCODEResponse()
}

// NewCodec returns the instance of Codec.
func (e *CODEEncoder) NewCodec() Codec {
	return NewCODECodec()
}

// NewCODECodec returns the instance of Codec.
func NewCODECodec() Codec {
	return &CODECodec{}
}

// NewCODERequest returns the instance of codeRequest.
func NewCODERequest() Request {
	return &request{}
}

//Reset resets the Request.
func (req *request) Reset() {
	*req = request{}
}

//SetSeq sets the value of Seq.
func (req *request) SetSeq(seq uint64) {
	req.Seq = seq
}

//GetSeq returns the value of Seq.
func (req *request) GetSeq() uint64 {
	return req.Seq
}

//SetUpgrade sets the value of Upgrade.
func (req *request) SetUpgrade(upgrade []byte) {
	req.Upgrade = upgrade
}

//GetUpgrade returns the value of Upgrade.
func (req *request) GetUpgrade() []byte {
	return req.Upgrade
}

//SetServiceMethod sets the value of ServiceMethod.
func (req *request) SetServiceMethod(serviceMethod string) {
	req.ServiceMethod = serviceMethod
}

//GetServiceMethod returns the value of ServiceMethod.
func (req *request) GetServiceMethod() string {
	return req.ServiceMethod
}

//SetArgs sets the value of Args.
func (req *request) SetArgs(args []byte) {
	req.Args = args
}

//GetArgs returns the value of Args.
func (req *request) GetArgs() []byte {
	return req.Args
}

// NewCODEResponse returns the instance of codeResponse.
func NewCODEResponse() Response {
	return &response{}
}

//Reset resets the Response.
func (res *response) Reset() {
	*res = response{}
}

//SetSeq sets the value of Seq.
func (res *response) SetSeq(seq uint64) {
	res.Seq = seq
}

//GetSeq returns the value of Seq.
func (res *response) GetSeq() uint64 {
	return res.Seq
}

//SetError sets the value of Error.
func (res *response) SetError(errorMsg string) {
	res.Error = errorMsg
}

//GetError returns the value of Error.
func (res *response) GetError() string {
	return res.Error
}

//SetReply sets the value of Reply.
func (res *response) SetReply(reply []byte) {
	res.Reply = reply
}

//GetReply returns the value of Reply.
func (res *response) GetReply() []byte {
	return res.Reply
}
