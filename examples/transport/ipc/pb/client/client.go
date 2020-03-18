package main

import (
	"fmt"
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/pb"
	"log"
)

func main() {
	MaxConnsPerHost := 2
	MaxIdleConnsPerHost := 0
	transport := rpc.NewTransport(MaxConnsPerHost, MaxIdleConnsPerHost, "ipc", "pb", rpc.DefaultOptions())
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	var err error
	err = transport.Call("/tmp/ipc1", "Arith.Multiply", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	err = transport.Call("/tmp/ipc", "Arith.Divide", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
}
