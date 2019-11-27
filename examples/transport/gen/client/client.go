package main
import (
	"hslam.com/git/x/rpc/examples/transport/gen/service"
	"hslam.com/git/x/rpc"
	"fmt"
	"log"
	"time"
)
func main()  {
	MaxConnsPerHost:=2
	MaxIdleConnsPerHost:=0
	transport:=rpc.NewTransport(MaxConnsPerHost,MaxIdleConnsPerHost,"tcp","gen",rpc.DefaultOptions())//tcp|ws|quic|http
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	var err error
	err = transport.Call("Arith.Multiply", req, &res,"127.0.0.1:8081")
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	err = transport.Call("Arith.Divide", req, &res,"127.0.0.1:8080")
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
	time.Sleep(time.Hour)
}