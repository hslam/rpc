package main

import (
	"context"
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/context/service"
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
	emptyCtx := context.Background()
	valueCtx := context.WithValue(emptyCtx, rpc.BufferContextKey, make([]byte, 64))
	ctx, cancel := context.WithTimeout(valueCtx, time.Minute)
	defer cancel()
	err = conn.CallWithContext(ctx, "Arith.Multiply", req, &res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
