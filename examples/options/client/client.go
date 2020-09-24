package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/options/service"
	"github.com/hslam/socket"
)

func main() {
	conn, err := rpc.DialWithOptions(":9999", &rpc.Options{
		NewSocket:        socket.NewTCPSocket,
		NewCodec:         rpc.NewPBCodec,
		NewHeaderEncoder: rpc.NewPBEncoder,
	})
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
