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

func RegisterTransport(name string, new transport.NewTransport) {
	transports.Store(name, new)
}

func NewTransport(name string) transport.Transport {
	if tran, ok := transports.Load(name); ok {
		return tran.(transport.NewTransport)()
	}
	return nil
}
