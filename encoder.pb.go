// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/codec"
)

// NewPBEncoder returns a header Encoder.
func NewPBEncoder() *Encoder {
	return NewEncoder(NewPBRequest(), NewPBResponse(), NewPBCodec())
}

// NewPBCodec returns the instance of Codec.
func NewPBCodec() codec.Codec {
	return &codec.GOGOPBCodec{}
}

// NewPBRequest returns the instance of pbRequest.
func NewPBRequest() Request {
	return &pbRequest{}
}

// Reset resets the pbRequest.
func (req *pbRequest) Reset() {
	*req = pbRequest{}
}

// SetSeq sets the value of Seq.
func (req *pbRequest) SetSeq(seq uint64) {
	req.Seq = seq
}

// GetSeq returns the value of Seq.
func (req *pbRequest) GetSeq() uint64 {
	return req.Seq
}

// SetUpgrade sets the value of Upgrade.
func (req *pbRequest) SetUpgrade(upgrade []byte) {
	req.Upgrade = upgrade
}

// GetUpgrade returns the value of Upgrade.
func (req *pbRequest) GetUpgrade() []byte {
	return req.Upgrade
}

// SetServiceMethod sets the value of ServiceMethod.
func (req *pbRequest) SetServiceMethod(serviceMethod string) {
	req.ServiceMethod = serviceMethod
}

// GetServiceMethod returns the value of ServiceMethod.
func (req *pbRequest) GetServiceMethod() string {
	return req.ServiceMethod
}

// SetArgs sets the value of Args.
func (req *pbRequest) SetArgs(args []byte) {
	req.Args = args
}

// GetArgs returns the value of Args.
func (req *pbRequest) GetArgs() []byte {
	return req.Args
}

// NewPBResponse returns the instance of pbResponse.
func NewPBResponse() Response {
	return &pbResponse{}
}

// Reset resets the pbResponse.
func (res *pbResponse) Reset() {
	*res = pbResponse{}
}

// SetSeq sets the value of Seq.
func (res *pbResponse) SetSeq(seq uint64) {
	res.Seq = seq
}

// GetSeq returns the value of Seq.
func (res *pbResponse) GetSeq() uint64 {
	return res.Seq
}

// SetError sets the value of Error.
func (res *pbResponse) SetError(errorMsg string) {
	res.Error = errorMsg
}

// GetError returns the value of Error.
func (res *pbResponse) GetError() string {
	return res.Error
}

// SetReply sets the value of Reply.
func (res *pbResponse) SetReply(reply []byte) {
	res.Reply = reply
}

// GetReply returns the value of Reply.
func (res *pbResponse) GetReply() []byte {
	return res.Reply
}
