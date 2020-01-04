package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
	"fmt"
	"log"
)
func main()  {
	opts:=rpc.DefaultOptions()
	opts.SetMultiplexing(true)
	conn, err:= rpc.DialWithOptions("tcp","127.0.0.1:8080","pb",opts)//tcp|ws|quic|http
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