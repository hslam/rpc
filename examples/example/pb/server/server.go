package main

import (
	"hslam.com/git/x/rpc"
	"hslam.com/git/x/rpc/examples/example/pb/service"
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
var saddr string
func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|fast|http|http2|quic|udp")
	flag.BoolVar(&debug, "debug", true, "debug: -debug=false")
	flag.IntVar(&debug_port, "dp", 6060, "debug_port: -dp=6060")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.BoolVar(&async, "async", false, "async: -async=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "multiplexing: -multiplexing=false")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	go func() {if debug{log.Println(http.ListenAndServe(":"+strconv.Itoa(debug_port), nil))}}()
	rpc.Register(new(service.Arith))
	rpc.SetLogLevel(4)
	if async{
		rpc.EnableAsyncHandle()
	}
	if multiplexing{
		rpc.EnableMultiplexing()
	}
	rpc.ListenAndServe(network,saddr)
}
