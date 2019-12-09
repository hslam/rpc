package main
import (
	"github.com/hslam/rpc/examples/helloworld/json_js/service"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("http",":8080")
}
