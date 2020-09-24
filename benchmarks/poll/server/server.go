package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/benchmarks/poll/service"
)

var network string
var addr string
var codec string

func init() {
	flag.StringVar(&network, "network", "tcp", "-network=tcp")
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.StringVar(&codec, "codec", "pb", "-codec=code")
	flag.Parse()
}

func main() {
	rpc.Register(new(service.Arith))
	rpc.SetPoll(true)
	rpc.Listen(network, addr, codec)
}
