package main

import (
	"hslam.com/git/x/rpc/examples/example/pb/service"
	"hslam.com/git/x/rpc"
	_ "net/http/pprof"
	"math/rand"
	"net/http"
	"strconv"
	"runtime"
	"flag"
	"time"
	"log"
	"fmt"
	"sync"
)

var countchan chan int
var count int

var network string
var codec string
var compress string
var debug bool
var debug_port int
var host string
var port int
var log_once bool
var run_time_second int64
var addr string
var batch bool
var batch_async bool
var pipelining bool
var multiplexing bool
var noresponse bool
var clients int

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|fast|http|http2|quic|udp")
	flag.StringVar(&codec, "codec", "pb", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&compress, "compress", "no", "compress: -compress=no|flate|zlib|gzip")
	flag.BoolVar(&debug, "debug", true, "debug: -debug=false")
	flag.IntVar(&debug_port, "dp", 6061, "debug_port: -dp=6060")
	flag.StringVar(&host, "h", "localhost", "host: -h=localhost")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.BoolVar(&log_once, "log_once", false, "log_once: -log_once=false")
	flag.Int64Var(&run_time_second, "ts", 180, "run_time_second: -ts=60")
	flag.BoolVar(&batch, "batch", true, "batch: -batch=false")
	flag.BoolVar(&batch_async, "batch_async", true, "batch_async: -batch_async=false")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&multiplexing, "multiplexing", true, "pipelining: -pipelining=false")
	flag.BoolVar(&noresponse, "noresponse", false, "noresponse: -noresponse=false")
	flag.IntVar(&clients, "clients", 1, "num: -clients=1")
	log.SetFlags(0)
	flag.Parse()
	addr=host+":"+strconv.Itoa(port)
	countchan=make(chan int,500000)
}

func main()  {
	go func() {if debug{log.Println(http.ListenAndServe(":"+strconv.Itoa(debug_port), nil))}}()
	go func() {
		for c :=range countchan{
			count+=c
		}
	}()
	if clients>1{
		pool,err := rpc.Dials(clients,network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		pool.SetCompressType(compress)
		if batch {pool.EnableBatch()}
		if batch_async{pool.EnableBatchAsync()}
		if multiplexing{pool.EnableMultiplexing()}
		for i:=0;i<clients;i++{
			go run(pool.All()[i])
		}

	}else if clients==1 {
		conn, err:= rpc.Dial(network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		conn.SetCompressType(compress)
		if batch {conn.EnableBatch()}
		if multiplexing{conn.EnableMultiplexing()}
		go run(conn)
	}else {
		return
	}
	for i:=1;i<int(run_time_second);i++ {
		time.Sleep(time.Second)
		fmt.Println(i,count/i,count)
	}
	time.Sleep(time.Hour*3)
	fmt.Println(count/int(run_time_second),int(run_time_second*1000000000)/count)
}
func run(conn rpc.Client)  {
	var err error
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	if log_once{
		err = conn.Call("Arith.Multiply", req, &res) // 乘法运算
		if err != nil {
			fmt.Println("arith error: ", err)
		}
		fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)

		err = conn.Call("Arith.Divide", req, &res)
		if err != nil {
			fmt.Println("arith error: ", err)
		}
		fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
	}
	parallel:=1
	if batch{
		parallel=conn.GetMaxBatchRequest()
	}else if pipelining{
		parallel=conn.GetMaxRequests()
	}else if multiplexing{
		parallel=conn.GetMaxRequests()
	}
	if log_once{
		fmt.Println("parallel - ",parallel)
	}
	wg :=&sync.WaitGroup{}
	for i:=0;i<parallel;i++{
		go work(conn,countchan,wg)
	}
	wg.Wait()
	defer conn.Close()
}
func work(conn rpc.Client, countchan chan int,wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	start_time:=time.Now().UnixNano()
	var err error
	for{
		A:= rand.Int31n(1000)
		B:= rand.Int31n(1000)
		req := &service.ArithRequest{A:A,B:B}
		if noresponse{
			conn.CallNoResponse("Arith.Multiply", req)
			countchan<-1
		}else {
			var res service.ArithResponse
			err=conn.Call("Arith.Multiply", req, &res)
			if res.Pro==A*B{
				countchan<-1
			}else {
				fmt.Printf("err %d * %d = %d\n",A,B,res.Pro,)
			}
		}
		if err != nil {
			fmt.Println("arith error: ", err)
		}
		if time.Now().UnixNano()-start_time>run_time_second*1000000000{
			break
		}
	}
}