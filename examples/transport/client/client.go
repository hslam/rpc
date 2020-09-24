package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/transport/service"
)

func main() {
	trans := &rpc.Transport{
		MaxConnsPerHost:     1,
		MaxIdleConnsPerHost: 1,
		Options:             &rpc.Options{Network: "tcp", Codec: "pb"},
	}
	defer trans.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	var err error
	err = trans.Call(":9999", "Arith.Multiply", req, &res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
