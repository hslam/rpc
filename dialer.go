// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"crypto/tls"
	"errors"
	"github.com/hslam/socket"
)

// Dial connects to an RPC server at the specified network address.
func Dial(network, address, codec string) (*Conn, error) {
	if newSocket := NewSocket(network); newSocket != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return NewConn().Dial(newSocket(nil), address, func(messages socket.Messages) ClientCodec {
				return NewClientCodec(newCodec(), nil, messages, bufferSize)
			})
		}
		return nil, errors.New("unsupported codec: " + codec)
	}
	return nil, errors.New("unsupported protocol scheme: " + network)
}

// DialTLS connects to an RPC server at the specified network address with tls.Config.
func DialTLS(network, address, codec string, config *tls.Config) (*Conn, error) {
	if newSocket := NewSocket(network); newSocket != nil {
		if newCodec := NewCodec(codec); newCodec != nil {
			return NewConn().Dial(newSocket(config), address, func(messages socket.Messages) ClientCodec {
				return NewClientCodec(newCodec(), nil, messages, bufferSize)
			})
		}
		return nil, errors.New("unsupported codec: " + codec)
	}
	return nil, errors.New("unsupported protocol scheme: " + network)
}

// DialWithOptions connects to an RPC server at the specified network address with Options.
func DialWithOptions(address string, opts *Options) (*Conn, error) {
	if opts.NewCodec == nil && opts.NewHeaderEncoder == nil && opts.Codec == "" {
		return nil, errors.New("need opts.NewCodec, opts.NewHeaderEncoder or opts.Codec")
	}
	if opts.NewSocket == nil && opts.Network == "" {
		return nil, errors.New("need opts.NewSocket or opts.Network")
	}
	var sock socket.Socket
	if newSocket := NewSocket(opts.Network); newSocket != nil {
		sock = newSocket(opts.TLSConfig)
	} else if opts.NewSocket != nil {
		sock = opts.NewSocket(opts.TLSConfig)
	}
	return NewConn().Dial(sock, address, func(messages socket.Messages) ClientCodec {
		var bodyCodec Codec
		if newCodec := NewCodec(opts.Codec); newCodec != nil {
			bodyCodec = newCodec()
		} else if opts.NewCodec != nil {
			bodyCodec = opts.NewCodec()
		}
		var headerEncoder Encoder
		if newEncoder := NewHeaderEncoder(opts.HeaderEncoder); newEncoder != nil {
			headerEncoder = newEncoder()
		} else if opts.NewHeaderEncoder != nil {
			headerEncoder = opts.NewHeaderEncoder()
		}
		return NewClientCodec(bodyCodec, headerEncoder, messages, opts.ClientBufferSize)
	})
}
