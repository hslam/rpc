package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/pb/service"
	"hslam.com/git/x/rpc"
	"fmt"
	"log"
)
func main()  {
	conn, err:= rpc.Dial("tcp","127.0.0.1:9999","pb")//tcp|ws|http|http2|quic
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