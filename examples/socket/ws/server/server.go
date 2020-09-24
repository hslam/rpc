package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/socket/ws/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.Listen("ws", ":9999", "pb")
}
