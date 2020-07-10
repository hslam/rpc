package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/socket"
)

func Listen(network, address string, codec string) error {
	if newSocket := NewSocket(network); newSocket != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return DefaultServer.Listen(newSocket(), address, func(messages socket.Messages) ServerCodec {
				return NewServerCodec(newCodec(), nil, messages)
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
	if opts.NewSocket == nil && opts.NewMessages == nil {
		return errors.New("need opts.NewSocket or opts.NewMessages")
	}
	var bodyCodec codec.Codec
	if opts.NewCodec != nil {
		bodyCodec = opts.NewCodec()
	}
	var headerEncoder *encoder.Encoder
	if opts.NewEncoder != nil {
		headerEncoder = opts.NewEncoder()
	}
	if opts.NewMessages != nil {
		if messages := opts.NewMessages(); messages == nil {
			return errors.New("NewMessages failed")
		} else {
			if codec := NewServerCodec(opts.NewCodec(), opts.NewEncoder(), messages); codec == nil {
				return errors.New("NewClientCodec failed")
			} else {
				DefaultServer.ServeCodec(codec)
				return nil
			}
		}
	}
	return DefaultServer.Listen(opts.NewSocket(), address, func(messages socket.Messages) ServerCodec {
		return NewServerCodec(bodyCodec, headerEncoder, messages)
	})
}
