package main

import (
	"flag"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/benchmarks/reverseproxy/service"
	"github.com/hslam/socket"
)

var addr string
var target string

func init() {
	flag.StringVar(&addr, "addr", ":8888", "-addr=:8888")
	flag.StringVar(&target, "target", ":9999", "-addr=:9999")
	flag.Parse()
}

//Arith defines the struct of arith.
type Arith struct {
	ReverseProxy *rpc.ReverseProxy
}

//Multiply operation
func (a *Arith) Multiply(req *service.ArithRequest, res *service.ArithResponse) error {
	a.ReverseProxy.Call("Arith.Multiply", req, res)
	return nil
}

func main() {
	opts := &rpc.Options{NewSocket: socket.NewTCPSocket, NewCodec: rpc.NewPBCodec, NewHeaderEncoder: rpc.NewPBEncoder}
	arith := new(Arith)
	arith.ReverseProxy = rpc.NewSingleHostReverseProxy(target)
	arith.ReverseProxy.Transport = &rpc.Transport{Options: opts}
	rpc.Register(arith)
	rpc.ListenWithOptions(addr, opts)
}
