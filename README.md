# rpc
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/rpc)](https://pkg.go.dev/github.com/hslam/rpc)
[![Build Status](https://github.com/hslam/rpc/workflows/build/badge.svg)](https://github.com/hslam/rpc/actions)
[![codecov](https://codecov.io/gh/hslam/rpc/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/rpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/rpc?v=7e100)](https://goreportcard.com/report/github.com/hslam/rpc)
[![LICENSE](https://img.shields.io/github/license/hslam/rpc.svg?style=flat-square)](https://github.com/hslam/rpc/blob/master/LICENSE)

Package rpc implements a remote procedure call over TCP, UNIX, HTTP and WS. The rpc improves throughput and reduces latency. Up to 4 times faster than net/rpc.

## Feature
* More throughput and less latency.
* **[Netpoll](https://github.com/hslam/netpoll "netpoll")** epoll/kqueue/net
* **[Network](https://github.com/hslam/socket "socket")** tcp/unix/http/[ws](https://github.com/hslam/websocket "websocket")
* **[Codec](https://github.com/hslam/codec "codec")** json/[code](https://github.com/hslam/code "code")/pb
* Multiplexing/Pipelining
* [Auto batching](https://github.com/hslam/writer "writer")
* Call/Go/RoundTrip/Ping/Watch/CallWithContext
* Server push
* Conn/Transport/Client
* TLS

**Comparison to other packages**

|Package| [netrpc](https://github.com/golang/go/tree/master/src/net/rpc "netrpc")| [jsonrpc](https://github.com/golang/go/tree/master/src/net/rpc/jsonrpc "jsonrpc")|[rpc](https://github.com/hslam/rpc "rpc")|[grpc](https://google.golang.org/grpc "grpc")|[rpcx](https://github.com/smallnest/rpcx "rpcx")|
|:--:|:--|:--|:--|:--|:--|
|Epoll/Kqueue|No|No|Yes|No|No|
|Multiplexing|Yes|Yes|Yes|Yes|Yes|
|Pipelining|No|No|Yes|No|No|
|Auto Batching|No|No|Yes|No|No|
|Transport|No|No|Yes|No|No|
|Server Push|No|No|Yes|Yes|Yes|

## [Benchmark](http://github.com/hslam/rpc-benchmark "rpc-benchmark")

##### Low Concurrency

<img src="https://raw.githubusercontent.com/hslam/rpc-benchmark/master/rpc-bar-qps.png" width = "400" height = "300" alt="rpc" align=center><img src="https://raw.githubusercontent.com/hslam/rpc-benchmark/master/rpc-bar-p99.png" width = "400" height = "300" alt="rpc" align=center>

##### High Concurrency

<img src="https://raw.githubusercontent.com/hslam/rpc-benchmark/master/rpc-curve-qps.png" width = "400" height = "300" alt="rpc" align=center><img src="https://raw.githubusercontent.com/hslam/rpc-benchmark/master/rpc-curve-p99.png" width = "400" height = "300" alt="rpc" align=center>

## Get started

### Install
```
go get github.com/hslam/rpc
```
### Import
```
import "github.com/hslam/rpc"
```

### Usage
#### [Examples](https://github.com/hslam/rpc/tree/master/examples "examples")
arith.proto
```
syntax = "proto3";
package service;

message ArithRequest {
    int32 a = 1;
    int32 b = 2;
}

message ArithResponse {
    int32 pro = 1;
}
```

**[GoGo Protobuf ](https://github.com/gogo/protobuf "gogoprotobuf")**
```
protoc ./arith.proto --gogofaster_out=./
```

arith.go
```go
package service

type Arith struct{}

func (a *Arith) Multiply(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}
```

server.go
```go
package main

import (
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/codec/pb/service"
)

func main() {
	rpc.Register(new(service.Arith))
	rpc.Listen("tcp", ":9999", "pb")
}
```

conn.go
```go
package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/codec/pb/service"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "pb")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	if err = conn.Call("Arith.Multiply", req, &res); err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
```

transport.go
```go
package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/codec/pb/service"
)

func main() {
	trans := &rpc.Transport{
		MaxConnsPerHost:     1,
		MaxIdleConnsPerHost: 1,
		Options:             &rpc.Options{Network: "tcp", Codec: "pb"},
	}
	defer trans.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	if err := trans.Call(":9999", "Arith.Multiply", req, &res); err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
```

client.go
```go
package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/codec/pb/service"
)

func main() {
	opts := &rpc.Options{Network: "tcp", Codec: "pb"}
	client := rpc.NewClient(opts, ":9997", ":9998", ":9999")
	client.Scheduling = rpc.LeastTimeScheduling
	defer client.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	if err := client.Call("Arith.Multiply", req, &res); err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
```

context.go
```go
package main

import (
	"context"
	"fmt"
	"github.com/hslam/rpc"
	"github.com/hslam/rpc/examples/codec/pb/service"
	"time"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "pb")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	req := &service.ArithRequest{A: 9, B: 2}
	var res service.ArithResponse
	valueCtx := context.WithValue(context.Background(), rpc.ContextKeyBuffer, make([]byte, 64))
	ctx, cancel := context.WithTimeout(valueCtx, time.Minute)
	defer cancel()
	err = conn.CallWithContext(ctx, "Arith.Multiply", req, &res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
```

#### Output
```
9 * 2 = 18
```


### License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)


### Author
rpc was written by Meng Huang.

