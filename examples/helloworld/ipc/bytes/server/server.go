package main

import (
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/bytes"
)

func main() {
	rpc.Register(new(service.Echo))
	rpc.ListenAndServe("ipc", "/tmp/ipc")
}
