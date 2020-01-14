package main

import (
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/json"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("ipc", "/tmp/ipc")
}
