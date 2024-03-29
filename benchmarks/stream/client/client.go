package main

import (
	"flag"
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/benchmarks/stream/service"
	"github.com/hslam/stats"
	"log"
	"math/rand"
)

var network string
var addr string
var codec string
var clients int
var total int
var parallel int
var bar bool

func init() {
	flag.StringVar(&network, "network", "tcp", "-network=tcp")
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.StringVar(&codec, "codec", "pb", "-codec=code")
	flag.IntVar(&total, "total", 100000, "-total=100000")
	flag.IntVar(&parallel, "parallel", 1, "-parallel=1")
	flag.IntVar(&clients, "clients", 1, "-clients=1")
	flag.BoolVar(&bar, "bar", true, "-bar=true")
	flag.Parse()
	stats.SetBar(bar)
	fmt.Printf("./client -network=%s -addr=%s -codec=%s -total=%d -parallel=%d -clients=%d\n", network, addr, codec, total, parallel, clients)
}

func main() {
	if clients < 1 || parallel < 1 || total < 1 {
		return
	}
	var wrkClients []stats.Client
	for i := 0; i < clients; i++ {
		conn, err := rpc.Dial(network, addr, codec)
		if err != nil {
			log.Fatalln("dailing error: ", err)
		}
		conn.SetNoCopy(true)
		for j := 0; j < parallel; j++ {
			stream, err := conn.NewStream("Arith.StreamMultiply")
			if err != nil {
				panic(err)
			}
			wrkClients = append(wrkClients, &WrkClient{Stream: stream, buf: make([]byte, 8)})
		}
	}
	parallel = 1
	stats.StartPrint(parallel, total, wrkClients)
}

type WrkClient struct {
	rpc.Stream
	buf []byte
}

func (c *WrkClient) Call() (int64, int64, bool) {
	A := rand.Int31n(100)
	B := rand.Int31n(100)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	if err := c.Stream.WriteMessage(req); err != nil {
		fmt.Println(err)
		return 0, 0, false
	}
	if err := c.Stream.ReadMessage(c.buf, &res); err != nil {
		fmt.Println(err)
		return 0, 0, false
	}
	if res.Pro == A*B {
		return 0, 0, true
	}
	fmt.Printf("err %d * %d = %d\n", A, B, res.Pro)
	return 0, 0, false
}
