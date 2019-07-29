### def
```
t   Thread
mt  Multi Thread
c   Concurrent
b   Batch
n   Noresponse
```
### MAC qps
```
package     Protocol    encoding    1t      mt      1t_con  2t_con  1t_bat  2t_bat  1t_no   2t_no   1t_cb   2t_cb   1t_bn   2t_bn   1t_cn   2t_cn   1t_cbn  2t_cbn
MRPC        UDP         protobuf    7219    30725   35001   34850   34671   66934   47850   44730   33035   65859   31548   38715   43091   43378   29318   43364
MRPC        TCP         protobuf    9090    32303   40862   51109   201106  279069  91044   87044   209760  267739  496407  635820  76158   83042   486769  632212
MRPC        WS          protobuf    8041    28174   34658   39505   186921  270323  71917   77072   200805  245743  480326  606449  73493   76328   479299  621587
MRPC        FASTHTTP    protobuf    9112    27224   25831   26952   191358  246002  28638   29242   191255  248876  567767  596883  28392   28392   557057  592714
MRPC        QUIC        protobuf    4127    11258   13207   13445   24557   49931   26405   26531   24425   47796   140531  170448  25867   26912   155083  170559
RPC         TCP         gob         9623    31480   44996   51807
RPC         HTTP        gob         9602    24343   33435   43911
JSONRPC     TCP         json        9099    28999   43981   48143
GRPC        HTTP2       protobuf    6008    21172   44233   50994
```

### MAC 99th percentile time  (ms)
```
package     Protocol    encoding    1t      mt      1t_con  2t_con  1t_bat  2t_bat  1t_no   2t_no   1t_cb   2t_cb   1t_bn   2t_bn   1t_cn   2t_cn   1t_cbn  2t_cbn
MRPC        UDP         protobuf    0.21    3.33    1.82    4.13    0.47    0.59    0.33    0.62    0.51    0.63    0.81    0.83    4.57    10.35   0.83    0.88
MRPC        TCP         protobuf    0.17    3.85    1.46    2.83    11.29   21.44   0.13    0.21    8.95    22.15   13.23   27.19   1.29    2.25    13.57   27.40
MRPC        WS          protobuf    0.20    4.88    1.97    4.99    12.67   21.48   0.16    0.24    9.49    24.22   14.16   26.74   0.93    2.25    13.38   26.57
MRPC        FASTHTTP    protobuf    0.18    7.36    2.71    6.30    10.83   23.81   1.03    2.03    12.87   26.91   15.08   35.46   3.23    7.66   15.63    37.57
MRPC        QUIC        protobuf    0.42    13.42   5.64    9.23    0.64    0.76    0.21    0.43    0.65    0.92    0.45    0.84    2.52    7.34    0.37    0.90
RPC         TCP         gob         0.18    3.78    1.52    3.22
RPC         HTTP        gob         0.18    7.76    4.38    5.51
JSONRPC     TCP         json        0.17    4.32    1.58    3.40
GRPC        HTTP2       protobuf    0.27    5.61    2.14    3.92
```

### server
```go
package main
import (
	"hslam.com/mgit/Mort/rpc/examples/helloworld/pb/service"
	"hslam.com/mgit/Mort/rpc"
	"strconv"
	"flag"
)
var network string
var port int
var saddr string
func init()  {
	flag.StringVar(&network, "network", "tcp", "network: -network=fast;ws;tcp;quic;udp")
	flag.IntVar(&port, "p", 9999, "port: -p=9999")
	flag.Parse()
	saddr = ":"+strconv.Itoa(port)
}
func main()  {
	rpc.Register(new(service.Arith))
	rpc.ListenAndServe(network,saddr)
}
```

### client
```go
package main
import (
	"hslam.com/mgit/Mort/rpc/examples/helloworld/pb/service"
	"hslam.com/mgit/Mort/rpc"
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
	flag.StringVar(&network, "network", "tcp", "network: -network=fast|ws|tcp|quic|udp")
	flag.StringVar(&codec, "codec", "pb", "codec: -codec=pb|json|xml")
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