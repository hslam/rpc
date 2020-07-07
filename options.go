package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/transport"
	"github.com/hslam/transport/tcp"
)

//Options defines the struct of options.
type Options struct {
	NewTransport   func() transport.Transport
	NewCodec       func() codec.Codec
	NewEncoder     func() *encoder.Encoder
	NewMessageConn func() MessageConn
}

//DefaultOptions returns a default options.
func DefaultOptions() *Options {
	return &Options{
		NewTransport: tcp.NewTransport,
	}
}
