package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/transport"
	"github.com/hslam/transport/tcp"
)

//Options defines the struct of options.
type Options struct {
	Transport transport.Transport
	Codec     codec.Codec
	Encoder   *encoder.Encoder
}

//DefaultOptions returns a default options.
func DefaultOptions() *Options {
	return &Options{
		Transport: tcp.NewTransport(),
	}
}
