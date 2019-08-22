package main

import (
	"hslam.com/mgit/Mort/rpc"
	"strconv"
	"flag"
	"log"
	"fmt"
)

var network string
var codec string
var host string
var port int
var addr string

func init()  {
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|fast|http|http2|quic|udp")
	flag.StringVar(&codec, "codec", "bytes", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&host, "h", "localhost", "host: -h=localhost")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Parse()
	addr=host+":"+strconv.Itoa(port)
}

func main()  {
	conn, err:= rpc.Dial(network,addr,codec)
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	defer conn.Close()
	var req =[]byte("Hello World")
	var res []byte
	err = conn.Call("Echo.ToLower", &req, &res)
	if err != nil {
		log.Fatalln("Echo error: ", err)
	}
	fmt.Printf("Echo.ToLower : %s\n", string(res))

}