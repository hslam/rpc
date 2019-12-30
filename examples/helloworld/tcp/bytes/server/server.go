package main
import (
	service "github.com/hslam/rpc/examples/service/bytes"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Echo))
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http
}
