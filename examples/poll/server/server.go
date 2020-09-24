package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/poll/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.SetPoll(true)
	rpc.Listen("tcp", ":9999", "pb")
}
