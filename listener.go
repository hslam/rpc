// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

func Listen(network, address string, codec string) error {
	return DefaultServer.Listen(network, address, codec)
}

func ListenWithOptions(address string, opts *Options) error {
	return DefaultServer.ListenWithOptions(address, opts)

}
