```
go-torch -u http://localhost:6060/debug/pprof/ -p > rpc-tcp-gob-1-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > rpc-tcp-gob-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > rpc-tcp-gob-1-heap.svg

go-torch -u http://localhost:6060/debug/pprof/ -p > rpc-tcp-gob-32-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > rpc-tcp-gob-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > rpc-tcp-gob-32-heap.svg


go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-tcp-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-tcp-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-tcp-pb-1-heap.svg

go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-tcp-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-tcp-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-tcp-pb-32-heap.svg


go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-quic-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-quic-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-quic-pb-1-heap.svg

go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-quic-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-quic-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-quic-pb-32-heap.svg


go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-udp-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-udp-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-udp-pb-1-heap.svg

go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-udp-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-udp-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-udp-pb-32-heap.svg


go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-ws-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-ws-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-ws-pb-1-heap.svg

go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-ws-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-ws-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-ws-pb-32-heap.svg


go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-fast-pb-1-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-fast-pb-1-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-fast-pb-1-heap.svg

go-torch -u http://localhost:6060/debug/pprof/ -p > mrpc-fast-pb-32-profile-framegraph.svg
go-torch -u http://localhost:6060/debug/pprof/heap -p > mrpc-fast-pb-32-heap-framegraph.svg
go tool pprof -alloc_space -cum -svg http://127.0.0.1:6060/debug/pprof/heap > mrpc-fast-pb-32-heap.svg
```


