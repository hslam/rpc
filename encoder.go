// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"sync"
)

var encoders = sync.Map{}

func init() {
	RegisterHeaderEncoder("json", NewJSONEncoder)
	RegisterHeaderEncoder("code", NewCODEEncoder)
	RegisterHeaderEncoder("pb", NewPBEncoder)
}

// RegisterHeaderEncoder registers a header Encoder.
func RegisterHeaderEncoder(name string, New func() Encoder) {
	encoders.Store(name, New)
}

// NewHeaderEncoder returns a new header Encoder.
func NewHeaderEncoder(name string) func() Encoder {
	if c, ok := encoders.Load(name); ok {
		return c.(func() Encoder)
	}
	return nil
}

// DefaultEncoder returns a default header Encoder.
func DefaultEncoder() Encoder {
	return NewPBEncoder()
}

// Request defines the interface of request.
type Request interface {
	SetSeq(uint64)
	GetSeq() uint64
	SetUpgrade([]byte)
	GetUpgrade() []byte
	SetServiceMethod(string)
	GetServiceMethod() string
	SetArgs([]byte)
	GetArgs() []byte
	Reset()
}

// Response defines the interface of response.
type Response interface {
	SetSeq(uint64)
	GetSeq() uint64
	SetError(string)
	GetError() string
	SetReply([]byte)
	GetReply() []byte
	Reset()
}

// Encoder defines the interface of Encoder.
type Encoder interface {
	NewRequest() Request
	NewResponse() Response
	NewCodec() Codec
}
