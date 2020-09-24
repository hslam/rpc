package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/reverseproxy/service"
)

func main() {
	opts := &rpc.Options{Network: "tcp", Codec: "pb"}
	conn, err := rpc.DialWithOptions(":8888", opts)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	err = conn.Call("Arith.Multiply", req, &res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
