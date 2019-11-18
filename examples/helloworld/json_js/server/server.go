package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/json_js/service"
	"hslam.com/git/x/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("http1",":9999")
}
