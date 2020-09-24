// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"crypto/tls"
	"github.com/hslam/socket"
	"sync"
)

var sockets = sync.Map{}

func init() {
	RegisterSocket("http", socket.NewHTTPSocket)
	RegisterSocket("tcp", socket.NewTCPSocket)
	RegisterSocket("unix", socket.NewUNIXSocket)
	RegisterSocket("ws", socket.NewWSSocket)
}

// RegisterSocket registers a network socket.
func RegisterSocket(network string, New func(config *tls.Config) socket.Socket) {
	sockets.Store(network, New)
}

// NewSocket returns a new Socket by network.
func NewSocket(network string) func(config *tls.Config) socket.Socket {
	if s, ok := sockets.Load(network); ok {
		return s.(func(config *tls.Config) socket.Socket)
	}
	return nil
}
