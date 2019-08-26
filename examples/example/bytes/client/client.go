package main

import (
	"hslam.com/mgit/Mort/rpc"
	_ "net/http/pprof"
	"math/rand"
	"net/http"
	"strconv"
	"runtime"
	"strings"
	"flag"
	"time"
	"log"
	"fmt"
	"bytes"
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
var pipelining bool
var noresponse bool
var norequest bool
var onlycall bool
var clients int

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|fast|http|http2|quic|udp")
	flag.StringVar(&codec, "codec", "bytes", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&compress, "compress", "no", "compress: -compress=no|flate|zlib|gzip")
	flag.BoolVar(&debug, "debug", false, "debug: -debug=false")
	flag.IntVar(&debug_port, "dp", 6060, "debug_port: -dp=6060")
	flag.StringVar(&host, "h", "localhost", "host: -h=localhost")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.BoolVar(&log_once, "log_once", true, "log_once: -log_once=false")
	flag.Int64Var(&run_time_second, "ts", 180, "run_time_second: -ts=60")
	flag.BoolVar(&batch, "batch", false, "batch: -batch=false")
	flag.BoolVar(&pipelining, "pipelining", false, "pipelining: -pipelining=false")
	flag.BoolVar(&noresponse, "noresponse", false, "noresponse: -noresponse=false")
	flag.BoolVar(&norequest, "norequest", false, "norequest: -norequest=false")
	flag.BoolVar(&onlycall, "onlycall", false, "onlycall: -onlycall=false")
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
		if batch {pool.EnabledBatch()}
		for i:=0;i<clients;i++{
			go run(pool.Get())
		}

	}else if clients==1 {
		conn, err:= rpc.Dial(network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		conn.SetCompressType(compress)
		if batch {conn.EnabledBatch()}
		go run(conn)
	}else {
		return
	}
	for i:=1;i<int(run_time_second);i++ {
		time.Sleep(time.Second)
		fmt.Println(i,count/i)
	}
	time.Sleep(time.Second*3)
	fmt.Println(count/int(run_time_second),int(run_time_second*1000000000)/count)
}
func run(conn *rpc.Client)  {
	var err error
	var req =[]byte("Hello World")
	var res []byte
	if log_once{
		err = conn.CallNoResponse("Echo.Set", &req)
		fmt.Printf("Echo.Set : %s\n", string(req))
		err = conn.CallNoRequest("Echo.Get", &res)
		fmt.Printf("Echo.Get : %s\n", string(res))
		err = conn.OnlyCall("Echo.Clear")
		fmt.Printf("Echo.Clear\n")
		err = conn.CallNoRequest("Echo.Get", &res)
		fmt.Printf("Echo.Get : %s\n", string(res))
		err = conn.CallNoResponse("Echo.Set", &req)
		err = conn.Call("Echo.ToLower", &req, &res)
		if err != nil {
			log.Fatalln("Echo error: ", err)
		}
		fmt.Printf("Echo.ToLower : %s\n", string(res))
	}
	parallel:=1
	if batch{
		parallel=conn.GetMaxBatchRequest()
	}else if pipelining{
		parallel=conn.GetMaxPipelineRequest()
	}
	if log_once{
		fmt.Println("parallel - ",parallel)
	}
	for i:=0;i<parallel;i++{
		go work(conn,countchan)
	}
	defer conn.Close()
	select {}
}
func work(conn *rpc.Client, countchan chan int) {
	start_time:=time.Now().UnixNano()
	var err error
	len:=10
	for{
		req := make([]byte, len)
		for i := 0; i < len; i++ {
			b := rand.Intn(26) + 65
			req[i] = byte(b)
		}
		if noresponse{
			err = conn.CallNoResponse("Echo.Set", &req)
			countchan<-1
		}else if onlycall{
			err = conn.OnlyCall("Echo.Clear")
			countchan<-1
		}else if norequest{
			var req =[]byte("Hello World")
			var res []byte
			err = conn.CallNoRequest("Echo.Get", &res)
			if bytes.Equal(res,[]byte("Hello World")){
				countchan<-1
			}else {
				fmt.Printf("Echo.Get is not equal: req-%s res-%s\n",string(req),string(res))
			}
		}else {
			var res []byte
			err = conn.Call("Echo.ToLower", &req, &res)
			if err != nil {
				log.Fatalln("Echo error: ", err)
			}
			if bytes.Equal(res,[]byte(strings.ToLower(string(req)))){
				countchan<-1
			}else {
				fmt.Printf("Echo.ToLower is not equal: req-%s res-%s\n",string(req),string(res))
			}
		}
		if err != nil {
			log.Fatalln("Echo error: ", err)
		}
		if time.Now().UnixNano()-start_time>run_time_second*1000000000{
			break
		}
	}
}
