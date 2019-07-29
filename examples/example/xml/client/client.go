package main

import (
	"hslam.com/mgit/Mort/rpc/examples/example/json/service"
	"hslam.com/mgit/Mort/rpc"
	_ "net/http/pprof"
	"math/rand"
	"net/http"
	"strconv"
	"runtime"
	"flag"
	"time"
	"log"
	"fmt"
)

var countchan chan int
var count int

var network string
var codec string
var debug bool
var debug_port int
var host string
var port int
var run_time_second int64
var addr string
var batch bool
var concurrent bool
var noresponse bool
var clients int

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&network, "network", "tcp", "network: -network=fast|ws|tcp|quic|udp")
	flag.StringVar(&codec, "codec", "xml", "codec: -codec=pb|json|xml")
	flag.BoolVar(&debug, "debug", false, "debug: -debug=false")
	flag.IntVar(&debug_port, "dp", 6060, "debug_port: -dp=6060")
	flag.StringVar(&host, "h", "localhost", "host: -h=localhost")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Int64Var(&run_time_second, "ts", 60, "run_time_second: -ts=60")
	flag.BoolVar(&batch, "batch", false, "batch: -batch=false")
	flag.BoolVar(&concurrent, "concurrent", false, "concurrent: -concurrent=false")
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
		if batch {pool.EnabledBatch()}
		for i:=0;i<clients;i++{
			go run(pool.Get())
		}

	}else if clients==1 {
		conn, err:= rpc.Dial(network,addr,codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
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
func run(conn rpc.Client)  {
	var err error
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	err = conn.Call("Arith.Multiply", req, &res) // 乘法运算
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)

	err = conn.Call("Arith.Divide", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
	if batch{
		for i:=0;i<conn.GetMaxBatchRequest();i++ {
			go work(conn,countchan)
		}
	}else if concurrent{
		for i:=0;i<conn.GetMaxConcurrentRequest();i++ {
			go work(conn,countchan)
		}
	}else {
		go work(conn,countchan)
	}
	defer conn.Close()
	select {}
}
func work(conn rpc.Client, countchan chan int) {
	start_time:=time.Now().UnixNano()
	var err error
	for{
		A:= rand.Int31n(100)
		B:= rand.Int31n(100)
		req := &service.ArithRequest{A:A,B:B}
		if noresponse{
			err = conn.CallNoResponse("Arith.Multiply", req) // 乘法运算
			countchan<-1
		}else {
			var res service.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res) // 乘法运算
			if res.Pro==A*B{
				countchan<-1
			}else {
				fmt.Printf("err %d * %d = %d\n",A,B,res.Pro,)
			}
		}
		if err != nil {
			log.Fatalln("arith error: ", err)
		}
		if time.Now().UnixNano()-start_time>run_time_second*1000000000{
			break
		}
	}
}