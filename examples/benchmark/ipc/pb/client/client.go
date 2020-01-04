package main

import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
	"github.com/hslam/stats"
	"math/rand"
	"runtime"
	"flag"
	"log"
	"fmt"
)
var network string
var codec string
var compress string
var addr string
var batching bool
var pipelining bool
var multiplexing bool
var noresponse bool
var mock string
var clients int
var total_calls int
var bar bool
func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "ipc", "network: -network=ipc")
	flag.StringVar(&codec, "codec", "pb", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&compress, "compress", "no", "compress: -compress=no|flate|zlib|gzip")
	flag.StringVar(&addr, "address", "/tmp/ipc", "address: -address=/tmp/ipc")
	flag.IntVar(&total_calls, "total", 1000000, "total_calls: -total=10000")
	flag.BoolVar(&batching, "batching", false, "batching: -batching=false")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "pipelining: -pipelining=false")
	flag.BoolVar(&noresponse, "noresponse", false, "noresponse: -noresponse=false")
	flag.StringVar(&mock, "mock", "", "mock: -mock=|0ms|1ms|10ms|50ms")
	flag.IntVar(&clients, "clients", 1, "clients: -clients=1")
	flag.BoolVar(&bar, "bar", false, "bar: -bar=true")
	log.SetFlags(0)
	flag.Parse()
	stats.SetBar(bar)
}

func main()  {
	fmt.Printf("./client -network=%s -codec=%s -compress=%s -address=%s -total=%d -pipelining=%t -multiplexing=%t -batching=%t -noresponse=%t -clients=%d\n",network,codec,compress,addr,total_calls,pipelining,multiplexing,batching,noresponse,clients)
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
		wrkClients=make([]stats.Client,len(pool.All()))
		for i:=0; i<len(pool.All());i++  {
			wrkClients[i]=&WrkClient{pool.All()[i]}
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
		err = c.Conn.CallNoResponse("Arith.Multiply"+mock, req) // 乘法运算
		if err==nil{
			return 0,0,true
		}
	}else {
		var res service.ArithResponse
		err = c.Conn.Call("Arith.Multiply"+mock, req, &res) // 乘法运算
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