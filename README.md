# rpc
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/rpc)](https://pkg.go.dev/github.com/hslam/rpc)
[![Build Status](https://travis-ci.org/hslam/rpc.svg?branch=master)](https://travis-ci.org/hslam/rpc)
[![codecov](https://codecov.io/gh/hslam/rpc/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/rpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/rpc?v=7e100)](https://goreportcard.com/report/github.com/hslam/rpc)
[![LICENSE](https://img.shields.io/github/license/hslam/rpc.svg?style=flat-square)](https://github.com/hslam/rpc/blob/master/LICENSE)

Package rpc provides access to the exported methods of an object across a network or other I/O connection.

## Feature
* More throughput and less latency.
* **[Netpoll](https://github.com/hslam/netpoll "netpoll")** epoll/kqueue/net
* **[Network](https://github.com/hslam/socket "socket")** tcp/unix/http/[ws](https://github.com/hslam/websocket "websocket")
* **[Codec](https://github.com/hslam/codec "codec")** json/[code](https://github.com/hslam/code "code")/pb
* Multiplexing/Pipelining
* [Auto batching](https://github.com/hslam/writer "writer")
* Call/Go/RoundTrip/Ping/Watch
* Server push
* Client/Transport
* TLS

**Comparison to other packages**

|Package| [netrpc](https://github.com/golang/go/tree/master/src/net/rpc "netrpc")| [jsonrpc](https://github.com/golang/go/tree/master/src/net/rpc/jsonrpc "jsonrpc")|[rpc](https://github.com/hslam/rpc "rpc")|[grpc](https://google.golang.org/grpc "grpc")|[rpcx](https://github.com/smallnest/rpcx "rpcx")|
|:--:|:--|:--|:--|:--|:--|
|Epoll/Kqueue|No|No|Yes|No|No|
|Multiplexing|Yes|Yes|Yes|Yes|Yes|
|Pipelining|No|No|Yes|No|No|
|Auto Batching|No|No|Yes|No|No|
|Client|Yes|Yes|Yes|Yes|Yes|
|Transport|No|No|Yes|No|No|
|Websocket|No|No|Yes|No|No|
|Server Push|No|No|Yes|Yes|Yes|

## [Benchmark](http://github.com/hslam/rpc-benchmark "rpc-benchmark")

##### Low Concurrency

<img src="https://raw.githubusercontent.com/hslam/rpc/master/rpc-bar-qps.png" width = "400" height = "300" alt="rpc" align=center><img src="https://raw.githubusercontent.com/hslam/rpc/master/rpc-bar-p99.png" width = "400" height = "300" alt="rpc" align=center>

##### High Concurrency

<img src="https://raw.githubusercontent.com/hslam/rpc/master/rpc-curve-qps.png" width = "400" height = "300" alt="rpc" align=center><img src="https://raw.githubusercontent.com/hslam/rpc/master/rpc-curve-p99.png" width = "400" height = "300" alt="rpc" align=center>

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


client.go
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
	err = conn.Call("Arith.Multiply", req, &res)
	if err != nil {
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
	var err error
	err = trans.Call(":9999", "Arith.Multiply", req, &res)
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

