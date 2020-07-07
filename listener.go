package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"io"
)

func Listen(network, address string, codec string) error {
	if newTransport := NewTransport(network); newTransport != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return DefaultServer.Listen(newTransport(), address, func(conn io.ReadWriteCloser) ServerCodec {
				return NewServerCodec(conn, newCodec(), nil, nil)
			})
		}
		return errors.New("unsupported codec: " + codec)
	}
	return errors.New("unsupported protocol scheme: " + network)
}

func ListenWithOptions(address string, opts *Options) error {
	if opts.NewCodec == nil && opts.NewEncoder == nil {
		return errors.New("need opts.Codec or opts.Encoder")
	}
	return DefaultServer.Listen(opts.NewTransport(), address, func(conn io.ReadWriteCloser) ServerCodec {
		var bodyCodec codec.Codec
		if opts.NewCodec != nil {
			bodyCodec = opts.NewCodec()
		}
		var headerEncoder *encoder.Encoder
		if opts.NewEncoder != nil {
			headerEncoder = opts.NewEncoder()
		}
		var messageConn MessageConn
		if opts.NewMessageConn != nil {
			messageConn = opts.NewMessageConn()
		}
		return NewServerCodec(conn, bodyCodec, headerEncoder, messageConn)
	})
}
