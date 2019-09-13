package rpc

import (
	"hslam.com/mgit/Mort/rpc/log"
	"sync"
)

type Pipeline struct {
	retryMu 			sync.RWMutex
	client 				*client
	requestChan 		RequestChan
	actionChan			RequestChan
	noResponseChan		RequestChan
	readChan 			chan []byte
	writeChan 			chan []byte
	maxRequests			int
	closeChan 			chan bool
}
func NewPipeline(maxRequests int,readChan  chan []byte,writeChan  chan []byte) *Pipeline {
	p:= &Pipeline{
		requestChan:make(RequestChan,maxRequests),
		actionChan:make(RequestChan,maxRequests*2),
		noResponseChan:make(RequestChan,maxRequests*2),
		readChan :readChan,
		writeChan :writeChan,
		maxRequests:maxRequests,
		closeChan:make(chan bool,1),
	}
	go p.run()
	return p
}
func (c *Pipeline)NewRequest(priority uint8,data []byte,noResponse bool,cbChan chan []byte) *IORequest {
	r:= &IORequest{
		priority:priority,
		data:data,
		noResponse:noResponse,
		cbChan:cbChan,
	}
	return r
}
func (c *Pipeline)RequestChan()RequestChan{
	return c.requestChan
}
func (c *Pipeline)ResetMaxRequests(max int) {
	c.maxRequests=max
}
func (c *Pipeline)Reset(readChan chan []byte,writeChan chan []byte) {
	c.readChan=readChan
	c.writeChan=writeChan
}
func (c *Pipeline)run() {
	writeStop:=false
	go func() {
		for cr:=range c.requestChan{
			c.retryMu.RLock()
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Errorln("v.reply err", err)
					}
				}()
				c.writeChan<-cr.data
				if cr.noResponse==false{
					if len(c.actionChan)<=c.maxRequests{
						c.actionChan<-cr
					}
				}else {
					if len(c.noResponseChan)>=c.maxRequests{
						<-c.noResponseChan
					}
					c.noResponseChan<-cr
					cr.cbChan<-[]byte("0")
				}
			}()
			c.retryMu.RUnlock()
			if writeStop{
				goto endfor
			}
		}
		endfor:
	}()
	for{
		select {
		case <-c.closeChan:
			writeStop=true
			goto endfor
		case b:=<-c.readChan:
			c.retryMu.RLock()
			cr:=<-c.actionChan
			cr.cbChan<-b
			c.retryMu.RUnlock()
		}
	}
	endfor:
}
func (c *Pipeline)Retry() {
	c.retryMu.Lock()
	if len(c.readChan)>0{
		for i:=0;i<len(c.readChan);i++{
			<-c.readChan
		}
	}
	if len(c.writeChan)>0{
		for i:=0;i<len(c.writeChan);i++{
			<-c.writeChan
		}
	}
	if len(c.noResponseChan)>0{
		for i:=0;i<len(c.noResponseChan);i++{
			cr:=<-c.noResponseChan
			c.writeChan<-cr.data
			c.noResponseChan<-cr
		}
	}
	if len(c.actionChan)>0{
		for i:=0;i<len(c.actionChan);i++{
			cr:=<-c.actionChan
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Errorln("Pipeline.retry", err)
					}
				}()
				c.writeChan<-cr.data
				if cr.noResponse==false{
					c.actionChan<-cr
				}
			}()

		}
	}
	c.retryMu.Unlock()
}
func (c *Pipeline)Close() {
	c.closeChan<-true
}
