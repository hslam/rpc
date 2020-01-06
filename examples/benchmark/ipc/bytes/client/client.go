package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/stats"
	"runtime"
	"flag"
	"log"
	"fmt"
	"bytes"
)
var network string
var codec string
var compress string
var addr string
var batching bool
var pipelining bool
var multiplexing bool
var noresponse bool
var norequest bool
var onlycall bool
var clients int
var total_calls int
var bar bool

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "ipc", "network: -network=ipc")
	flag.StringVar(&codec, "codec", "bytes", "codec: -codec=pb|json|xml|bytes|gen")
	flag.StringVar(&compress, "compress", "no", "compress: -compress=no|flate|zlib|gzip")
	flag.StringVar(&addr, "addr", "/tmp/ipc", "addr: -addr=/tmp/ipc")
	flag.IntVar(&total_calls, "total", 1000000, "total_calls: -total=10000")
	flag.BoolVar(&batching, "batching", false, "batching: -batching=false")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "pipelining: -pipelining=false")
	flag.BoolVar(&norequest, "norequest", false, "norequest: -norequest=false")
	flag.BoolVar(&noresponse, "noresponse", false, "noresponse: -noresponse=false")
	flag.BoolVar(&onlycall, "onlycall", false, "onlycall: -onlycall=false")
	flag.IntVar(&clients, "clients", 1, "num: -clients=1")
	flag.BoolVar(&bar, "bar", true, "bar: -bar=true")
	log.SetFlags(0)
	flag.Parse()
	stats.SetBar(bar)
}

func main()  {
	fmt.Printf("./client -network=%s -codec=%s -compress=%s -addr=%s -total=%d -pipelining=%t -multiplexing=%t -batching=%t -norequest=%t -noresponse=%t -onlycall=%t -clients=%d\n",network,codec,compress,addr,total_calls,pipelining,multiplexing,batching,norequest,noresponse,onlycall,clients)
	var wrkClients []stats.Client
	opts:=rpc.DefaultOptions()
	opts.SetCompressType(compress)
	opts.SetBatching(batching)
	opts.SetPipelining(pipelining)
	opts.SetMultiplexing(multiplexing)
	parallel:=1
	if clients>1{
		pool,err := rpc.DialsWithOptions(clients,network,addr,codec,opts)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		var req =[]byte("Hello World")
		for i:=0; i<len(pool.All());i++  {
			pool.All()[i].CallNoResponse("Echo.Set", &req)
		}
		wrkClients=make([]stats.Client,len(pool.All()))
		for i:=0; i<len(pool.All());i++  {
			var body=[]byte("HelloWorld")
			var result=[]byte("helloworld")
			wrkClients[i]=&WrkClient{pool.All()[i],body,result}
		}
		if batching{
			parallel=pool.GetMaxBatchRequest()
		}
		if pipelining{
			if pool.GetMaxRequests()>parallel{
				parallel=pool.GetMaxRequests()
			}
		}else if multiplexing{
			if pool.GetMaxRequests()>parallel{
				parallel=pool.GetMaxRequests()
			}
		}
	}else if clients==1 {
		conn, err:= rpc.DialWithOptions(network,addr,codec,opts)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		if batching{
			parallel=conn.GetMaxBatchRequest()
		}
		if pipelining{
			if conn.GetMaxRequests()>parallel{
				parallel=conn.GetMaxRequests()
			}
		}else if multiplexing{
			if conn.GetMaxRequests()>parallel{
				parallel=conn.GetMaxRequests()
			}
		}
		wrkClients=make([]stats.Client,1)
		var body=[]byte("HelloWorld")
		var result=[]byte("helloworld")
		wrkClients[0]= &WrkClient{conn,body,result}
		var req =[]byte("Hello World")
		conn.CallNoResponse("Echo.Set", &req)
	}else {
		return
	}
	stats.StartPrint(parallel,total_calls,wrkClients)
}

type WrkClient struct {
	Conn rpc.Client
	body []byte
	result []byte
}

func (c *WrkClient)Call()(int64,int64,bool){
	var err error
	var req =c.body
	if noresponse{
		err = c.Conn.CallNoResponse("Echo.Set", &req)
		if err==nil{
			return int64(len(c.body)),0,true
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
		if bytes.Equal(res,c.result){
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