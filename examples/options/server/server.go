package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/options/service"
	"github.com/hslam/socket"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenWithOptions(":9999", &rpc.Options{
		NewSocket:        socket.NewTCPSocket,
		NewCodec:         rpc.NewPBCodec,
		NewHeaderEncoder: rpc.NewPBEncoder,
	})
}
