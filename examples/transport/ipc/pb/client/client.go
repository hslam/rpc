package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
	"fmt"
	"log"
	"time"
)
func main()  {
	MaxConnsPerHost:=2
	MaxIdleConnsPerHost:=0
	transport:=rpc.NewTransport(MaxConnsPerHost,MaxIdleConnsPerHost,"ipc","pb",rpc.DefaultOptions())
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	var err error
	err = transport.Call("Arith.Multiply", req, &res,"/tmp/ipc1")
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	err = transport.Call("Arith.Divide", req, &res,"/tmp/ipc")
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
	time.Sleep(time.Hour)
}