package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/pipelining/service"
)

func main() {
	rpc.Register(new(service.Seq))
	rpc.SetPipelining(true)
	rpc.Listen("tcp", ":9999", "json")
}
