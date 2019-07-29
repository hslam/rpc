### go test -v -run="none" -bench=. -benchtime=3s  -benchmem
```
goos: darwin
goarch: amd64
pkg: hslam.com/mgit/Mort/transporter/protocol
BenchmarkRTO-4             	500000000	         8.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint16ToBytes-4   	2000000000	         0.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkBytesToUint16-4   	500000000	         3.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint32ToBytes-4   	2000000000	         0.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkBytesToUint32-4   	200000000	         5.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket64-4        	30000000	        42.3 ns/op	      80 B/op	       1 allocs/op
BenchmarkUnpack64-4        	20000000	        76.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket512-4       	10000000	       121 ns/op	     576 B/op	       1 allocs/op
BenchmarkUnpack512-4       	20000000	        74.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket1500-4      	 5000000	       277 ns/op	    1536 B/op	       1 allocs/op
BenchmarkUnpack1500-4      	20000000	        76.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPacket65536-4     	  200000	      7782 ns/op	   73728 B/op	       1 allocs/op
BenchmarkUnpack65536-4     	20000000	        75.2 ns/op	       0 B/op	       0 allocs/op
```
