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

## example benchmark
### ENV

```
Mac 4 CPU 8 GiB
```

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
```
Summary:
	Clients:	1
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	5.13s
	Requests per second:	194881.33
	Fastest time for request:	0.69ms
	Average time per request:	2.62ms
	Slowest time for request:	15.55ms

Time:
	0.1%	time for request:	0.88ms
	1%	time for request:	1.24ms
	5%	time for request:	1.39ms
	10%	time for request:	1.49ms
	25%	time for request:	1.95ms
	50%	time for request:	2.31ms
	75%	time for request:	2.98ms
	90%	time for request:	3.71ms
	95%	time for request:	4.78ms
	99%	time for request:	8.03ms
	99.9%	time for request:	12.51ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=8
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	3.27s
	Requests per second:	306081.85
	Fastest time for request:	0.86ms
	Average time per request:	13.28ms
	Slowest time for request:	97.11ms

Time:
	0.1%	time for request:	1.72ms
	1%	time for request:	3.36ms
	5%	time for request:	5.75ms
	10%	time for request:	6.84ms
	25%	time for request:	8.66ms
	50%	time for request:	11.37ms
	75%	time for request:	15.59ms
	90%	time for request:	22.68ms
	95%	time for request:	27.97ms
	99%	time for request:	40.22ms
	99.9%	time for request:	58.74ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=1
```
Summary:
	Clients:	1
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	1.92s
	Requests per second:	521845.22
	Fastest time for request:	0.00ms
	Average time per request:	0.98ms
	Slowest time for request:	30.47ms

Time:
	0.1%	time for request:	0.00ms
	1%	time for request:	0.00ms
	5%	time for request:	0.00ms
	10%	time for request:	0.00ms
	25%	time for request:	0.00ms
	50%	time for request:	0.00ms
	75%	time for request:	0.78ms
	90%	time for request:	3.71ms
	95%	time for request:	5.89ms
	99%	time for request:	9.15ms
	99.9%	time for request:	15.02ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)

```


### def
```
t   Thread
mt  Multi Thread
p   Pipeline
b   Batch
n   Noresponse
```

Batch is only useful when there are multiple goroutines calling it.

### Mac 4 CPU 8 GiB Localhost  qps
```
package     transport   codec       1t      mt      1t_pipe 2t_pipe 1t_bat  2t_bat  1t_no   2t_no   1t_pb   2t_pb   1t_bn   2t_bn   1t_pn   2t_pn   1t_pbn  2t_pbn
MRPC        UDP         protobuf    7219    30725   35001   34850   34671   66934   47850   44730   33035   65859   31548   38715   43091   43378   29318   43364
MRPC        TCP         protobuf    9090    32303   40862   51109   201106  279069  91044   87044   209760  267739  496407  635820  76158   83042   486769  632212
MRPC        WS          protobuf    8041    28174   34658   39505   186921  270323  71917   77072   200805  245743  480326  606449  73493   76328   479299  621587
MRPC        QUIC        protobuf    4127    11258   13207   13445   24557   49931   26405   26531   24425   47796   140531  170448  25867   26912   155083  170559
net/rpc     TCP         gob         9623    31480   44996   51807
net/rpc     HTTP        gob         9602    24343   33435   43911
JSONRPC     TCP         json        9099    28999   43981   48143
GRPC        HTTP2       protobuf    6008    21172   44233   50994
```

### Mac 4 CPU 8 GiB Localhost 99th percentile time (ms)
```
package     transport   codec       1t      mt      1t_pipe 2t_pipe 1t_bat  2t_bat  1t_no   2t_no   1t_pb   2t_pb   1t_bn   2t_bn   1t_pn   2t_pn   1t_pbn  2t_pbn
MRPC        UDP         protobuf    0.21    3.33    1.82    4.13    0.47    0.59    0.33    0.62    0.51    0.63    0.81    0.83    4.57    10.35   0.83    0.88
MRPC        TCP         protobuf    0.17    3.85    1.46    2.83    11.29   21.44   0.13    0.21    8.95    22.15   13.23   27.19   1.29    2.25    13.57   27.40
MRPC        WS          protobuf    0.20    4.88    1.97    4.99    12.67   21.48   0.16    0.24    9.49    24.22   14.16   26.74   0.93    2.25    13.38   26.57
MRPC        QUIC        protobuf    0.42    13.42   5.64    9.23    0.64    0.76    0.21    0.43    0.65    0.92    0.45    0.84    2.52    7.34    0.37    0.90
net/rpc     TCP         gob         0.18    3.78    1.52    3.22
net/rpc     HTTP        gob         0.18    7.76    4.38    5.51
JSONRPC     TCP         json        0.17    4.32    1.58    3.40
GRPC        HTTP2       protobuf    0.27    5.61    2.14    3.92
```
### Linux 12 vCPU 12 GiB Localhost  qps
```
package     transport   codec       1t      mt      1t_pipe 2t_pipe 1t_bat  12t_bat 1t_no  12t_no   1t_bn   2t_bn
MRPC        UDP         protobuf    12452   177727  70023   158784  51998   296501  181198  159354  610008  363127
MRPC        TCP         protobuf    15188   259303  55518   314560  243115  446229  197284  740138  568227  821322
MRPC        WS          protobuf    28583   280422  73847   259468  270826  582654  222924  821275  623432  839396
MRPC        QUIC        protobuf    8848    39362
net/rpc     TCP         gob         20935   275122
net/rpc     HTTP        gob         21415   283631
JSONRPC     TCP         json        19160   224116
GRPC        HTTP2       protobuf    12059   113275
```

### Linux 12 vCPU 12 GiB Between  qps
```
package     transport   codec       1t      mt      1t_pipe 2t_pipe 1t_bat  12t_bat 1t_no   12t_no  1t_bn   2t_bn
MRPC        UDP         protobuf    8456    207902  42920   195943  38454   310601  177132  143414  771688  572685
MRPC        TCP         protobuf    9015    351772  45942   374813  247619  616358  384285  1234356 638546  649993
MRPC        WS          protobuf    8815    381351  51022   298416  259927  640804  279694  832946  734921  892733
MRPC        QUIC        protobuf    5476    49323
net/rpc     TCP         gob         8166    270813
net/rpc     HTTP        gob         8247    293821
JSONRPC     TCP         json        7915    239660
GRPC        HTTP2       protobuf    6325    163704
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