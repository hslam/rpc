package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/client/transport/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenWithOptions(":9999", &rpc.Options{Network: "tcp", Codec: "pb"})
}
