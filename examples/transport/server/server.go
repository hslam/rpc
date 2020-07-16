package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/encoder/codepb"
	"github.com/hslam/rpc/encoder/pb"
	"github.com/hslam/rpc/examples/transport/service"
	"github.com/hslam/socket/tcp"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.Parse()
}
func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenWithOptions(addr, &rpc.Options{NewSocket: tcp.NewSocket, NewCodec: pb.NewCodec, NewEncoder: codepb.NewEncoder})
}
