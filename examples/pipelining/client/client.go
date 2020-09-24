package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/pipelining/service"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "pb")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	req := &service.Message{}
	res := &service.Message{}
	conn.Call("Seq.Reset", req, res)
	if res.Value != 1 {
		panic("reset error")
	}
	total := int32(10000)
	ch := make(chan *rpc.Call, total)
	for i := int32(0); i < total; i++ {
		ch <- conn.Go("Seq.Check", &service.Message{Value: i}, &service.Message{}, make(chan *rpc.Call, 1))
	}
	for i := int32(0); i < total; i++ {
		call := <-ch
		<-call.Done
		if call.Error != nil {
			panic(call.Error)
		}
		fmt.Printf("req:%d res:%d\n", call.Args, call.Reply)
	}
}
