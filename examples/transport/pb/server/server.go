package main
import (
	"github.com/hslam/rpc/examples/transport/pb/service"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	go func() {
		rpc.ListenAndServe("tcp",":8081")//tcp|ws|quic|http
	}()
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http
}
