package main

import (
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/pb"
)

func main() {
	rpc.Register(new(service.Arith))
	go rpc.ListenAndServe("ipc", "/tmp/ipc1")
	rpc.ListenAndServe("ipc", "/tmp/ipc")
}
