package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/reverseproxy/service"
)

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
	opts := &rpc.Options{Network: "tcp", Codec: "pb"}
	arith := new(Arith)
	arith.ReverseProxy = rpc.NewSingleHostReverseProxy(":9999")
	arith.ReverseProxy.Transport = &rpc.Transport{Options: opts}
	rpc.Register(arith)
	rpc.ListenWithOptions(":8888", opts)
}
