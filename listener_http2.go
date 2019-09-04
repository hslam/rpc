package rpc

import (
	"hslam.com/mgit/Mort/rpc/log"
	"net/http"
	"golang.org/x/net/http2"
	"net"
)

type HTTP2Listener struct {
	server			*Server
	address			string
	netListener		net.Listener
	httpServer			http.Server
}
func ListenHTTP2(address string,server *Server) (Listener, error) {
	var httpServer http.Server
	handler:=new(Handler)
	handler.server=server
	httpServer.Handler=handler
	httpServer.Addr = address
	httpServer.TLSConfig=DefalutTLSConfig()
	s2 := &http2.Server{}
	http2.ConfigureServer(&httpServer, s2)
	var netListener net.Listener
	var err error
	if netListener, err = net.Listen("tcp", address); err != nil {
		return nil,err
	}
	listener:=  &HTTP2Listener{address:address,netListener:netListener,httpServer:httpServer,server:server}
	return listener,nil
}

func (l *HTTP2Listener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	err:=l.httpServer.ServeTLS(l.netListener,"","")
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return err
	}
	return nil
}
func (l *HTTP2Listener)Addr() (string) {
	return l.address
}
