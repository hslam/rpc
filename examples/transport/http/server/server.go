package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/codec/pb"
	"github.com/hslam/rpc/examples/transport/http/service"
	"github.com/hslam/transport/http"
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
	rpc.ListenAndServe(http.NewTransport(), addr, pb.NewCodec())
}
