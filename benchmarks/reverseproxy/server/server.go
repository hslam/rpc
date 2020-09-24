package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/benchmarks/reverseproxy/service"
	"github.com/hslam/socket"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.Parse()
}

func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenWithOptions(addr, &rpc.Options{NewSocket: socket.NewTCPSocket, NewCodec: rpc.NewPBCodec, NewHeaderEncoder: rpc.NewPBEncoder})
}
