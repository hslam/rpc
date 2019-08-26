package main

import (
	"hslam.com/mgit/Mort/rpc"
	"hslam.com/mgit/Mort/rpc/stats"
	"math/rand"
	"strconv"
	"runtime"
	"flag"
	"log"
	"fmt"
	"bytes"
	"strings"
)
var network string
var codec string
var compress string
var host string
var port int
var addr string
var batch bool
var concurrent bool
var noresponse bool
var clients int
var total_calls int
var bar bool

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|fast|http|http2|quic|udp")
	flag.StringVar(&codec, "codec", "bytes", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&compress, "compress", "no", "compress: -compress=no|flate|zlib|gzip")
	flag.StringVar(&host, "h", "127.0.0.1", "host: -h=127.0.0.1")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.IntVar(&total_calls, "total", 1000000, "total_calls: -total=10000")
	flag.BoolVar(&batch, "batch", true, "batch: -batch=false")
	flag.BoolVar(&concurrent, "concurrent", false, "concurrent: -concurrent=false")
	flag.BoolVar(&noresponse, "noresponse", false, "noresponse: -noresponse=false")
	flag.IntVar(&clients, "clients", 1, "num: -clients=1")
	flag.BoolVar(&bar, "bar", true, "bar: -bar=true")
	log.SetFlags(0)
	flag.Parse()
	addr=host+":"+strconv.Itoa(port)
	stats.SetLog(bar)
}

func main()  {
	fmt.Printf("./client -network=%s -codec=%s -h=%s -p=%d -total=%d -concurrent=%t -batch=%t -noresponse=%t -clients=%d\n",network,codec,host,port,total_calls,concurrent,batch,noresponse,clients)
	var wrkClients []stats.Client
	parallel:=1
	if clients>1{
		pool,err := rpc.Dials(clients,network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		pool.SetCompressType(compress)
		if batch {pool.EnabledBatch()}

		wrkClients=make([]stats.Client,len(pool.All()))
		for i:=0; i<len(pool.All());i++  {
			wrkClients[i]=&WrkClient{pool.All()[i]}
		}
		if batch{
			parallel=pool.All()[0].GetMaxBatchRequest()
		}else if concurrent{
			parallel=pool.All()[0].GetMaxConcurrentRequest()
		}
	}else if clients==1 {
		conn, err:= rpc.Dial(network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		conn.SetCompressType(compress)
		if batch {conn.EnabledBatch()}
		if batch{
			parallel=conn.GetMaxBatchRequest()
		}else if concurrent{
			parallel=conn.GetMaxConcurrentRequest()
		}
		wrkClients=make([]stats.Client,1)
		wrkClients[0]= &WrkClient{conn}
	}else {
		return
	}
	stats.StartClientStats(parallel,total_calls,wrkClients)
}

type WrkClient struct {
	Conn *rpc.Client
}

func (c *WrkClient)Call()(int64,bool){
	var err error
	len:=10
	req := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		req[i] = byte(b)
	}
	if noresponse{
		err = c.Conn.CallNoResponse("Echo.ToLower", &req)
		if err==nil{
			return 0,true
		}
	}else {
		var res []byte
		err = c.Conn.Call("Echo.ToLower", &req, &res)
		if bytes.Equal(res,[]byte(strings.ToLower(string(req)))){
			return 0,true
		}else {
			fmt.Printf("Echo.ToLower is not equal: req-%s res-%s\n",string(req),string(res))
		}
		if err != nil {
			fmt.Println("Echo error: ", err)
		}
	}
	return 0,false
}