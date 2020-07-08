package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/socket"
)

func Dial(network, address, codec string) (*Client, error) {
	if newSocket := NewSocket(network); newSocket != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return NewClient().Dial(newSocket(), address, func(message socket.Message) ClientCodec {
				return NewClientCodec(newCodec(), nil, message)
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
	if opts.NewSocket == nil && opts.NewMessage == nil {
		return nil, errors.New("need opts.NewSocket or opts.NewMessage")
	}
	var bodyCodec codec.Codec
	if opts.NewCodec != nil {
		bodyCodec = opts.NewCodec()
	}
	var headerEncoder *encoder.Encoder
	if opts.NewEncoder != nil {
		headerEncoder = opts.NewEncoder()
	}
	if opts.NewMessage != nil {
		if message := opts.NewMessage(); message == nil {
			return nil, errors.New("NewMessage failed")
		} else {
			if codec := NewClientCodec(bodyCodec, headerEncoder, opts.NewMessage()); codec == nil {
				return nil, errors.New("NewClientCodec failed")
			} else {
				return NewClientWithCodec(codec), nil
			}
		}
	}
	return NewClient().Dial(opts.NewSocket(), address, func(message socket.Message) ClientCodec {
		return NewClientCodec(bodyCodec, headerEncoder, message)
	})

}
