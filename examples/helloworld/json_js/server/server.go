package main

import (
	"hslam.com/mgit/Mort/rpc/examples/helloworld/json/service"
	"hslam.com/mgit/Mort/rpc"
	"strconv"
	"flag"
)

var network string
var port int
var saddr string
func init()  {
	flag.StringVar(&network, "network", "http", "network: -network=tcp|ws|fast|http|http2|quic|udp")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe(network,saddr)
}
