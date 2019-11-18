package rpc

import (
	"errors"
)
//Dialer
type Conn interface {
	Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool)
	Buffer(enable bool)
	Multiplexing(enable bool)
	TickerFactor()(int)
	BatchFactor()(int)
	Retry()(error)
	Close()(error)
	Closed()(bool)
}


func dial(network,address string) (Conn, error) {
	switch network {
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
	case HTTP2:
		return DialHTTP2(address)
	default:
		return nil, errors.New("this network is not suported")
	}
}