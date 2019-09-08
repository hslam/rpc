### go test -v -run="none" -bench=. -benchtime=3s  -benchmem
```
goos: darwin
goarch: amd64
pkg: hslam.com/mgit/Mort/rpc/protocol
BenchmarkRTO-4             	300000000	        11.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint16ToBytes-4   	5000000000	         0.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkBytesToUint16-4   	1000000000	         8.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint32ToBytes-4   	5000000000	         0.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkBytesToUint32-4   	500000000	        10.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket64-4        	100000000	        64.2 ns/op	      80 B/op	       1 allocs/op
BenchmarkUnpack64-4        	30000000	       159 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket512-4       	30000000	       193 ns/op	     576 B/op	       1 allocs/op
BenchmarkUnpack512-4       	30000000	       162 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket1500-4      	10000000	       384 ns/op	    1536 B/op	       1 allocs/op
BenchmarkUnpack1500-4      	30000000	       165 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket65536-4     	  300000	     14084 ns/op	   73728 B/op	       1 allocs/op
BenchmarkUnpack65536-4     	30000000	       150 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	hslam.com/mgit/Mort/rpc/protocol	68.890s
```
