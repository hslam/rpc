// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"crypto/tls"
	"github.com/hslam/codec"
	"github.com/hslam/socket"
)

//Options defines the struct of options.
type Options struct {
	NewSocket        func(*tls.Config) socket.Socket
	NewCodec         func() codec.Codec
	NewHeaderEncoder func() *Encoder
	Network          string
	Codec            string
	HeaderEncoder    string
	TLSConfig        *tls.Config
}

//DefaultOptions returns a default options.
func DefaultOptions() *Options {
	return &Options{
		NewSocket: socket.NewTCPSocket,
		NewCodec:  NewJSONCodec,
	}
}
