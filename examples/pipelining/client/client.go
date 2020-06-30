package main

import (
	"flag"
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/codec/pb"
	"github.com/hslam/rpc/examples/pipelining/service"
	"github.com/hslam/transport/tcp"
	"log"
)

var network string
var addr string

func init() {
	flag.StringVar(&network, "network", "tcp", "-network=tcp")
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.Parse()
}

func main() {
	conn, err := rpc.Dial(tcp.NewTransport(), addr, pb.NewCodec())
	if err != nil {
		log.Fatalln("dailing error: ", err)
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
			log.Fatalln("Seq error: ", call.Error)
		}
		fmt.Printf("req:%d res:%d\n", call.Args, call.Reply)
	}
}
