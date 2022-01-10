package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/tls/service"
)

func main() {
	conn, err := rpc.DialTLS("tcp", ":9999", "pb", rpc.SkipVerifyTLSConfig())
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
