package rpc

import (
	"errors"
)

//Conn defines the interface of conn.
type Conn interface {
	Handle(readChan chan []byte, writeChan chan []byte, stopChan chan bool, finishChan chan bool)
	NoDelay(enable bool)
	Multiplexing(enable bool)
	TickerFactor() int
	BatchFactor() int
	Retry() error
	Close() error
	Closed() bool
}

func dial(network, address string) (Conn, error) {
	switch network {
	case IPC:
		return dialIPC(address)
	case TCP:
		return dialTCP(address)
	case UDP:
		return dialUDP(address)
	case QUIC:
		return dialQUIC(address)
	case WS:
		return dialWS(address)
	case HTTP:
		return dialHTTP(address)
	case HTTP1:
		return dialHTTP1(address)
	default:
		return nil, errors.New("this network is not suported")
	}
}
