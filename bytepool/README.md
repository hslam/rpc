### go test -v -run="none" -bench=. -benchtime=3s  -benchmem
```
goos: darwin
goarch: amd64
pkg: hslam.com/mgit/Mort/rpc/bytepool
BenchmarkBytePool-4   	50000000	        71.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkByteMake-4   	 1000000	      5592 ns/op	   65536 B/op	       1 allocs/op
PASS
ok  	hslam.com/mgit/Mort/rpc/bytepool	12.217s
```