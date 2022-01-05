package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/return/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.SetContextBuffer(true)
	rpc.Listen("tcp", ":9999", "pb")
}
