package rpc

import (
	"errors"
)

//Listener defines the interface of listener.
type Listener interface {
	Serve() error
	Addr() string
}

// Listen announces on the local network address.
// The network must be "ipc", "tcp", "udp", "quic", "ws", "http" or "http1".
func Listen(network, address string, server *Server) (Listener, error) {
	logger.Noticef("network - %s", network)
	logger.Noticef("listening on %s", address)
	switch network {
	case IPC:
		return listenIPC(address, server)
	case TCP:
		return listenTCP(address, server)
	case UDP:
		return listenUDP(address, server)
	case QUIC:
		return listenQUIC(address, server)
	case WS:
		return listenWS(address, server)
	case HTTP:
		return listenHTTP(address, server)
	case HTTP1:
		return listenHTTP(address, server)
	}
	return nil, errors.New("this network is not suported")
}
