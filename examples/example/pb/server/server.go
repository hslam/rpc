package main

import (
	"hslam.com/mgit/Mort/rpc"
	"hslam.com/mgit/Mort/rpc/examples/example/pb/service"
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
var num int
var max int
var useWorkerPool bool
var saddr string
func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "udp", "network: -network=fast;ws;tcp;quic;udp")
	flag.BoolVar(&debug, "debug", true, "debug: -debug=false")
	flag.IntVar(&debug_port, "dp", 6060, "debug_port: -dp=6060")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.IntVar(&num, "n", 65536, "num: -n=1024")
	flag.IntVar(&max, "m", runtime.NumCPU()*64, "max: -m=64")
	flag.BoolVar(&useWorkerPool, "pool", false, "useWorkerPool: -pool=false")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	go func() {if debug{log.Println(http.ListenAndServe(":"+strconv.Itoa(debug_port), nil))}}()
	rpc.Register(new(service.Arith))
	rpc.SetLogLevel(6)
	if useWorkerPool{
		rpc.EnabledWorkerPoolWithSize(num,max)
	}
	rpc.ListenAndServe(network,saddr)
}
