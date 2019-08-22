package rpc

import (
	"hslam.com/mgit/Mort/rpc/log"
	"net/http"
	"golang.org/x/net/http2"
	"net"
)

type HTTP2Listener struct {
	address			string
	netListener		net.Listener
	server			http.Server
}
func ListenHTTP2(address string) (Listener, error) {
	var server http.Server
	server.Addr = address
	server.TLSConfig=DefalutTLSConfig()
	http.HandleFunc("/", IndexHandler)
	s2 := &http2.Server{}
	http2.ConfigureServer(&server, s2)
	var netListener net.Listener
	var err error
	if netListener, err = net.Listen("tcp", address); err != nil {
		return nil,err
	}
	listener:=  &HTTP2Listener{address:address,netListener:netListener,server:server}
	return listener,nil
}

func (l *HTTP2Listener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	err:=l.server.ServeTLS(l.netListener,"","")
	if err!=nil{
		log.Fatalf("fatal error: %s", err)
	}
	return nil
}
func (l *HTTP2Listener)Addr() (string) {
	return l.address
}
