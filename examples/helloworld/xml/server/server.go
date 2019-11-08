package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/xml/service"
	"hslam.com/git/x/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("tcp",":9999")//tcp|ws|http|http2|quic|udp
}
