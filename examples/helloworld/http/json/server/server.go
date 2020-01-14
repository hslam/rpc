package main

import (
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/json"
)

func main() {
	rpc.SETRPCCODEC(rpc.RPC_CODEC_PROTOBUF)
	rpc.Register(new(service.Arith))
	rpc.SetPipelining(true)
	rpc.ListenAndServe("http", ":8080")
}
