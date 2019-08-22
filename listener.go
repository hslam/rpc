package rpc

import (
	"errors"
	"hslam.com/mgit/Mort/rpc/log"
)

type Listener interface {
	Serve() error
	Addr() string
}

func Listen(network,address string) (Listener, error) {
	log.Allf( "network - %s", network)
	log.Allf( "listening on %s", address)
	switch network {
	case TCP:
		return ListenTCP(address)
	case UDP:
		return ListenUDP(address)
	case QUIC:
		return ListenQUIC(address)
	case WS:
		return ListenWS(address)
	case FASTHTTP:
		return ListenFASTHTTP(address)
	case HTTP:
		return ListenHTTP(address)
	case HTTP2:
		return ListenHTTP2(address)
	}
	return nil, errors.New("this network is not suported")
}