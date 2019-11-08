# rpc
A RPC implementation written in Golang over TCP UDP QUIC WS HTTP HTTP2

## Server Feature
* Pipelining async/sync
* Multiplexing
## Client Feature
* Network tcp/ws/http/http2/quic/udp
* Codec json/protobuf/xml
* Compress flate/zlib/gzip/no
* Pipelining
* Multiplexing
* Batch async/sync
* Call/CallNoRequest/CallNoResponse/OnlyCall
* Protocal stream/message/frame/conn
* Pool

## Benchmark
Batch is only useful when there are multiple goroutines calling it.

### Env

```
Mac 4 CPU 8 GiB
```

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
```
Summary:
	Clients:	1
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	4.75s
	Requests per second:	210329.36
	Fastest time for request:	0.70ms
	Average time per request:	2.43ms
	Slowest time for request:	14.70ms

Time:
	0.1%	time for request:	0.94ms
	1%	time for request:	1.26ms
	5%	time for request:	1.41ms
	10%	time for request:	1.53ms
	25%	time for request:	1.92ms
	50%	time for request:	2.29ms
	75%	time for request:	2.86ms
	90%	time for request:	3.38ms
	95%	time for request:	3.73ms
	99%	time for request:	5.11ms
	99.9%	time for request:	6.20ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=8
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	3.30s
	Requests per second:	303380.20
	Fastest time for request:	0.70ms
	Average time per request:	13.39ms
	Slowest time for request:	72.79ms

Time:
	0.1%	time for request:	1.79ms
	1%	time for request:	3.80ms
	5%	time for request:	5.98ms
	10%	time for request:	7.03ms
	25%	time for request:	8.92ms
	50%	time for request:	11.74ms
	75%	time for request:	16.04ms
	90%	time for request:	21.95ms
	95%	time for request:	26.25ms
	99%	time for request:	37.58ms
	99.9%	time for request:	60.24ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=8
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	1.75s
	Requests per second:	570064.98
	Fastest time for request:	0.00ms
	Average time per request:	6.77ms
	Slowest time for request:	290.49ms

Time:
	0.1%	time for request:	0.00ms
	1%	time for request:	0.00ms
	5%	time for request:	0.00ms
	10%	time for request:	0.00ms
	25%	time for request:	0.00ms
	50%	time for request:	0.00ms
	75%	time for request:	2.37ms
	90%	time for request:	17.91ms
	95%	time for request:	45.54ms
	99%	time for request:	105.99ms
	99.9%	time for request:	206.51ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

### Define
* c Client
* p Pipeline
* m Multiplex
* b Batch
* n No Response


### Mac 4 CPU 8 GiB Localhost  qps
```
package     transport   codec       1c      2c      1c_p    2c_p    1c_n    2c_n    1c_pb   2c_pb   1c_m    2c_m    1c_mb   2c_mb   1c_mbn  2c_mbn
RPC         TCP         protobuf    7640    14899   38460   40417   62296   68947   142992  165612  35846   37042   201432  302815  517066  534081
RPC         WS          protobuf    7084    13762   36852   38483   58523   68162   147906  189202  32265   35351   201778  295464  541553  538206
RPC         UDP         protobuf    5136    9739    19807   17739   23525   18901   21351   45916   19731   18831   28739   51115   19513   35540
RPC         QUIC        protobuf    3959    7879    13673   12980   27098   26474   22274   46708   13333   12960   23503   51506   152819  151515
RPC         HTTP        protobuf    3694    7530    4953    8470    5348    8908    130125  163166  4962    8107    149764  228230  301091  690483
RPC         HTTP2       protobuf    2536    4304    5933    6575    7213    6215    109387  146322  6707    6695    154804  189480  380008  351820
```

### Mac 4 CPU 8 GiB Localhost 99th percentile time (ms)
```
package     transport   codec       1c      2c      1c_p    2c_p    1c_n    2c_n    1c_pb   2c_pb   1c_m    2c_m    1c_mb   2c_mb   1c_mbn  2c_mbn
RPC         TCP         protobuf    0.22    0.22    1.91    4.22    0.09    0.21    7.00    13.53   1.66    4.52    5.35    7.92    9.23    22.31
RPC         WS          protobuf    0.25    0.31    1.71    3.95    0.09    0.18    6.59    10.98   1.96    4.35    5.40    8.04    7.62    29.04
RPC         UDP         protobuf    0.31    0.42    4.70    17.67   0.42    1.17    1.35    0.93    5.56    14.92   0.59    0.75    0.77    3.48
RPC         QUIC        protobuf    0.38    0.58    3.54    8.92    0.19    0.50    0.65    0.65    3.73    8.28    0.64    0.60    0.36    1.04
RPC         HTTP        protobuf    0.47    0.53    8.73    10.30   0.38    0.56    7.73    13.27   9.30    11.75   7.48    13.02   7.91    31.05
RPC         HTTP2       protobuf    0.68    1.07    13.63   19.82   3.07    8.64    9.81    14.35   9.65    20.36   7.70    11.72   22.76   87.47
```

## Example
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
    int32 quo = 2;
    int32 rem = 3;
}
```

### arith.pb.go
```
protoc ./arith.proto --go_out=./

```
### arith.go
```
package service
import (
	"errors"
)
type Arith struct {}
func (this *Arith) Multiply(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}
func (this *Arith) Divide(req *ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}
```

### server.go
```go
package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/pb/service"
	"hslam.com/git/x/rpc"
	"strconv"
	"flag"
)
var network string
var port int
var saddr string
func init()  {
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|http|http2|quic|udp")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	rpc.Register(new(service.Arith))
	rpc. (network,saddr)
}
```

### client.go
```go
package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/pb/service"
	"hslam.com/git/x/rpc"
	"strconv"
	"flag"
	"log"
	"fmt"
)
var network string
var codec string
var host string
var port int
var addr string
func init()  {
	flag.StringVar(&network, "network", "tcp", "network: -network=tcp|ws|http|http2|quic|udp")
	flag.StringVar(&codec, "codec", "pb", "codec: -codec=pb|json|xml|bytes")
	flag.StringVar(&host, "h", "localhost", "host: -h=localhost")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Parse()
	addr=host+":"+strconv.Itoa(port)
}
func main()  {
	conn, err:= rpc.Dial(network,addr,codec)
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

	err = conn.Call("Arith.Divide", req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
}
```

## Example http json javascript

### http-json-server.go
```go
package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/json/service"
	"hslam.com/git/x/rpc"
	"strconv"
	"flag"
)
var network string
var port int
var saddr string
func init()  {
	flag.StringVar(&network, "network", "http", "network: -network=tcp|ws|http|http2|quic|udp")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe(network,saddr)
}
```
### http-json-client-javascript
```
<script type="text/javascript" src="./js/lib/rpc/rpc.min.js"></script>
<script type="text/javascript">
    ArithRequest = function(A,B) {
        this.a=A;
        this.b=B;
    }
    var client = new rpc.Dial("127.0.0.1:9999");
    var req = new ArithRequest(9,2)
    var reply=client.Call("Arith.Multiply",req)
    console.log(req.a.toString()+" * "+req.b.toString()+" = "+reply.pro.toString());
    var reply=client.Call("Arith.Divide",req)
    console.log(req.a.toString()+" / "+req.b.toString()+", quo is "+reply.quo.toString()+", rem is "+reply.rem.toString());
</script>
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)


### Authors
workerpool was written by Mort Huang.