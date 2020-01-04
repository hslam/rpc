package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Seq))
	rpc.SetPipelining(true)
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http
}
