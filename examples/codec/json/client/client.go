package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/codec/json/service"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "json")
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
