package main

import (
	"fmt"
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/pb"
	"log"
)

func main() {
	opts := rpc.DefaultOptions()
	opts.SetPipelining(true)
	conn, err := rpc.DialWithOptions("tcp", "127.0.0.1:8080", "pb", opts) //tcp|ws|quic|http
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	defer conn.Close()
	req := &service.SeqResponse{}
	conn.CallNoRequest("Seq.Reset", req)
	if req.B != 1 {
		panic("reset error")
	}
	total := int32(10000)
	ch := make(chan *rpc.Call, total)
	for i := int32(0); i < total; i++ {
		ch <- conn.Go("Seq.Check", &service.SeqRequest{A: i}, &service.SeqResponse{}, make(chan *rpc.Call, 1))
	}
	for i := int32(0); i < total; i++ {
		call := <-ch
		<-call.Done
		if call.Error != nil {
			log.Fatalln("Seq error: ", call.Error)
		}
		fmt.Printf("req:%d res:%d\n", call.Args, call.Reply)
	}
}
