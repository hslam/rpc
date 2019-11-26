package main
import (
	"hslam.com/git/x/rpc/examples/transport/pb/service"
	"hslam.com/git/x/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	go func() {
		rpc.ListenAndServe("tcp",":8081")//tcp|ws|quic|http
	}()
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http
}
