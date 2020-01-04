package rpc

import (
	"errors"
)
//Dialer
type Conn interface {
	Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool)
	NoDelay(enable bool)
	Multiplexing(enable bool)
	TickerFactor()(int)
	BatchFactor()(int)
	Retry()(error)
	Close()(error)
	Closed()(bool)
}


func dial(network,address string) (Conn, error) {
	switch network {
	case IPC:
		return DialIPC(address)
	case TCP:
		return DialTCP(address)
	case UDP:
		return DialUDP(address)
	case QUIC:
		return DialQUIC(address)
	case WS:
		return DialWS(address)
	case HTTP:
		return DialHTTP(address)
	case HTTP1:
		return DialHTTP1(address)
	default:
		return nil, errors.New("this network is not suported")
	}
}