package main

import (
	"hslam.com/git/x/rpc/examples/benchmark/pb/service"
	"hslam.com/git/x/rpc"
	"hslam.com/git/x/stats"
	"math/rand"
	"strconv"
	"runtime"
	"flag"
	"log"
	"fmt"
)
var network string
var codec string
var compress string
var host string
var port int
var addr string
var batch bool
var batch_async bool
var pipelining bool
var multiplexing bool
var noresponse bool
var clients int
var total_calls int
var bar bool

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|quic|http|http1|http2")
	flag.StringVar(&codec, "codec", "pb", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&compress, "compress", "no", "compress: -compress=no|flate|zlib|gzip")
	flag.StringVar(&host, "h", "127.0.0.1", "host: -h=127.0.0.1")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.IntVar(&total_calls, "total", 100000, "total_calls: -total=10000")
	flag.BoolVar(&batch, "batch", false, "batch: -batch=false")
	flag.BoolVar(&batch_async, "batch_async", false, "batch_async: -batch_async=false")
	flag.BoolVar(&pipelining, "pipelining", true, "pipelining: -pipelining=false")
	flag.BoolVar(&multiplexing, "multiplexing", false, "pipelining: -pipelining=false")
	flag.BoolVar(&noresponse, "noresponse", false, "noresponse: -noresponse=false")
	flag.IntVar(&clients, "clients", 1, "num: -clients=1")
	flag.BoolVar(&bar, "bar", false, "bar: -bar=true")
	log.SetFlags(0)
	flag.Parse()
	addr=host+":"+strconv.Itoa(port)
	stats.SetBar(bar)
}

func main()  {
	fmt.Printf("./client -network=%s -codec=%s -compress=%s -h=%s -p=%d -total=%d -pipelining=%t -multiplexing=%t -batch=%t -batch_async=%t -noresponse=%t -clients=%d\n",network,codec,compress,host,port,total_calls,pipelining,multiplexing,batch,batch_async,noresponse,clients)
	var wrkClients []stats.Client
	opts:=rpc.DefaultOptions()
	opts.SetCompressType(compress)
	opts.SetBatch(batch)
	opts.SetBatchAsync(batch_async)
	opts.SetPipelining(pipelining)
	opts.SetMultiplexing(multiplexing)
	parallel:=1
	if clients>1{
		pool,err := rpc.DialsWithOptions(clients,network,addr,codec,opts)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		wrkClients=make([]stats.Client,len(pool.All()))
		for i:=0; i<len(pool.All());i++  {
			wrkClients[i]=&WrkClient{pool.All()[i]}
		}
		if batch{
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
		if batch{
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
		wrkClients[0]= &WrkClient{conn}
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
	A:= rand.Int31n(1000)
	B:= rand.Int31n(1000)
	req := &service.ArithRequest{A:A,B:B}
	if noresponse{
		err = c.Conn.CallNoResponse("Arith.Multiply", req) // 乘法运算
		if err==nil{
			return 0,0,true
		}
	}else {
		var res service.ArithResponse
		err = c.Conn.Call("Arith.Multiply", req, &res) // 乘法运算
		if res.Pro==A*B{
			return 0,0,true
		}else {
			fmt.Printf("err %d * %d = %d\n",A,B,res.Pro,)
		}
		if err != nil {
			fmt.Println("arith error: ", err)
		}
	}
	return 0,0,false
}