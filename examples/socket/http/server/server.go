package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/socket/http/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.Listen("http", ":9999", "pb")
}
