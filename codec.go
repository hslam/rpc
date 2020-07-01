package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/funcs"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/rpc/encoder/code"
	"sync"
)

var codecs = sync.Map{}

func init() {
	RegisterCodec("json", &codec.JSONCodec{})
	RegisterCodec("code", &codec.CODECodec{})
	RegisterCodec("pb", &codec.GOGOPBCodec{})

}

func RegisterCodec(name string, codec codec.Codec) {
	codecs.Store(name, codec)
}

func NewCodec(name string) codec.Codec {
	if c, ok := codecs.Load(name); ok {
		return c.(codec.Codec)
	}
	return nil
}

func NewDefaultEncoder() *encoder.Encoder {
	return encoder.NewEncoder(code.NewRequest(), code.NewResponse(), &codec.CODECodec{})
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
