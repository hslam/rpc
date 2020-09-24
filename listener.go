// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"crypto/tls"
)

// Listen announces on the local network address.
func Listen(network, address string, codec string) error {
	return DefaultServer.Listen(network, address, codec)
}

// ListenTLS announces on the local network address with tls.Config.
func ListenTLS(network, address string, codec string, config *tls.Config) error {
	return DefaultServer.ListenTLS(network, address, codec, config)
}

// ListenWithOptions announces on the local network address with Options.
func ListenWithOptions(address string, opts *Options) error {
	return DefaultServer.ListenWithOptions(address, opts)

}
