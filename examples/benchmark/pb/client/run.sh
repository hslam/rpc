#!/bin/sh

sh ./command.sh tcp >> tcp.log.out
sh ./command.sh ws >> ws.log.out
sh ./command.sh quic >> quic.log.out
sh ./command.sh http >> http.log.out
sh ./command.sh http2 >> http2.log.out
sh ./command.sh udp >> udp.log.out
