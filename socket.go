package rpc

import (
	"github.com/hslam/socket"
	"github.com/hslam/socket/http"
	"github.com/hslam/socket/ipc"
	"github.com/hslam/socket/tcp"
	"sync"
)

var sockets = sync.Map{}

func init() {
	RegisterSocket("tcp", tcp.NewSocket)
	RegisterSocket("ipc", ipc.NewSocket)
	RegisterSocket("http", http.NewSocket)
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
