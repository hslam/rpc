package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/codec/pb/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.Listen("tcp", ":9999", "pb")
}
