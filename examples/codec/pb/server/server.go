package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/codec/pb"
	"github.com/hslam/rpc/examples/codec/pb/service"
	"github.com/hslam/transport/tcp"
)

var network string
var addr string

func init() {
	flag.StringVar(&network, "network", "tcp", "-network=tcp")
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.Parse()
}
func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe(tcp.NewTransport(), addr, pb.NewCodec())
}
