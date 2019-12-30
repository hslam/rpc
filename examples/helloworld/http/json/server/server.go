package main
import (
	service "github.com/hslam/rpc/examples/service/json"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.SetMultiplexing(false)
	rpc.SetPipelining(true)
	rpc.ListenAndServe("http",":8080")
}
