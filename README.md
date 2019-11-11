# rpc
A RPC implementation written in Golang over TCP UDP QUIC WS HTTP HTTP2

## Server Feature

* **Network** tcp/ws/http/http2/quic/udp
* **Pipelining** async/sync
* **Multiplexing**

## Client Feature
* **Network** tcp/ws/http/http2/quic
* **Codec** json/protobuf/xml/gob/bytes
* **Compress** flate/zlib/gzip/no
* **Pipelining**
* **Multiplexing**
* **Batch** async/sync
* **Go/Call/CallNoRequest/CallNoResponse/OnlyCall**
* **Protocal** stream/message/frame/conn
* **Pool**

## Benchmark
Batch is only useful when there are multiple goroutines calling it.

### Mac Environment
* **CPU** 4 Cores 2.9 GHz
* **Memory** 8 GiB

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
```
Summary:
	Clients:	1
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	4.85s
	Requests per second:	206070.68
	Fastest time for request:	0.64ms
	Average time per request:	2.48ms
	Slowest time for request:	17.50ms

Time:
	0.1%	time for request:	0.80ms
	1%	time for request:	1.19ms
	5%	time for request:	1.37ms
	10%	time for request:	1.47ms
	25%	time for request:	1.94ms
	50%	time for request:	2.27ms
	75%	time for request:	2.88ms
	90%	time for request:	3.40ms
	95%	time for request:	4.02ms
	99%	time for request:	6.70ms
	99.9%	time for request:	10.88ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=2
```
Summary:
	Clients:	2
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	3.44s
	Requests per second:	291068.62
	Fastest time for request:	0.46ms
	Average time per request:	3.51ms
	Slowest time for request:	20.23ms

Time:
	0.1%	time for request:	0.86ms
	1%	time for request:	1.18ms
	5%	time for request:	1.49ms
	10%	time for request:	1.79ms
	25%	time for request:	2.62ms
	50%	time for request:	3.29ms
	75%	time for request:	4.04ms
	90%	time for request:	5.12ms
	95%	time for request:	6.42ms
	99%	time for request:	9.28ms
	99.9%	time for request:	15.57ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=1
```
Summary:
	Clients:	1
	Parallel calls per client:	512
	Total calls:	1000000
	Total time:	1.85s
	Requests per second:	541076.06
	Fastest time for request:	0.00ms
	Average time per request:	0.94ms
	Slowest time for request:	31.52ms

Time:
	0.1%	time for request:	0.00ms
	1%	time for request:	0.00ms
	5%	time for request:	0.00ms
	10%	time for request:	0.00ms
	25%	time for request:	0.00ms
	50%	time for request:	0.00ms
	75%	time for request:	0.74ms
	90%	time for request:	3.41ms
	95%	time for request:	5.95ms
	99%	time for request:	9.66ms
	99.9%	time for request:	20.41ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
### Linux Environment
* **CPU** 12 Cores 3.1 GHz
* **Memory** 24 GiB

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
```
Summary:
        Clients:        1
        Parallel calls per client:      512
        Total calls:    1000000
        Total time:     3.27s
        Requests per second:    306227.97
        Fastest time for request:       0.43ms
        Average time per request:       1.67ms
        Slowest time for request:       4.39ms

Time:
        0.1%    time for request:       0.61ms
        1%      time for request:       0.88ms
        5%      time for request:       1.19ms
        10%     time for request:       1.34ms
        25%     time for request:       1.53ms
        50%     time for request:       1.59ms
        75%     time for request:       1.72ms
        90%     time for request:       2.15ms
        95%     time for request:       2.48ms
        99%     time for request:       2.90ms
        99.9%   time for request:       3.81ms

Result:
        Response ok:    1000000 (100.00%)
        Errors: 0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=6
```
Summary:
        Clients:        6
        Parallel calls per client:      512
        Total calls:    1000000
        Total time:     1.25s
        Requests per second:    800650.13
        Fastest time for request:       0.20ms
        Average time per request:       3.81ms
        Slowest time for request:       13.27ms

Time:
        0.1%    time for request:       0.67ms
        1%      time for request:       1.12ms
        5%      time for request:       1.82ms
        10%     time for request:       2.19ms
        25%     time for request:       2.81ms
        50%     time for request:       3.52ms
        75%     time for request:       4.50ms
        90%     time for request:       5.79ms
        95%     time for request:       6.82ms
        99%     time for request:       8.87ms
        99.9%   time for request:       10.70ms

Result:
        Response ok:    1000000 (100.00%)
        Errors: 0 (0.00%)
```
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=9999 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=1
```
Summary:
        Clients:        1
        Parallel calls per client:      512
        Total calls:    1000000
        Total time:     1.00s
        Requests per second:    1004994.82
        Fastest time for request:       0.00ms
        Average time per request:       0.50ms
        Slowest time for request:       26.51ms

Time:
        0.1%    time for request:       0.00ms
        1%      time for request:       0.00ms
        5%      time for request:       0.00ms
        10%     time for request:       0.00ms
        25%     time for request:       0.00ms
        50%     time for request:       0.14ms
        75%     time for request:       0.54ms
        90%     time for request:       0.98ms
        95%     time for request:       2.27ms
        99%     time for request:       5.65ms
        99.9%   time for request:       12.01ms

Result:
        Response ok:    1000000 (100.00%)
        Errors: 0 (0.00%)
```


### Define

* **c**     client
* **p**     pipeline
* **m**     multiplex
* **b**     batch
* **n**     no response
* **trans** transport
* **pkg**   package

