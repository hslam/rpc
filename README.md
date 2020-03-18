# rpc
A Golang implementation of RPC over TCP,QUIC,HTTP and IPC

## Feature
* **Network** tcp|quic|http|ipc
* **Codec** json/protobuf/xml/bytes/code
* **Compress** flate/zlib/gzip/no
* **Pipelining**
* **Multiplexing**
* **Batching**
* **Go/Call/CallNoRequest/CallNoResponse/OnlyCall**
* **Protocal** stream/message/frame
* **Pool/Transport**

Pipelining still requires the requests to be returned in the order requested, then there will be the head of line blocking problem.

Multiplexing allows the requests responses to be returned in an intermingled fashion so avoiding head of line blocking.

Batching is only useful when there are multiple goroutines calling it.

## Benchmark

### Linux Environment
* **CPU** 12 Cores 3.1 GHz
* **Memory** 24 GiB

#### Define

* **pkg**   package
* **trans** transport
* **c**     client
* **p**     pipelining
* **m**     multiplexing
* **b**     batching
* **n**     no response

#### Linux 12 CPU 3.1 GHz 24 GiB Requests per second
```
pkg     trans codec 1c    8c     1c_p   8c_p   1c_m   8c_m   1c_pb  8c_pb  1c_mb  8c_mb  1c_mbn
RPC     TCP   pb    29296 140652 279863 638999 259555 476963 429438 865441 424627 957064 1288990
RPC     HTTP  pb    29533 139781 283524 672418 264161 537103 428327 880574 434149 949597 1283761
RPC     QUIC  pb    10936 57506  198306 383062 181747 305278 177623 676530 174130 679271 1141024
NETRPC  TCP   gob   30931 144303 -      -      96493  391953 -      -      -      -      -
NETRPC  HTTP  gob   30847 144550 -      -      95691  383312 -      -      -      -      -
JSONRPC HTTP  json  26630 127039 -      -      93800  308515 -      -      -      -      -
RPCX    TCP   pb    28632 126713 -      -      92070  308490 -      -      -      -      -
GRPC    HTTP2 pb    15525 63556  -      -      116194 156065 -      -      -      -      -
```

#### Linux 12 CPU 3.1 GHz 24 GiB 99th percentile time (ms)
```
pkg     trans codec 1c    8c     1c_p   8c_p   1c_m   8c_m   1c_pb  8c_pb  1c_mb  8c_mb  1c_mbn
RPC     TCP   pb    0.04  0.10   0.57   2.47   0.74   3.69   2.36   30.45  2.21   20.56  6.61
RPC     HTTP  pb    0.04  0.10   0.54   2.42   0.66   3.34   2.35   32.05  2.17   24.93  6.29
RPC     QUIC  pb    0.12  0.58   0.81   3.60   1.00   4.34   0.64   2.51   0.63   2.21   0.94
NETRPC  TCP   gob   0.04  0.10   -      -      1.28   4.02   -      -      -      -      -
NETRPC  HTTP  gob   0.04  0.10   -      -      1.27   4.18   -      -      -      -      -
JSONRPC HTTP  json  0.05  0.13   -      -      1.27   5.12   -      -      -      -      -
RPCX    TCP   pb    0.05  0.21   -      -      1.52   5.86   -      -      -      -      -
GRPC    HTTP2 pb    0.10  0.43   -      -      1.24   8.86   -      -      -      -      -
```
./server -network=tcp -pipelining=false -multiplexing=true -batching=true

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batching=true -noresponse=false -clients=1
```
Summary:
	Clients:	1
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	2.36s
	Requests per second:	424627.55
	Fastest time for request:	0.31ms
	Average time per request:	1.20ms
	Slowest time for request:	3.71ms

Time:
	0.1%	time for request:	0.44ms
	1%	time for request:	0.56ms
	5%	time for request:	0.67ms
	10%	time for request:	0.76ms
	25%	time for request:	0.95ms
	50%	time for request:	1.17ms
	75%	time for request:	1.41ms
	90%	time for request:	1.67ms
	95%	time for request:	1.82ms
	99%	time for request:	2.21ms
	99.9%	time for request:	2.77ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batching=true -noresponse=false -clients=8
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	1.04s
	Requests per second:	957064.19
	Fastest time for request:	0.07ms
	Average time per request:	4.22ms
	Slowest time for request:	64.40ms

Time:
	0.1%	time for request:	0.21ms
	1%	time for request:	0.35ms
	5%	time for request:	0.56ms
	10%	time for request:	0.73ms
	25%	time for request:	1.35ms
	50%	time for request:	3.04ms
	75%	time for request:	5.70ms
	90%	time for request:	8.60ms
	95%	time for request:	12.41ms
	99%	time for request:	20.56ms
	99.9%	time for request:	31.83ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

## Get started

### Install
```
go get github.com/hslam/rpc
```
### Import
```
import "github.com/hslam/rpc"
```

## [Example](https://github.com/hslam/rpc/tree/masterll/examples "examples")
### arith.proto
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
### arith.pb.go
**[protobuf](https://github.com/protocolbuffers/protobuf "protobuf")**
```
protoc ./arith.proto --go_out=./
```

**[gogoproto](https://github.com/gogo/protobuf "gogoproto")**
```
protoc ./arith.proto --gofast_out=./
```
### arith.go
```
package service
type Arith struct {}
func (this *Arith) Multiply(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}
```

### server.go
```go
package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("tcp",":8080")//tcp|quic|http
}
```

### client.go
```go
package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
	"fmt"
	"log"
)
func main()  {
	conn, err:= rpc.Dial("tcp","127.0.0.1:8080","pb")//tcp|quic|http
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	defer conn.Close()
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	err = conn.Call("Arith.Multiply", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
```
### transport.go
```go
package main
import (
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc"
	"fmt"
	"log"
	"time"
)
func main()  {
	MaxConnsPerHost:=2
	MaxIdleConnsPerHost:=0
	transport:=rpc.NewTransport(MaxConnsPerHost,MaxIdleConnsPerHost,"tcp","pb",rpc.DefaultOptions())//tcp|quic|http
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	var err error
	err = transport.Call("127.0.0.1:8081","Arith.Multiply", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
}
```
### Output
```
9 * 2 = 18
```
### License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)


### Authors
rpc was written by Meng Huang.