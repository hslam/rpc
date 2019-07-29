package rpc

import (
	"errors"
)

type Listener interface {
	Serve() error
	Addr() string
}

func Listen(network,address string) (Listener, error) {
	Allf( "network - %s", network)
	Allf( "listening on %s", address)
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
	}
	return nil, errors.New("this network is not suported")
}