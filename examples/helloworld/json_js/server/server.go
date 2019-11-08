package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/json/service"
	"hslam.com/git/x/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("http",":9999")
}
