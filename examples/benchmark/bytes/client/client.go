package main

import (
	"hslam.com/git/x/rpc"
	"hslam.com/git/x/stats"
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
var pipelining bool
var multiplexing bool
var noresponse bool
var norequest bool
var onlycall bool
var clients int
var total_calls int
var bar bool
var batch_async bool

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|http|http2|quic|udp")
	flag.StringVar(&codec, "codec", "bytes", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&compress, "compress", "no", "compress: -compress=no|flate|zlib|gzip")
	flag.StringVar(&host, "h", "127.0.0.1", "host: -h=127.0.0.1")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.IntVar(&total_calls, "total", 1000000, "total_calls: -total=10000")
	flag.BoolVar(&batch, "batch", true, "batch: -batch=false")
	flag.BoolVar(&batch_async, "batch_async", true, "batch_async: -batch_async=false")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "pipelining: -pipelining=false")
	flag.BoolVar(&norequest, "norequest", false, "norequest: -norequest=false")
	flag.BoolVar(&noresponse, "noresponse", false, "noresponse: -noresponse=false")
	flag.BoolVar(&onlycall, "onlycall", false, "onlycall: -onlycall=false")
	flag.IntVar(&clients, "clients", 1, "num: -clients=1")
	flag.BoolVar(&bar, "bar", true, "bar: -bar=true")
	log.SetFlags(0)
	flag.Parse()
	addr=host+":"+strconv.Itoa(port)
	stats.SetLog(bar)
}

func main()  {
	fmt.Printf("./client -network=%s -codec=%s -compress=%s -h=%s -p=%d -total=%d -pipelining=%t -batch=%t -batch_async=%t -norequest=%t -noresponse=%t -onlycall=%t -clients=%d\n",network,codec,compress,host,port,total_calls,pipelining,batch,batch_async,norequest,noresponse,onlycall,clients)
	var wrkClients []stats.Client
	parallel:=1
	if clients>1{
		pool,err := rpc.Dials(clients,network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		pool.SetCompressType(compress)
		if batch {pool.EnableBatch()}
		if batch_async{pool.EnableBatchAsync()}
		if multiplexing{pool.EnableMultiplexing()}
		var req =[]byte("Hello World")
		for i:=0; i<len(pool.All());i++  {
			pool.All()[i].CallNoResponse("Echo.Set", &req)
		}
		wrkClients=make([]stats.Client,len(pool.All()))
		for i:=0; i<len(pool.All());i++  {
			wrkClients[i]=&WrkClient{pool.All()[i]}
		}
		if batch{
			parallel=pool.GetMaxBatchRequest()
		}else if pipelining{
			parallel=pool.GetMaxRequests()
		}else if multiplexing{
			parallel=pool.GetMaxRequests()
		}
	}else if clients==1 {
		conn, err:= rpc.Dial(network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		conn.SetCompressType(compress)
		if batch {conn.EnableBatch()}
		if batch_async{conn.EnableBatchAsync()}
		if multiplexing{conn.EnableMultiplexing()}
		if batch{
			parallel=conn.GetMaxBatchRequest()
		}else if pipelining{
			parallel=conn.GetMaxRequests()
		}else if multiplexing{
			parallel=conn.GetMaxRequests()
		}
		wrkClients=make([]stats.Client,1)
		wrkClients[0]= &WrkClient{conn}
		var req =[]byte("Hello World")
		conn.CallNoResponse("Echo.Set", &req)
	}else {
		return
	}
	stats.StartPrint(parallel,total_calls,wrkClients)
}

type WrkClient struct {
	Conn rpc.Client
}

func (c *WrkClient)Call()(int64,int64,bool){
	var err error
	length:=10
	req := make([]byte, length)
	for i := 0; i < length; i++ {
		b := rand.Intn(26) + 65
		req[i] = byte(b)
	}
	if noresponse{
		err = c.Conn.CallNoResponse("Echo.Set", &req)
		if err==nil{
			return int64(length),0,true
		}
	}else if onlycall{
		err = c.Conn.OnlyCall("Echo.Clear")
		if err==nil{
			return 0,0,true
		}
	}else if norequest{
		var req =[]byte("Hello World")
		var res []byte
		err = c.Conn.CallNoRequest("Echo.Get", &res)
		if bytes.Equal(res,[]byte("Hello World")){
			return 0,int64(len(res)),true
		}else {
			fmt.Printf("Echo.Get is not equal: req-%s res-%s\n",string(req),string(res))
		}
	}else {
		var res []byte
		err = c.Conn.Call("Echo.ToLower", &req, &res)
		if bytes.Equal(res,[]byte(strings.ToLower(string(req)))){
			return int64(len(req)),int64(len(res)),true
		}else {
			fmt.Printf("Echo.ToLower is not equal: req-%s res-%s\n",string(req),string(res))
		}
		if err != nil {
			fmt.Println("Echo error: ", err)
		}
	}
	return 0,0,false
}