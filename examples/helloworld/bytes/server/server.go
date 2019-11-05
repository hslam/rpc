package main

import (
	"hslam.com/git/x/rpc/examples/helloworld/bytes/service"
	"hslam.com/git/x/rpc"
	"strconv"
	"flag"
)

var network string
var port int
var saddr string
func init()  {
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|http|http2|quic|udp")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	rpc.Register(new(service.Echo))
	rpc.ListenAndServe(network,saddr)
}
