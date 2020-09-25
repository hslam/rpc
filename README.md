# rpc
Package rpc provides access to the exported methods of an object across a network or other I/O connection.

## Feature
* **[Netpoll](https://github.com/hslam/netpoll "netpoll")** epoll/kqueue/net
* **[Network](https://github.com/hslam/socket "socket")** tcp/unix/http/[ws](https://github.com/hslam/websocket "websocket")
* **[Codec](https://github.com/hslam/codec "codec")** json/[code](https://github.com/hslam/code "code")/pb
* Multiplexing/Pipelining
* [Auto Batching](https://github.com/hslam/writer "writer")
* Call/Go/RoundTrip/Ping
* Client/Transport
* TLS

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
```
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

