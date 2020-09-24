// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/funcs"
	"github.com/hslam/rpc/encoder/code"
	"github.com/hslam/rpc/encoder/codepb"
	"github.com/hslam/rpc/encoder/json"
	"github.com/hslam/rpc/encoder/pb"
	"sync"
)

var codecs = sync.Map{}

func init() {
	RegisterCodec("json", json.NewCodec)
	RegisterCodec("code", code.NewCodec)
	RegisterCodec("codepb", codepb.NewCodec)
	RegisterCodec("pb", pb.NewCodec)
}

// RegisterCodec registers a codec.
func RegisterCodec(name string, New func() codec.Codec) {
	codecs.Store(name, New)
}

// NewCodec returns a new Codec.
func NewCodec(name string) func() codec.Codec {
	if c, ok := codecs.Load(name); ok {
		return c.(func() codec.Codec)
	}
	return nil
}

// Context is an RPC context for codec.
type Context struct {
	Seq           uint64
	Upgrade       []byte
	ServiceMethod string
	Error         string
	heartbeat     bool
	noRequest     bool
	noResponse    bool
	keepReading   bool
	f             *funcs.Func
	args          funcs.Value
	reply         funcs.Value
}

// Reset resets the Context.
func (ctx *Context) Reset() {
	*ctx = Context{}
}

// ServerCodec implements reading of RPC requests and writing of
// RPC responses for the server side of an RPC session.
// The server calls ReadRequestHeader and ReadRequestBody in pairs
// to read requests from the connection, and it calls WriteResponse to
// write a response back. The server calls Close when finished with the
// connection. ReadRequestBody may be called with a nil
// argument to force the body of the request to be read and discarded.
type ServerCodec interface {
	ReadRequestHeader(*Context) error
	ReadRequestBody(interface{}) error
	WriteResponse(*Context, interface{}) error
	Close() error
}

// ClientCodec implements writing of RPC requests and
// reading of RPC responses for the client side of an RPC session.
// The client calls WriteRequest to write a request to the connection
// and calls ReadResponseHeader and ReadResponseBody in pairs
// to read responses. The client calls Close when finished with the
// connection. ReadResponseBody may be called with a nil
// argument to force the body of the response to be read and then
// discarded.
type ClientCodec interface {
	WriteRequest(*Context, interface{}) error
	ReadResponseHeader(*Context) error
	ReadResponseBody(interface{}) error
	Close() error
}
