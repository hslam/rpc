#!/bin/sh

sh ./command.sh tcp >> ./log.out/tcp.log.out
sh ./command.sh ws >> ./log.out/ws.log.out
sh ./command.sh quic >> ./log.out/quic.log.out
sh ./command.sh http >> ./log.out/http.log.out
sh ./command.sh http2 >> ./log.out/http2.log.out
