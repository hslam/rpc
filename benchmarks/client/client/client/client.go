package main

import (
	"flag"
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/benchmarks/client/client/service"
	"github.com/hslam/socket"
	"github.com/hslam/stats"
	"math/rand"
)

var addr string
var clients int
var total int
var parallel int
var bar bool

func init() {
	flag.StringVar(&addr, "addr", ":9999", "-addr=:9999")
	flag.IntVar(&total, "total", 100000, "-total=100000")
	flag.IntVar(&parallel, "parallel", 1, "-parallel=1")
	flag.IntVar(&clients, "clients", 1, "-clients=1")
	flag.BoolVar(&bar, "bar", true, "-bar=true")
	flag.Parse()
	stats.SetBar(bar)
	fmt.Printf("./client -addr=%s -total=%d -parallel=%d -clients=%d\n", addr, total, parallel, clients)
}

func main() {
	if clients < 1 || parallel < 1 || total < 1 {
		return
	}
	var wrkClients []stats.Client
	opts := &rpc.Options{NewSocket: socket.NewTCPSocket, NewCodec: rpc.NewPBCodec, NewHeaderEncoder: rpc.NewPBEncoder}
	for i := 0; i < clients; i++ {
		client := rpc.NewClient(opts, ":9999")
		wrkClients = append(wrkClients, &WrkClient{client})
	}
	stats.StartPrint(parallel, total, wrkClients)
}

type WrkClient struct {
	*rpc.Client
}

func (c *WrkClient) Call() (int64, int64, bool) {
	A := rand.Int31n(100)
	B := rand.Int31n(100)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	if err := c.Transport.Call(addr, "Arith.Multiply", req, &res); err != nil {
		fmt.Println(err)
		return 0, 0, false
	}
	if res.Pro == A*B {
		return 0, 0, true
	} else {
		fmt.Printf("err %d * %d = %d != %d\n", A, B, A*B, res.Pro)
	}
	return 0, 0, false
}
