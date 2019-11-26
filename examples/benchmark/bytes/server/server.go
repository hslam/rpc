package main

import (
	"hslam.com/git/x/rpc"
	"hslam.com/git/x/rpc/examples/benchmark/bytes/service"
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
var port int
var async bool
var multiplexing bool
var pipelining bool
var batch bool
var saddr string
func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|quic|http")
	flag.BoolVar(&debug, "debug", true, "debug: -debug=false")
	flag.IntVar(&debug_port, "dp", 6060, "debug_port: -dp=6060")
	flag.IntVar(&port, "p", 8080, "port: -p=8080")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&async, "async", false, "async: -async=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "multiplexing: -multiplexing=false")
	flag.BoolVar(&batch, "batch", true, "batch: -batch=false")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	go func() {if debug{log.Println(http.ListenAndServe(":"+strconv.Itoa(debug_port), nil))}}()
	rpc.Register(new(service.Echo))
	rpc.SetLogLevel(rpc.InfoLevel)
	rpc.SetPipelining(pipelining)
	rpc.SetPipeliningAsync(async)
	rpc.SetMultiplexing(multiplexing)
	rpc.SetBatch(batch)
	rpc.ListenAndServe(network,saddr)
}
