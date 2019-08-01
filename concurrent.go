package rpc
import (
	"sync"
)

type ConcurrentChan chan *ConcurrentRequest

type ConcurrentRequest struct {
	id  uint64
	data []byte
	cbChan chan []byte
}

type Concurrent struct {
	mut sync.Mutex
	concurrentChan chan *ConcurrentRequest
	actionConcurrentChan chan *ConcurrentRequest
	noResponseConcurrentChan chan *ConcurrentRequest

	readChan chan []byte
	writeChan chan []byte
	conn Conn
	maxConcurrentRequest	int
	returnid		int64
	stop			bool
}
func NewConcurrent(maxConcurrentRequest int,readChan  chan []byte,writeChan  chan []byte,conn Conn) *Concurrent {
	c:= &Concurrent{
		concurrentChan:make(chan *ConcurrentRequest,maxConcurrentRequest),
		actionConcurrentChan:make(chan *ConcurrentRequest,maxConcurrentRequest*2),
		noResponseConcurrentChan:make(chan *ConcurrentRequest,maxConcurrentRequest*2),
		readChan :readChan,
		writeChan :writeChan,
		conn:conn,
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
			for {
				if !c.stop{
					c.writeChan<-cr.data
					if cr.cbChan!=nil{
						if len(c.actionConcurrentChan)<=c.maxConcurrentRequest{
							c.actionConcurrentChan<-cr
						}
					}else {
						if len(c.noResponseConcurrentChan)>=c.maxConcurrentRequest{
							<-c.noResponseConcurrentChan
						}
						c.noResponseConcurrentChan<-cr
					}
					break
				}
			}
		case b:=<-c.readChan:
			if !c.stop{
				cr:=<-c.actionConcurrentChan
				cr.cbChan<-b
			}
		}
	}
}
func (c *Concurrent)retry() {
	c.stop=true
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
	if len(c.noResponseConcurrentChan)>0{
		for i:=0;i<len(c.noResponseConcurrentChan);i++{
			cr:=<-c.noResponseConcurrentChan
			c.writeChan<-cr.data
			c.noResponseConcurrentChan<-cr
		}
	}
	if len(c.actionConcurrentChan)>0{
		for i:=0;i<len(c.actionConcurrentChan);i++{
			cr:=<-c.actionConcurrentChan
			c.writeChan<-cr.data
			if cr.cbChan!=nil{
				c.actionConcurrentChan<-cr
			}
		}
	}
	c.stop=false
}
