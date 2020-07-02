### BENCHMARK
go test -bench=. -benchmem
```
goos: darwin
goarch: amd64
pkg: github.com/hslam/rpc
BenchmarkRequestMarshalCODE-4        	33482854	        31.1 ns/op	17055.25 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestMarshalCODEPB-4      	34093851	        34.6 ns/op	15421.15 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestMarshalGOGOPB-4      	17655142	        66.6 ns/op	8022.43 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestMarshalJSON-4        	  894202	      1247 ns/op	 579.00 MB/s	    1472 B/op	       2 allocs/op
BenchmarkRequestUnmarshalCODE-4      	46397149	        25.0 ns/op	21232.64 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestUnmarshalCODEPB-4    	38636697	        30.6 ns/op	17476.32 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestUnmarshalGOGOPB-4    	16259637	        71.6 ns/op	7453.41 MB/s	      16 B/op	       1 allocs/op
BenchmarkRequestUnmarshalJSON-4      	  183936	      6344 ns/op	 113.81 MB/s	     800 B/op	       6 allocs/op
BenchmarkRequestRoundtripCODE-4      	19764566	        58.9 ns/op	9017.95 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestRoundtripCODEPB-4    	17977524	        64.9 ns/op	8226.33 MB/s	       0 B/op	       0 allocs/op
BenchmarkRequestRoundtripGOGOPB-4    	 8483186	       140 ns/op	3807.32 MB/s	      16 B/op	       1 allocs/op
BenchmarkRequestRoundtripJSON-4      	  149979	      7850 ns/op	  91.97 MB/s	    2272 B/op	       8 allocs/op
BenchmarkResponseMarshalCODE-4       	43076028	        27.8 ns/op	18626.74 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseMarshalCODEPB-4     	41994508	        27.7 ns/op	18668.33 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseMarshalGOGOPB-4     	20898134	        56.2 ns/op	9221.58 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseMarshalJSON-4       	  972465	      1225 ns/op	 577.77 MB/s	    1472 B/op	       2 allocs/op
BenchmarkResponseUnmarshalCODE-4     	57120253	        20.3 ns/op	25487.98 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseUnmarshalCODEPB-4   	52179763	        22.8 ns/op	22700.75 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseUnmarshalGOGOPB-4   	31025642	        39.5 ns/op	13116.48 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseUnmarshalJSON-4     	  166291	      6240 ns/op	 113.46 MB/s	     784 B/op	       5 allocs/op
BenchmarkResponseRoundtripCODE-4     	23979372	        49.2 ns/op	10497.67 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseRoundtripCODEPB-4   	22384860	        52.3 ns/op	9901.09 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseRoundtripGOGOPB-4   	12373190	        95.2 ns/op	5439.52 MB/s	       0 B/op	       0 allocs/op
BenchmarkResponseRoundtripJSON-4     	  149359	      7613 ns/op	  93.00 MB/s	    2256 B/op	       7 allocs/op
BenchmarkRoundtripCODE-4             	11107009	       105 ns/op	9985.33 MB/s	       0 B/op	       0 allocs/op
BenchmarkRoundtripCODEPB-4           	10183730	       116 ns/op	9085.41 MB/s	       0 B/op	       0 allocs/op
BenchmarkRoundtripGOGOPB-4           	 4641254	       259 ns/op	4063.59 MB/s	      16 B/op	       1 allocs/op
BenchmarkRoundtripJSON-4             	   76912	     15405 ns/op	  92.83 MB/s	    4529 B/op	      15 allocs/op
PASS
ok  	github.com/hslam/rpc	34.985s
```