package rpc

import (
	"errors"
	"hslam.com/git/x/rpc/log"
)

type Listener interface {
	Serve() error
	Addr() string
}

func Listen(network,address string,server *Server) (Listener, error) {
	log.Allf( "network - %s", network)
	log.Allf( "listening on %s", address)
	switch network {
	case TCP:
		return ListenTCP(address,server)
	case UDP:
		return ListenUDP(address,server)
	case QUIC:
		return ListenQUIC(address,server)
	case WS:
		return ListenWS(address,server)
	case HTTP:
		return ListenHTTP(address,server)
	case HTTP1:
		return ListenHTTP(address,server)
	case HTTP2:
		return ListenHTTP2(address,server)
	}
	return nil, errors.New("this network is not suported")
}