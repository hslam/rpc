package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/client/client/service"
)

func main() {
	opts := &rpc.Options{Network: "tcp", Codec: "pb"}
	client := rpc.NewClient(opts, ":9997", ":9998", ":9999")
	client.Scheduling = rpc.LeastTimeScheduling
	defer client.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	if err := client.Call("Arith.Multiply", req, &res); err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
