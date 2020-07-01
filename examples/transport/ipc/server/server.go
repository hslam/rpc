package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/transport/ipc/service"
)

var network string
var addr string
var codec string

func init() {
	flag.StringVar(&network, "network", "ipc", "-network=tcp")
	flag.StringVar(&addr, "addr", "/tmp/ipc", "-addr=:9999")
	flag.StringVar(&codec, "codec", "pb", "-codec=code")
	flag.Parse()
}
func main() {
	rpc.Register(new(service.Arith))
	rpc.Listen(network, addr, codec)
}
