package rpc
import (
	"sync"
)

type ConcurrentChan chan *ConcurrentRequest

type ConcurrentRequest struct {
	data []byte
	cbChan chan []byte
}

type Concurrent struct {
	mut sync.Mutex
	concurrentChan chan *ConcurrentRequest
	actionConcurrentChan chan *ConcurrentRequest
	readChan chan []byte
	writeChan chan []byte
	maxConcurrentRequest	int
}
func NewConcurrent(maxConcurrentRequest int,readChan  chan []byte,writeChan  chan []byte) *Concurrent {
	c:= &Concurrent{
		concurrentChan:make(chan *ConcurrentRequest,maxConcurrentRequest),
		actionConcurrentChan:make(chan *ConcurrentRequest,maxConcurrentRequest),
		readChan :readChan,
		writeChan :writeChan,
		maxConcurrentRequest:maxConcurrentRequest,
	}
	go c.run()
	return c
}
func NewConcurrentRequest(data []byte,cbChan  chan []byte) *ConcurrentRequest {
	c:= &ConcurrentRequest{
		data:data,
		cbChan:cbChan,
	}
	return c
}
func (c *Concurrent)GetMaxConcurrentRequest()int {
	return c.maxConcurrentRequest
}
func (c *Concurrent)run() {
	for{
		select {
		case cr:=<-c.concurrentChan:
			c.writeChan<-cr.data
			if cr.cbChan!=nil{
				c.actionConcurrentChan<-cr
			}
		case b:=<-c.readChan:
			cr:=<-c.actionConcurrentChan
			cr.cbChan<-b
		}
	}
}
