#!/bin/sh

net="$1"
host="127.0.0.1"
pbar="false"
cod="pb"
com="no"
p="8080"
t1="10000"
t2="12000"
t3="100000"
t4="120000"
t5="250000"
t6="300000"
t7="500000"
t8="600000"
c="2"

sleep 1s
nohup ./rpc_server -network=$net -p=$p -async=false -pipelining=false -multiplexing=false -batch=false > ./tmp/log.rpc_server &
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t1 -pipelining=false -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=1
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t2 -pipelining=false -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=$c
killall rpc_server

sleep 1s
nohup ./rpc_server -network=$net -p=$p -async=false -pipelining=true -multiplexing=false -batch=false > ./tmp/log.rpc_server &
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t3 -pipelining=true -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=1
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=true -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=$c
killall rpc_server

sleep 1s
nohup ./rpc_server -network=$net -p=$p -async=false -pipelining=false -multiplexing=true -batch=false > ./tmp/log.rpc_server &
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t3 -pipelining=false -multiplexing=true -batch=false -batch_async=false -noresponse=false -clients=1
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=false -multiplexing=true -batch=false -batch_async=false -noresponse=false -clients=$c
killall rpc_server

sleep 1s
nohup ./rpc_server -network=$net -p=$p -async=false -pipelining=true -multiplexing=false -batch=true > ./tmp/log.rpc_server &
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t5 -pipelining=true -multiplexing=false -batch=true -batch_async=false -noresponse=false -clients=1
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t6 -pipelining=true -multiplexing=false -batch=true -batch_async=false -noresponse=false -clients=$c
killall rpc_server


sleep 1s
nohup ./rpc_server -network=$net -p=$p -async=false -pipelining=false -multiplexing=true -batch=true > ./tmp/log.rpc_server &
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t5 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t6 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=$c
killall rpc_server


sleep 1s
nohup ./rpc_server -network=$net -p=$p -async=false -pipelining=false -multiplexing=true -batch=true > ./tmp/log.rpc_server &
sleep 1s
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t7 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=1
sleep 1s
if [ $net != "udp" ] ; then
./rpc_client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t8 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=$c
fi
killall rpc_server


