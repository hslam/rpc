package main
import (
	"hslam.com/git/x/rpc/examples/transport/pb/service"
	"hslam.com/git/x/rpc"
	"fmt"
	"log"
)
func main()  {
	maxConnsPerHost:=1
	transport:=rpc.NewTransport(maxConnsPerHost,"tcp","pb",rpc.DefaultOptions())
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	var err error
	err = transport.Call("Arith.Multiply", req, &res,"127.0.0.1:9998")
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	err = transport.Call("Arith.Divide", req, &res,"127.0.0.1:9999")
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
}