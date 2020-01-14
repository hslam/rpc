package main

import (
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/pb"
)

func main() {
	rpc.Register(new(service.Seq))
	rpc.SetPipelining(true)
	rpc.ListenAndServe("tcp", ":8080") //tcp|ws|quic|http
}
