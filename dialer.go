package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"io"
)

func Dial(network, address, codec string) (*Client, error) {
	if newTransport := NewTransport(network); newTransport != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return NewClient().Dial(newTransport(), address, func(conn io.ReadWriteCloser) ClientCodec {
				return NewClientCodec(conn, newCodec(), nil, nil)
			})
		}
		return nil, errors.New("unsupported codec: " + codec)
	}
	return nil, errors.New("unsupported protocol scheme: " + network)
}

func DialWithOptions(address string, opts *Options) (*Client, error) {
	if opts.NewCodec == nil && opts.NewEncoder == nil {
		return nil, errors.New("need opts.Codec or opts.Encoder")
	}
	return NewClient().Dial(opts.NewTransport(), address, func(conn io.ReadWriteCloser) ClientCodec {
		var bodyCodec codec.Codec
		if opts.NewCodec != nil {
			bodyCodec = opts.NewCodec()
		}
		var headerEncoder *encoder.Encoder
		if opts.NewEncoder != nil {
			headerEncoder = opts.NewEncoder()
		}
		var stream Stream
		if opts.NewStream != nil {
			stream = opts.NewStream()
		}
		return NewClientCodec(conn, bodyCodec, headerEncoder, stream)
	})

}
