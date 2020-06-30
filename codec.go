package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/funcs"
	"reflect"
)

var fast = true

func SetFast(enable bool) {
	fast = enable
}

// Codec defines the struct of codec.
type Codec struct {
	Request  Request
	Response Response
	Codec    codec.Codec
}

//NewCodec returns the instance of Codec.
func NewCodec(req Request, res Response, codec codec.Codec) *Codec {
	return &Codec{Request: req, Response: res, Codec: codec}
}

//Clone returns the copy of Codec.
func (c *Codec) Clone() *Codec {
	clone := reflect.New(reflect.Indirect(reflect.ValueOf(c)).Type()).Interface().(*Codec)
	clone.Request = reflect.New(reflect.Indirect(reflect.ValueOf(c.Request)).Type()).Interface().(Request)
	clone.Response = reflect.New(reflect.Indirect(reflect.ValueOf(c.Response)).Type()).Interface().(Response)
	clone.Codec = reflect.New(reflect.Indirect(reflect.ValueOf(c.Codec)).Type()).Interface().(codec.Codec)
	return clone
}

// Request defines the interface of request.
type Request interface {
	SetSeq(uint64)
	GetSeq() uint64
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

type Context struct {
	Seq           uint64
	ServiceMethod string
	Error         string
	decodeHeader  bool
	keepReading   bool
	f             *funcs.Func
	args          funcs.Value
	reply         funcs.Value
}

func (ctx *Context) Reset() {
	*ctx = Context{}
}

type ServerCodec interface {
	ReadRequestHeader(*Context) error
	ReadRequestBody(interface{}) error
	WriteResponse(*Context, interface{}) error
	Close() error
}

type ClientCodec interface {
	WriteRequest(*Context, interface{}) error
	ReadResponseHeader(*Context) error
	ReadResponseBody(interface{}) error
	Close() error
}
