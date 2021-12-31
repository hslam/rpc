// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/hslam/buffer"
	"github.com/hslam/funcs"
	"github.com/hslam/socket"
	"sync"
)

var codecs = sync.Map{}

func init() {
	RegisterCodec("json", NewJSONCodec)
	RegisterCodec("code", NewCODECodec)
	RegisterCodec("pb", NewPBCodec)
}

// RegisterCodec registers a codec.
func RegisterCodec(name string, New func() Codec) {
	codecs.Store(name, New)
}

// NewCodec returns a new Codec.
func NewCodec(name string) func() Codec {
	if c, ok := codecs.Load(name); ok {
		return c.(func() Codec)
	}
	return nil
}

// contextKey is a key for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

// String returns a context key.
func (k *contextKey) String() string { return "github.com/hslam/rpc context key " + k.name }

// BufferContextKey is a context key.
var BufferContextKey = &contextKey{"buffer"}

const (
	bufferSize = 65536
)

func checkBuffer(buf []byte, n int) []byte {
	if cap(buf) >= n {
		buf = buf[:n]
	} else {
		buf = make([]byte, n)
	}
	return buf
}

// GetBuffer gets a buffer from the pool.
func GetBuffer(size int) []byte {
	if size > 0 {
		return buffer.AssignPool(size).GetBuffer(size)
	}
	return nil
}

// PutBuffer puts a buffer to the pool.
func PutBuffer(buf []byte) {
	size := cap(buf)
	if size > 0 {
		buffer.AssignPool(size).PutBuffer(buf[:size])
	}
}

// BufferSize sets the buffer size.
type BufferSize interface {
	SetBufferSize(int)
}

// GetContextBuffer gets a buffer from the context.
func GetContextBuffer(ctx context.Context) (buffer []byte) {
	value := ctx.Value(BufferContextKey)
	if value != nil {
		if b, ok := value.([]byte); ok {
			buffer = b
		}
	}
	return
}

// FreeContextBuffer frees the context buffer to the pool.
func FreeContextBuffer(ctx context.Context) {
	PutBuffer(GetContextBuffer(ctx))
}

// Context is an RPC context for codec.
type Context struct {
	Seq           uint64
	Upgrade       []byte
	ServiceMethod string
	Error         string
	buffer        []byte
	data          []byte
	upgrade       *upgrade
	f             *funcs.Func
	args          funcs.Value
	reply         funcs.Value
	sending       *sync.Mutex
	codec         ServerCodec
	value         []byte
	stream        *stream
	ctx           *Context
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
	Messages() socket.Messages
	Concurrency() int
	ReadRequestHeader(*Context) error
	ReadRequestBody([]byte, interface{}) error
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
	Messages() socket.Messages
	Concurrency() int
	WriteRequest(*Context, interface{}) error
	ReadResponseHeader(*Context) error
	ReadResponseBody([]byte, interface{}) error
	Close() error
}

// Codec defines the interface for encoding/decoding.
type Codec interface {
	Marshal(buf []byte, v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// JSONCodec struct
type JSONCodec struct {
}

// Marshal returns the JSON encoding of v.
func (c *JSONCodec) Marshal(buf []byte, v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
func (c *JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// XMLCodec struct
type XMLCodec struct {
}

// Marshal returns the XML encoding of v.
func (c *XMLCodec) Marshal(buf []byte, v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// Unmarshal parses the XML-encoded data and stores the result in the value pointed to by v.
func (c *XMLCodec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// BYTESCodec struct
type BYTESCodec struct {
}

// Marshal returns the BYTES encoding of v.
func (c *BYTESCodec) Marshal(buf []byte, v interface{}) ([]byte, error) {
	return *v.(*[]byte), nil
}

// Unmarshal parses the BYTES-encoded data and stores the result in the value pointed to by v.
func (c *BYTESCodec) Unmarshal(data []byte, v interface{}) error {
	*v.(*[]byte) = data
	return nil
}

// GoGoProtobuf defines the interface for gogo's protobuf.
type GoGoProtobuf interface {
	Size() (n int)
	Marshal() (data []byte, err error)
	MarshalTo(buf []byte) (int, error)
	Unmarshal(data []byte) error
}

// GOGOPBCodec struct
type GOGOPBCodec struct {
}

//ErrorGOGOPB is the error that v is not GoGoProtobuf
var ErrorGOGOPB = errors.New("is not GoGoProtobuf")

// Marshal returns the GOGOPB encoding of v.
func (c *GOGOPBCodec) Marshal(buf []byte, v interface{}) ([]byte, error) {
	if p, ok := v.(GoGoProtobuf); ok {
		size := p.Size()
		if cap(buf) >= size {
			buf = buf[:size]
			n, err := p.MarshalTo(buf)
			return buf[:n], err
		}
		return p.Marshal()
	}
	return nil, ErrorGOGOPB
}

// Unmarshal parses the GOGOPB-encoded data and stores the result in the value pointed to by v.
func (c *GOGOPBCodec) Unmarshal(data []byte, v interface{}) error {
	if p, ok := v.(GoGoProtobuf); ok {
		return p.Unmarshal(data)
	}
	return ErrorGOGOPB
}

// Code defines the interface for code.
type Code interface {
	Marshal(buf []byte) ([]byte, error)
	Unmarshal(buf []byte) (uint64, error)
}

//ErrorCODE is the error that v is not Code
var ErrorCODE = errors.New("is not Code")

// CODECodec struct
type CODECodec struct {
}

// Marshal returns the CODE encoding of v.
func (c *CODECodec) Marshal(buf []byte, v interface{}) ([]byte, error) {
	if p, ok := v.(Code); ok {
		return p.Marshal(buf)
	}
	return nil, ErrorCODE
}

// Unmarshal parses the CODE-encoded data and stores the result in the value pointed to by v.
func (c *CODECodec) Unmarshal(data []byte, v interface{}) error {
	if p, ok := v.(Code); ok {
		_, err := p.Unmarshal(data)
		return err
	}
	return ErrorCODE
}

// MsgPack defines the interface for msgp.
type MsgPack interface {
	MarshalMsg(buf []byte) ([]byte, error)
	UnmarshalMsg(bts []byte) (o []byte, err error)
}

//ErrorMSGP is the error that v is not MSGP
var ErrorMSGP = errors.New("is not MSGP")

// MSGPCodec struct
type MSGPCodec struct {
}

// Marshal returns the MSGP encoding of v.
func (c *MSGPCodec) Marshal(buf []byte, v interface{}) ([]byte, error) {
	if p, ok := v.(MsgPack); ok {
		return p.MarshalMsg(buf)
	}
	return nil, ErrorMSGP
}

// Unmarshal parses the MSGP-encoded data and stores the result in the value pointed to by v.
func (c *MSGPCodec) Unmarshal(data []byte, v interface{}) error {
	if p, ok := v.(MsgPack); ok {
		_, err := p.UnmarshalMsg(data)
		return err
	}
	return ErrorMSGP
}
