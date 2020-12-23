package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/benchmarks/reverseproxy/service"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.Parse()
}

func main() {
	rpc.Register(new(service.Arith))
	rpc.ListenWithOptions(addr, &rpc.Options{Network: "tcp", Codec: "pb"})
}
