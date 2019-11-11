#!/bin/sh

net="$1"
host="127.0.0.1"
pbar="false"
cod="pb"
com="no"
p="9999"
t1="50000"
t2="100000"
t3="500000"
t4="1000000"

nohup ./server -network=$net -async=false -multiplexing=false > log.server &
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t1 -pipelining=false -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=1
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t2 -pipelining=false -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=8

sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t3 -pipelining=true -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=1
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=true -multiplexing=false -batch=false -batch_async=false -noresponse=false -clients=8

sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t2 -pipelining=false -multiplexing=false -batch=false -batch_async=false -noresponse=true -clients=1
sleep 3s
if [ $net != "udp" ]; then
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t3 -pipelining=false -multiplexing=false -batch=false -batch_async=false -noresponse=true -clients=8
fi
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t3 -pipelining=true -multiplexing=false -batch=true -batch_async=false -noresponse=false -clients=1
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=true -multiplexing=false -batch=true -batch_async=false -noresponse=false -clients=8
sleep 3s

killall server 
sleep 3s
nohup ./server -network=$net -async=false -multiplexing=true > log.server &
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t3 -pipelining=false -multiplexing=true -batch=false -batch_async=false -noresponse=false -clients=1
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=false -multiplexing=true -batch=false -batch_async=false -noresponse=false -clients=8

sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t3 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=1
sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=false -clients=8

sleep 3s
./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=1
sleep 3s
if [ $net != "udp" ] ; then
	./client -network=$net -codec=$cod -compress=$com -h=$host -p=$p -total=$t4 -pipelining=false -multiplexing=true -batch=true -batch_async=false -noresponse=true -clients=8
fi
sleep 3s
killall server 


