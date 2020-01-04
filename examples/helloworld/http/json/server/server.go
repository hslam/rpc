package main
import (
	service "github.com/hslam/rpc/examples/service/json"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.SETRPCCODEC(rpc.RPC_CODEC_PROTOBUF)
	rpc.Register(new(service.Arith))
	rpc.SetPipelining(true)
	rpc.ListenAndServe("http",":8080")
}
