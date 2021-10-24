// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

// JSONEncoder implements a header Encoder.
type JSONEncoder struct{}

// NewJSONEncoder returns a header Encoder.
func NewJSONEncoder() Encoder {
	return &JSONEncoder{}
}

// NewRequest returns the instance of Request.
func (e *JSONEncoder) NewRequest() Request {
	return NewJSONRequest()
}

// NewResponse returns the instance of Response.
func (e *JSONEncoder) NewResponse() Response {
	return NewJSONResponse()
}

// NewCodec returns the instance of Codec.
func (e *JSONEncoder) NewCodec() Codec {
	return NewJSONCodec()
}

// NewJSONCodec returns the instance of Codec.
func NewJSONCodec() Codec {
	return &JSONCodec{}
}

// NewJSONRequest returns the instance of jsonRequest.
func NewJSONRequest() Request {
	return &jsonRequest{}
}

// Reset resets the jsonRequest.
func (req *jsonRequest) Reset() {
	*req = jsonRequest{}
}

// SetSeq sets the value of Seq.
func (req *jsonRequest) SetSeq(seq uint64) {
	req.Seq = seq
}

// GetSeq returns the value of Seq.
func (req *jsonRequest) GetSeq() uint64 {
	return req.Seq
}

// SetUpgrade sets the value of Upgrade.
func (req *jsonRequest) SetUpgrade(upgrade []byte) {
	req.Upgrade = upgrade
}

// GetUpgrade returns the value of Upgrade.
func (req *jsonRequest) GetUpgrade() []byte {
	return req.Upgrade
}

// SetServiceMethod sets the value of ServiceMethod.
func (req *jsonRequest) SetServiceMethod(serviceMethod string) {
	req.ServiceMethod = serviceMethod
}

// GetServiceMethod returns the value of ServiceMethod.
func (req *jsonRequest) GetServiceMethod() string {
	return req.ServiceMethod
}

// SetArgs sets the value of Args.
func (req *jsonRequest) SetArgs(args []byte) {
	req.Args = args
}

// GetArgs returns the value of Args.
func (req *jsonRequest) GetArgs() []byte {
	return req.Args
}

// NewJSONResponse returns the instance of jsonResponse.
func NewJSONResponse() Response {
	return &jsonResponse{}
}

// Reset resets the jsonResponse.
func (res *jsonResponse) Reset() {
	*res = jsonResponse{}
}

// SetSeq sets the value of Seq.
func (res *jsonResponse) SetSeq(seq uint64) {
	res.Seq = seq
}

// GetSeq returns the value of Seq.
func (res *jsonResponse) GetSeq() uint64 {
	return res.Seq
}

// SetError sets the value of Error.
func (res *jsonResponse) SetError(errorMsg string) {
	res.Error = errorMsg
}

// GetError returns the value of Error.
func (res *jsonResponse) GetError() string {
	return res.Error
}

// SetReply sets the value of Reply.
func (res *jsonResponse) SetReply(reply []byte) {
	res.Reply = reply
}

// GetReply returns the value of Reply.
func (res *jsonResponse) GetReply() []byte {
	return res.Reply
}
