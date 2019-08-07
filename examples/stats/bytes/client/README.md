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
MRPC        FASTHTTP    protobuf    8925    25883   22977   23177   199272  246933  24114   25444   175341  264204  475338  504053  20526   24167   530874  560745
MRPC        QUIC        protobuf    4127    11258   13207   13445   24557   49931   26405   26531   24425   47796   140531  170448  25867   26912   155083  170559
RPC         TCP         gob         8373    27192
RPC         HTTP        gob         8341    27678
JSONRPC     TCP         json        7809    25390
GRPC        HTTP2       protobuf    5529    18405
```

### MAC 99th percentile time  (ms)
```
package     Protocol    encoding    1t      mt      1t_con  2t_con  1t_bat  2t_bat  1t_no   2t_no   1t_cb   2t_cb   1t_bn   2t_bn   1t_cn   2t_cn   1t_cbn  2t_cbn
MRPC        UDP         protobuf    0.21    3.33    1.82    4.13    0.47    0.59    0.33    0.62    0.51    0.63    0.81    0.83    4.57    10.35   0.83    0.88
MRPC        TCP         protobuf    0.17    3.85    1.46    2.83    11.29   21.44   0.13    0.21    8.95    22.15   13.23   27.19   1.29    2.25    13.57   27.40
MRPC        WS          protobuf    0.20    4.88    1.97    4.99    12.67   21.48   0.16    0.24    9.49    24.22   14.16   26.74   0.93    2.25    13.38   26.57
MRPC        FASTHTTP    protobuf    0.19    8.48    2.56    6.60    11.17   23.21   0.35    0.74    11.94   21.75   8.81    17.60   3.90    5.84    7.15    17.48
MRPC        QUIC        protobuf    0.42    13.42   5.64    9.23    0.64    0.76    0.21    0.43    0.65    0.92    0.45    0.84    2.52    7.34    0.37    0.90
RPC         TCP         gob         8373    27192
RPC         HTTP        gob         8341    27678
JSONRPC     TCP         json        7809    25390
GRPC        HTTP2       protobuf    5529    18405
```
