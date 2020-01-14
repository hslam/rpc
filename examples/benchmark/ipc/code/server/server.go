package main

import (
	"flag"
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/code"
	"runtime"
)

var network string
var addr string
var multiplexing bool
var pipelining bool
var batching bool

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "ipc", "network: -network=ipc")
	flag.StringVar(&addr, "addr", "/tmp/ipc", "addr: -addr=/tmp/ipc")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "multiplexing: -multiplexing=false")
	flag.BoolVar(&batching, "batching", false, "batching: -batching=false")
	flag.Parse()
}
func main() {
	rpc.Register(new(service.Arith))
	rpc.SetLogLevel(rpc.InfoLevel)
	rpc.SetPipelining(pipelining)
	rpc.SetMultiplexing(multiplexing)
	rpc.SetBatching(batching)
	rpc.ListenAndServe(network, addr)
}
