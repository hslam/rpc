package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/benchmarks/stream/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.SetNoCopy(true)
	rpc.Listen("tcp", ":9999", "pb")
}
