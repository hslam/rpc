package rpc

import "time"

type IO interface {
	NewRequest(priority uint8, data []byte, noResponse bool, cbChan chan []byte) *IORequest
	RequestChan() RequestChan
	ResetMaxRequests(max int)
	Reset(readChan chan []byte, writeChan chan []byte)
	Retry()
	Close()
}

type IORequest struct {
	id         uint32
	priority   uint8
	data       []byte
	noResponse bool
	cbChan     chan []byte
	startTime  time.Time
}

type RequestChan chan *IORequest
