### BENCHMARK
go test -bench=. -benchmem
```
goos: darwin
goarch: amd64
pkg: github.com/hslam/fastrpc/codec
BenchmarkRequestMarshalCODE-4      	33914080	        34.9 ns/op	15373.65 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestMarshalGOGOPB-4    	23503790	        51.5 ns/op	10368.83 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestMarshalJSON-4      	  750686	      1446 ns/op	 499.14 MB/s	    1472 B/op	       2 allocs/op
BenchmarkRequestUnmarshalCODE-4    	51390728	        22.5 ns/op	23852.02 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestUnmarshalPB-4      	16403053	        71.6 ns/op	7454.04 MB/s	      16 B/op	       1 allocs/op
BenchmarkRequestUnmarshalJSON-4    	  177550	      6656 ns/op	 108.47 MB/s	     800 B/op	       6 allocs/op
BenchmarkRequestRoundtripCODE-4    	20193453	        58.1 ns/op	9246.32 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestRoundtripPB-4      	 9720596	       122 ns/op	4393.69 MB/s	      16 B/op	       1 allocs/op
BenchmarkRequestRoundtripJSON-4    	  141364	      8214 ns/op	  87.90 MB/s	    2272 B/op	       8 allocs/op
BenchmarkResponseMarshalCODE-4     	35888379	        33.5 ns/op	15589.77 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseMarshalGOGOPB-4   	28574744	        41.6 ns/op	12460.89 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseMarshalJSON-4     	  814357	      1422 ns/op	 498.00 MB/s	    1472 B/op	       2 allocs/op
BenchmarkResponseUnmarshalCODE-4   	51008488	        22.7 ns/op	23015.73 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseUnmarshalPB-4     	31294515	        38.1 ns/op	13605.07 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseUnmarshalJSON-4   	  181924	      6473 ns/op	 109.38 MB/s	     784 B/op	       5 allocs/op
BenchmarkResponseRoundtripCODE-4   	20749662	        56.5 ns/op	9248.95 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseRoundtripPB-4     	14576751	        80.8 ns/op	6408.26 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseRoundtripJSON-4   	  143953	      8159 ns/op	  86.78 MB/s	    2256 B/op	       7 allocs/op
BenchmarkRoundtripCODE-4           	 9870566	       119 ns/op	4386.88 MB/s	       0 B/op	       0 allocs/op
BenchmarkRoundtripPB-4             	 5782466	       206 ns/op	2519.20 MB/s	      16 B/op	       1 allocs/op
BenchmarkRoundtripJSON-4           	   72177	     16584 ns/op	  42.69 MB/s	    4529 B/op	      15 allocs/op
PASS
ok  	github.com/hslam/fastrpc/codec	31.270s
```