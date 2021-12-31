package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/stream/service"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "pb")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	stream, err := conn.NewStream("Arith.StreamMultiply")
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	for i := 0; i < 3; i++ {
		req := &service.ArithRequest{A: int32(i), B: int32(i + 1)}
		var res service.ArithResponse
		if err := stream.WriteMessage(req); err != nil {
			panic(err)
		}
		buf := make([]byte, 64)
		if err := stream.ReadMessage(buf, &res); err != nil {
			panic(err)
		}
		fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	}
}
