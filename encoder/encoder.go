// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package encoder

import (
	"github.com/hslam/codec"
	"reflect"
)

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
	SetUpgrade([]byte)
	GetUpgrade() []byte
	SetError(string)
	GetError() string
	SetReply([]byte)
	GetReply() []byte
	Reset()
}

// Encoder defines the struct of Encoder.
type Encoder struct {
	Request  Request
	Response Response
	Codec    codec.Codec
}

//NewEncoder returns the instance of Encoder.
func NewEncoder(req Request, res Response, codec codec.Codec) *Encoder {
	return &Encoder{Request: req, Response: res, Codec: codec}
}

//Clone returns the copy of Codec.
func (c *Encoder) Clone() *Encoder {
	clone := reflect.New(reflect.Indirect(reflect.ValueOf(c)).Type()).Interface().(*Encoder)
	clone.Request = reflect.New(reflect.Indirect(reflect.ValueOf(c.Request)).Type()).Interface().(Request)
	clone.Response = reflect.New(reflect.Indirect(reflect.ValueOf(c.Response)).Type()).Interface().(Response)
	clone.Codec = reflect.New(reflect.Indirect(reflect.ValueOf(c.Codec)).Type()).Interface().(codec.Codec)
	return clone
}
