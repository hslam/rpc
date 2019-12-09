package main
import (
	"github.com/hslam/rpc/examples/helloworld/json/service"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http
}
