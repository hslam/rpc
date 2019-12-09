package main
import (
	"github.com/hslam/rpc/examples/helloworld/bytes/service"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Echo))
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http
}
