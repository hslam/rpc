# rpc
A Golang implementation of RPC over TCP,WS,QUIC and HTTP

## Server Feature

* **Network** tcp|quic|http|ws
* **Pipelining** async/sync
* **Multiplexing**

## Client Feature
* **Network** tcp|quic|http|ws
* **Codec** json/protobuf/xml/bytes/gencode
* **Compress** flate/zlib/gzip/no
* **Pipelining**
* **Multiplexing**
* **Batch** async/sync
* **Go/Call/CallNoRequest/CallNoResponse/OnlyCall**
* **Protocal** stream/message/frame
* **Pool/Transport**

Batch is only useful when there are multiple goroutines calling it.

## Benchmark

### Linux Environment
* **CPU** 12 Cores 3.1 GHz
* **Memory** 24 GiB

#### Define

* **pkg**   package
* **trans** transport
* **c**     client
* **p**     pipeline
* **m**     multiplex
* **b**     batch
* **n**     no response

#### Linux 12 CPU 3.1 GHz 24 GiB Requests per second
```
pkg     trans codec 1c    8c     1c_p   8c_p   1c_m   8c_m   1c_pb  8c_pb  1c_mb  8c_mb  1c_mbn
RPC     TCP   pb    27807 140157 189640 501542 175541 406636 383769 807446 381797 910869 1271944
RPC     HTTP  pb    27943 139138 189945 522737 178484 442490 385941 831232 389369 933167 1269670
RPC     QUIC  pb    10955 57936  128896 323773 115542 263563 127853 530690 127902 530925 1006680
NETRPC  TCP   gob   30340 141675 -      -      95282  364254 -      -      -      -      -
NETRPC  HTTP  gob   30325 141484 -      -      96946  398242 -      -      -      -      -
JSONRPC HTTP  json  26872 124867 -      -      92910  299619 -      -      -      -      -
GRPC    HTTP2 pb    15649 62359  -      -      90411  146727 -      -      -      -      -
```

#### Linux 12 CPU 3.1 GHz 24 GiB 99th percentile time (ms)
```
pkg     trans codec 1c    8c     1c_p   8c_p   1c_m   8c_m   1c_pb  8c_pb  1c_mb  8c_mb  1c_mbn
RPC     TCP   pb    0.05  0.10   0.36   1.40   0.45   1.80   2.60   26.45  2.44   14.63  7.04
RPC     HTTP  pb    0.05  0.10   0.33   1.38   0.36   1.58   2.56   31.02  2.36   22.15  6.77
RPC     QUIC  pb    0.12  0.57   0.64   2.30   0.79   2.50   0.43   1.40   0.44   1.45   0.48
NETRPC  TCP   gob   0.04  0.10   -      -      0.61   2.15   -      -      -      -      -
NETRPC  HTTP  gob   0.04  0.10   -      -      2.47   7.32   -      -      -      -      -
JSONRPC HTTP  json  0.04  0.12   -      -      0.62   2.78   -      -      -      -      -
GRPC    HTTP2 pb    0.11  0.43   -      -      0.80   4.81   -      -      -      -      -
```
./server -network=tcp -async=false -pipelining=false -multiplexing=true -batch=true

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
```
Summary:
	Clients:	1
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	2.62s
	Requests per second:	381797.86
	Fastest time for request:	0.33ms
	Average time per request:	1.34ms
	Slowest time for request:	3.81ms

Time:
	0.1%	time for request:	0.52ms
	1%	time for request:	0.70ms
	5%	time for request:	0.88ms
	10%	time for request:	0.98ms
	25%	time for request:	1.13ms
	50%	time for request:	1.26ms
	75%	time for request:	1.47ms
	90%	time for request:	1.86ms
	95%	time for request:	2.07ms
	99%	time for request:	2.44ms
	99.9%	time for request:	3.07ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=8
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	1.10s
	Requests per second:	910869.59
	Fastest time for request:	0.07ms
	Average time per request:	4.45ms
	Slowest time for request:	55.33ms

Time:
	0.1%	time for request:	0.37ms
	1%	time for request:	0.69ms
	5%	time for request:	1.21ms
	10%	time for request:	1.68ms
	25%	time for request:	2.61ms
	50%	time for request:	3.93ms
	75%	time for request:	5.57ms
	90%	time for request:	7.64ms
	95%	time for request:	9.53ms
	99%	time for request:	14.63ms
	99.9%	time for request:	23.34ms

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

## [Example](https://github.com/hslam/rpc/tree/master/examples "examples")
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
	"github.com/hslam/rpc/examples/helloworld/pb/service"
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
	"github.com/hslam/rpc/examples/helloworld/pb/service"
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
	"github.com/hslam/rpc/examples/helloworld/pb/service"
	"github.com/hslam/rpc"
	"fmt"
	"log"
)
func main()  {
	MaxConnsPerHost:=2
	MaxIdleConnsPerHost:=0
	transport:=rpc.NewTransport(MaxConnsPerHost,MaxIdleConnsPerHost,"tcp","pb",rpc.DefaultOptions())//tcp|quic|http
	req := &service.ArithRequest{A:9,B:2}
	var res service.ArithResponse
	var err error
	err = transport.Call("Arith.Multiply", req, &res,"127.0.0.1:8080")
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
### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Meng Huang)


### Authors
rpc was written by Meng Huang.