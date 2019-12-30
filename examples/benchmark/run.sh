#!/bin/sh

path="log.out"

if [ -d $path ]
then
        rm -rf $path
fi
mkdir $path

echo "Start Build"

go build -o rpc_server ./tcp/pb/server/server.go
go build -o rpc_client ./tcp/pb/client/client.go

echo "Start Benchmark"

echo "TCP Benchmark ..."

sh ./command.sh tcp > ./log.out/tcp.log.out

echo "QUIC Benchmark ..."

sh ./command.sh quic > ./log.out/quic.log.out

echo "HTTP Benchmark ..."

sh ./command.sh http > ./log.out/http.log.out

echo "WS Benchmark ..."

sh ./command.sh ws > ./log.out/ws.log.out

rm rpc_server rpc_client
rm -r tmp
echo ""
echo "Benchmark Finish!"
echo "Please Enable CTRL+C to Exit"