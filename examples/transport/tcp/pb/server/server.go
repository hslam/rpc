package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	go rpc.ListenAndServe("tcp",":8081")
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http
}
