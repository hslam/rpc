package rpc

import (
	"github.com/hslam/transport"
	"github.com/hslam/transport/http"
	"github.com/hslam/transport/ipc"
	"github.com/hslam/transport/tcp"
	"sync"
)

var transports = sync.Map{}

func init() {
	RegisterTransport("tcp", tcp.NewTransport)
	RegisterTransport("ipc", ipc.NewTransport)
	RegisterTransport("http", http.NewTransport)
}

func RegisterTransport(name string, New func() transport.Transport) {
	transports.Store(name, New)
}

func NewTransport(name string) func() transport.Transport {
	if tran, ok := transports.Load(name); ok {
		return tran.(func() transport.Transport)
	}
	return nil
}
