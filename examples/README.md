
## Benchmark
Batch is only useful when there are multiple goroutines calling it.

### Mac Environment
* **CPU** 4 Cores 2.9 GHz
* **Memory** 8 GiB

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
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
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=2
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
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=1
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

./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
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
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=6
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
./client -network=tcp -codec=pb -compress=no -h=127.0.0.1 -p=8080 -total=1000000 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=1
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

### Mac 4 CPU 2.9 GHz 8 GiB Requests per second
```
pkg     trans  codec 1c      2c      1c_p    2c_p    1c_m    2c_m    1c_pb   2c_pb   1c_mb   2c_mb   1c_mbn
RPC     TCP    pb    7934    15615   71117   92935   69100   79206   167454  205973  218024  312157  586399
RPC     WS     pb    6839    13787   32979   33019   29890   34117   144973  199951  180127  180127  581486
RPC     QUIC   pb    3869    7140    42769   52329   46133   50490   55937   95972   59992   114576  367816
RPC     HTTP   pb    7660    14944   71043   93771   68557   80717   167123  220909  209429  315826  580307
RPC     HTTP1  pb    4721    9209    5560    11343   5559    11540   159709  214649  200784  292199  311868
RPC     HTTP2  pb    2729    5176    8080    6636    8009    6986    140996  189238  181252  267546  529664
NETRPC  TCP    gob   7891    16635   -       -       47213   54125   -       -       -       -       -
NETRPC  HTTP   gob   8494    16033   -       -       51541   59532   -       -       -       -       -
JSONRPC HTTP   json  8101    14661   -       -       39201   43342   -       -       -       -       -
GRPC    HTTP2  pb    5877    9814    -       -       45546   54089   -       -       -       -       -
```

### Mac 4 CPU 2.9 GHz 8 GiB 99th percentile time (ms)
```
pkg     trans  codec 1c      2c      1c_p    2c_p    1c_m    2c_m    1c_pb   2c_pb   1c_mb   2c_mb   1c_mbn
RPC     TCP    pb    0.21    0.20    0.98    1.51    0.89    1.73    5.61    10.40   4.63    7.90    10.72
RPC     WS     pb    0.34    0.27    2.31    5.56    2.39    4.77    8.06    11.87   5.56    10.02   10.09
RPC     QUIC   pb    0.46    0.68    2.01    4.27    1.61    3.59    1.12    1.51    1.02    1.17    1.01
RPC     HTTP   pb    0.20    0.21    0.89    1.47    0.94    1.67    5.75    10.28   5.24    7.70    11.26
RPC     HTTP1  pb    0.30    0.34    8.15    7.71    8.34    7.59    6.38    10.13   5.05    8.01    4.47
RPC     HTTP2  pb    0.62    0.80    8.36    22.49   8.83    18.95   6.54    10.94   5.39    8.61    20.45
NETRPC  TCP    gob   0.26    0.21    -       -       1.48    2.99    -       -       -       -       -
NETRPC  HTTP   gob   0.22    0.22    -       -       5.23    11.97   -       -       -       -       -
JSONRPC HTTP   json  0.21    0.28    -       -       2.15    4.45    -       -       -       -       -
GRPC    HTTP2  pb    0.29    0.51    -       -       1.65    2.90    -       -       -       -       -
```

### Linux 12 CPU 3.1 GHz 24 GiB Requests per second
```
pkg     trans  codec 1c      8c      1c_p    8c_p    1c_m    8c_m    1c_pb   8c_pb   1c_mb   8c_mb   1c_mbn
RPC     TCP    pb    22080   119183  163335  361216  154558  313828  288621  578407  305782  777963  1181147
RPC     WS     pb    20109   100620  90610   222028  92005   204300  283193  581244  296610  761963  1035882
RPC     QUIC   pb    10011   53606   111137  251736  102591  219560  105142  426647  111252  492597  904398
RPC     HTTP   pb    21730   118649  162478  366485  154139  328157  287788  582813  305188  834926  1210480
RPC     HTTP1  pb    11415   58564   14238   56618   14127   56329   259239  576407  265513  801691  468574
RPC     HTTP2  pb    6276    30207   12378   55507   12357   55488   248576  553502  248488  747352  1119258
NETRPC  TCP    gob   30200   137407  -       -       94253   351391  -       -       -       -       -
NETRPC  HTTP   gob   29917   138912  -       -       95735   393183  -       -       -       -       -
JSONRPC HTTP   json  26664   123223  -       -       92579   294136  -       -       -       -       -
GRPC    HTTP2  pb    15386   61854   -       -       89989   145862  -       -       -       -       -
```

### Linux 12 CPU 3.1 GHz 24 GiB 99th percentile time (ms)
```
pkg     trans  codec 1c      8c      1c_p    8c_p    1c_m    8c_m    1c_pb   8c_pb   1c_mb   8c_mb   1c_mbn
RPC     TCP    pb    0.06    0.12    0.45    1.94    0.51    2.04    3.31    25.40   2.92    13.04   7.62
RPC     WS     pb    0.07    0.28    0.64    4.26    0.69    4.15    3.37    16.67   3.03    13.12   5.19
RPC     QUIC   pb    0.14    0.58    0.72    2.74    0.81    2.88    0.54    1.76    0.51    1.56    0.42
RPC     HTTP   pb    0.06    0.12    0.44    1.92    0.47    1.95    3.21    32.97   2.93    12.18   7.55
RPC     HTTP1  pb    0.12    0.33    2.81    6.60    2.83    6.64    3.54    21.53   3.45    13.14   2.73
RPC     HTTP2  pb    0.21    0.73    4.27    11.40   4.24    11.34   3.65    17.53   3.61    13.11   7.82
NETRPC  TCP    gob   0.04    0.10    -       -       0.61    2.33    -       -       -       -       -
NETRPC  HTTP   gob   0.04    0.10    -       -       2.52    7.77   -       -       -       -       -
JSONRPC HTTP   json  0.04    0.12    -       -       0.61    2.84    -       -       -       -       -
GRPC    HTTP2  pb    0.11    0.42    -       -       0.80    4.81    -       -       -       -       -
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
	"github.com/hslam/rpc/examples/helloworld/pb/service"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("tcp",":8080")//tcp|ws|quic|http|http1|http2
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
	conn, err:= rpc.Dial("tcp","127.0.0.1:8080","pb")//tcp|ws|quic|http|http1|http2
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
	transport:=rpc.NewTransport(MaxConnsPerHost,MaxIdleConnsPerHost,"tcp","pb",rpc.DefaultOptions())
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

## Example http json javascript

### http-json-server.go
```go
package main
import (
	"github.com/hslam/rpc/examples/helloworld/json_js/service"
	"github.com/hslam/rpc"
)
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe("http1",":8080")
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
    var client = new rpc.Dial("127.0.0.1:8080");
    var req = new ArithRequest(9,2)
    var reply=client.Call("Arith.Multiply",req)
    console.log(req.a.toString()+" * "+req.b.toString()+" = "+reply.pro.toString());
</script>
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Meng Huang)


### Authors
rpc was written by Meng Huang.