### Mac 4 CPU 2.9 GHz 8 GiB Localhost  qps
```
pkg     trans   codec   1c      2c      1c_p    2c_p    1c_n    2c_n    1c_pb   2c_pb   1c_m    2c_m    1c_mb   2c_mb   1c_mbn  2c_mbn
RPC     TCP     pb      7640    14899   38460   40417   62296   68947   142992  165612  35846   37042   201432  302815  517066  534081
RPC     WS      pb      7084    13762   36852   38483   58523   68162   147906  189202  32265   35351   201778  295464  541553  538206
RPC     UDP     pb      5136    9739    19807   17739   23525   -       21351   45916   19731   18831   28739   51115   19513   -
RPC     QUIC    pb      3959    7879    13673   12980   27098   26474   22274   46708   13333   12960   23503   51506   152819  151515
RPC     HTTP    pb      3694    7530    4953    8470    5348    8908    130125  163166  4962    8107    149764  228230  301091  690483
RPC     HTTP2   pb      2536    4304    5933    6575    7213    6215    109387  146322  6707    6695    154804  189480  380008  351820
```

### Mac 4 CPU 2.9 GHz 8 GiB Localhost 99th percentile time (ms)
```
pkg     trans   codec   1c      2c      1c_p    2c_p    1c_n    2c_n    1c_pb   2c_pb   1c_m    2c_m    1c_mb   2c_mb   1c_mbn  2c_mbn
RPC     TCP     pb      0.22    0.22    1.91    4.22    0.09    0.21    7.00    13.53   1.66    4.52    5.35    7.92    9.23    22.31
RPC     WS      pb      0.25    0.31    1.71    3.95    0.09    0.18    6.59    10.98   1.96    4.35    5.40    8.04    7.62    29.04
RPC     UDP     pb      0.31    0.42    4.70    17.67   0.42    -       1.35    0.93    5.56    14.92   0.59    0.75    0.77    -
RPC     QUIC    pb      0.38    0.58    3.54    8.92    0.19    0.50    0.65    0.65    3.73    8.28    0.64    0.60    0.36    1.04
RPC     HTTP    pb      0.47    0.53    8.73    10.30   0.38    0.56    7.73    13.27   9.30    11.75   7.48    13.02   7.91    31.05
RPC     HTTP2   pb      0.68    1.07    13.63   19.82   3.07    8.64    9.81    14.35   9.65    20.36   7.70    11.72   22.76   87.47
```

### Linux 12 CPU 3.1 GHz 24 GiB Localhost  qps
```
pkg     trans   codec   1c      8c      1c_p    8c_p    1c_n    8c_n    1c_pb   8c_pb   1c_m    8c_m    1c_mb   8c_mb   1c_mbn  8c_mbn
RPC     TCP     pb      22509   118317  105188  314515  213503  975411  286945  581101  100365  269593  304154  763328  1029544 1150087
RPC     WS      pb      19873   99763   90240   221460  184139  904936  279640  584942  91887   205025  291581  758001  953977  1063996
RPC     UDP     pb      15931   72314   89083   140257  6052    -       60258   233585  86391   138741  70512   278373  20220   -
RPC     QUIC    pb      9999    53056   39473   55485   68558   108047  44110   213551  39240   55671   52343   258951  468606  566404
RPC     HTTP    pb      10665   54212   13395   52496   15519   63847   263006  542608  13304   52317   275541  780726  456321  1067055
RPC     HTTP2   pb      5940    29049   18884   47305   22652   64481   241620  546947  17792   53155   253349  723360  941024  967815
```

### Linux 12 CPU 3.1 GHz 24 GiB Localhost 99th percentile time (ms)
```
pkg     trans   codec   1c      8c      1c_p    8c_p    1c_n    8c_n    1c_pb   8c_pb   1c_m    8c_m    1c_mb   8c_mb   1c_mbn  8c_mbn
RPC     TCP     pb      0.05    0.12    0.51    2.83    0.05    0.03    3.16    23.19   0.59    3.48    2.92    12.86   6.68    57.54
RPC     WS      pb      0.07    0.29    0.65    4.44    0.05    0.06    3.30    15.91   0.69    4.21    3.10    13.00   4.43    34.71
RPC     UDP     pb      0.09    0.22    0.81    5.03    0.03    -       0.24    0.70    0.82    5.14    0.21    0.65    0.12    -
RPC     QUIC    pb      0.14    0.59    1.35    7.46    0.06    0.65    0.33    1.24    1.38    7.44    0.27    0.94    0.10    1.48
RPC     HTTP    pb      0.13    0.35    3.01    7.15    0.10    0.36    3.50    32.19   2.99    6.94    3.31    13.12   2.58    55.59
RPC     HTTP2   pb      0.22    0.75    2.87    13.19   1.11    4.13    3.79    18.37   3.08    12.01   3.69    13.28   4.32    55.87
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
}
```
### arith.pb.go
```
protoc ./arith.proto --go_out=./
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
	"hslam.com/git/x/rpc/examples/helloworld/pb/service"
	"hslam.com/git/x/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("tcp",":9999")//tcp|ws|http|http2|quic
}
```

### client.go
```go
package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/pb/service"
	"hslam.com/git/x/rpc"
	"fmt"
	"log"
)
func main()  {
	conn, err:= rpc.Dial("tcp","127.0.0.1:9999","pb")//tcp|ws|http|http2|quic
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

## Example http json javascript

### http-json-server.go
```go
package main
import (
	"hslam.com/git/x/rpc/examples/helloworld/json_js/service"
	"hslam.com/git/x/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("http",":9999")
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
</script>
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)


### Authors
rpc was written by Mort Huang.