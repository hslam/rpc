package rpc

import "time"

//IO defines the interface of io
type IO interface {
	NewRequest(priority uint8, data []byte, noResponse bool, cbChan chan []byte) *ioRequest
	RequestChan() requestChan
	ResetMaxRequests(max int)
	Reset(readChan chan []byte, writeChan chan []byte)
	Retry()
	Close()
}

type ioRequest struct {
	id         uint32
	priority   uint8
	data       []byte
	noResponse bool
	cbChan     chan []byte
	startTime  time.Time
}

type requestChan chan *ioRequest
