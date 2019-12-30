package main

import (
	"github.com/hslam/rpc"
	service "github.com/hslam/rpc/examples/service/pb"
	_ "net/http/pprof"
	"net/http"
	"strconv"
	"runtime"
	"flag"
	"log"
)

var network string
var debug bool
var debug_port int
var addr string
var async bool
var multiplexing bool
var pipelining bool
var batch bool
func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "ipc", "network: -network=ipc")
	flag.BoolVar(&debug, "debug", true, "debug: -debug=false")
	flag.IntVar(&debug_port, "dp", 6060, "debug_port: -dp=6060")
	flag.StringVar(&addr, "address", "/tmp/ipc", "address: -address=/tmp/ipc")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&async, "async", false, "async: -async=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "multiplexing: -multiplexing=false")
	flag.BoolVar(&batch, "batch", false, "batch: -batch=false")
	flag.Parse()
}
func main()  {
	go func() {if debug{log.Println(http.ListenAndServe(":"+strconv.Itoa(debug_port), nil))}}()
	rpc.Register(new(service.Arith))
	rpc.SetLogLevel(rpc.InfoLevel)
	rpc.SetPipelining(pipelining)
	rpc.SetPipeliningAsync(async)
	rpc.SetMultiplexing(multiplexing)
	rpc.SetBatch(batch)
	rpc.ListenAndServe(network,addr)
}
