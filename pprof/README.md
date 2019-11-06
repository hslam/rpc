```
go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-tcp-gob-1-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-tcp-gob-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-tcp-gob-1-heap.svg

go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-tcp-gob-32-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-tcp-gob-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-tcp-gob-32-heap.svg


go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-tcp-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-tcp-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-tcp-pb-1-heap.svg

go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-tcp-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-tcp-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-tcp-pb-32-heap.svg


go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-quic-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-quic-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-quic-pb-1-heap.svg

go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-quic-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-quic-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-quic-pb-32-heap.svg


go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-udp-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-udp-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-udp-pb-1-heap.svg

go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-udp-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-udp-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-udp-pb-32-heap.svg


go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-ws-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-ws-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-ws-pb-1-heap.svg

go-torch -u http://localhost:6061/debug/pprof/ -p > rpc-ws-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6061/debug/pprof/heap -p > rpc-ws-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6061/debug/pprof/heap > rpc-ws-pb-32-heap.svg
```


