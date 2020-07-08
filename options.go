package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/socket"
	"github.com/hslam/socket/tcp"
)

//Options defines the struct of options.
type Options struct {
	NewSocket  func() socket.Socket
	NewCodec   func() codec.Codec
	NewEncoder func() *encoder.Encoder
	NewMessage func() socket.Message
}

//DefaultOptions returns a default options.
func DefaultOptions() *Options {
	return &Options{
		NewSocket: tcp.NewSocket,
	}
}
