package main
import (
	service "github.com/hslam/rpc/examples/service/json"
	"github.com/hslam/rpc"
	"log"
	"fmt"
)
func main()  {
	rpc.SETRPCCODEC(rpc.RPC_CODEC_PROTOBUF)
	opts:=rpc.DefaultOptions()
	opts.SetPipelining(true)
	conn, err:= rpc.DialWithOptions("http","127.0.0.1:8080","json",opts)
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	defer conn.Close()
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	err = conn.Call("Arith.Multiply", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)

	err = conn.Call("Arith.Divide", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
}