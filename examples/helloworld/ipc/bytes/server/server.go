package main
import (
	service "github.com/hslam/rpc/examples/service/bytes"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Echo))
	rpc.ListenAndServe("ipc","/tmp/ipc")
}
