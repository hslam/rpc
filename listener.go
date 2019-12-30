package rpc

import (
	"errors"
)

type Listener interface {
	Serve() error
	Addr() string
}

func Listen(network,address string,server *Server) (Listener, error) {
	Allf( "network - %s", network)
	Allf( "listening on %s", address)
	switch network {
	case IPC:
		return ListenIPC(address,server)
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
	}
	return nil, errors.New("this network is not suported")
}