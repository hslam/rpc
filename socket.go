// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/socket"
	"github.com/hslam/socket/http"
	"github.com/hslam/socket/tcp"
	"github.com/hslam/socket/unix"
	"github.com/hslam/socket/ws"
	"sync"
)

var sockets = sync.Map{}

func init() {
	RegisterSocket("http", http.NewSocket)
	RegisterSocket("tcp", tcp.NewSocket)
	RegisterSocket("unix", unix.NewSocket)
	RegisterSocket("ws", ws.NewSocket)
}

func RegisterSocket(network string, New func() socket.Socket) {
	sockets.Store(network, New)
}

func NewSocket(network string) func() socket.Socket {
	if s, ok := sockets.Load(network); ok {
		return s.(func() socket.Socket)
	}
	return nil
}
