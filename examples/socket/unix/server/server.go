package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/socket/unix/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.Listen("unix", "/tmp/unix", "pb")
}
