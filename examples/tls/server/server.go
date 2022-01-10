package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/tls/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenTLS("tcp", ":9999", "pb", rpc.DefalutServerTLSConfig())
}
