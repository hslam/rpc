package main

import (
	"flag"
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/bytes"
	"runtime"
	"strconv"
)

var network string
var port int
var multiplexing bool
var pipelining bool
var batching bool
var saddr string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|quic|http")
	flag.IntVar(&port, "p", 8080, "port: -p=8080")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "multiplexing: -multiplexing=false")
	flag.BoolVar(&batching, "batching", false, "batching: -batching=false")
	flag.Parse()
	saddr = ":" + strconv.Itoa(port)
}
func main() {
	rpc.Register(new(service.Echo))
	rpc.SetLogLevel(rpc.InfoLevel)
	rpc.SetPipelining(pipelining)
	rpc.SetMultiplexing(multiplexing)
	rpc.SetBatching(batching)
	rpc.ListenAndServe(network, saddr)
}
