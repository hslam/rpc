package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("ipc","/tmp/ipc")
}
