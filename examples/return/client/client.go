package main

import (
	"context"
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/return/service"
	"time"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "pb")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = conn.CallWithContext(ctx, "Arith.Multiply", req, &res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
