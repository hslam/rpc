package rpc
import (
	"sync"
)

type PipelineRequestChan chan *PipelineRequest

type PipelineRequest struct {
	data []byte
	noResponse  bool
	cbChan chan []byte
}

type Pipeline struct {
	mut sync.Mutex
	pipelineRequestChan PipelineRequestChan
	actionPipelineRequestChan PipelineRequestChan
	noResponsePipelineRequestChan PipelineRequestChan
	readChan chan []byte
	writeChan chan []byte
	maxPipelineRequest	int
	returnid		int64
	stop			bool
}
func NewPipeline(maxPipelineRequest int,readChan  chan []byte,writeChan  chan []byte) *Pipeline {
	c:= &Pipeline{
		pipelineRequestChan:make(PipelineRequestChan,maxPipelineRequest),
		actionPipelineRequestChan:make(PipelineRequestChan,maxPipelineRequest*2),
		noResponsePipelineRequestChan:make(PipelineRequestChan,maxPipelineRequest*2),
		readChan :readChan,
		writeChan :writeChan,
		maxPipelineRequest:maxPipelineRequest,
	}
	go c.run()
	return c
}
func NewPipelineRequest(data []byte,noResponse  bool,cbChan  chan []byte) *PipelineRequest {
	c:= &PipelineRequest{
		data:data,
		noResponse:noResponse,
		cbChan:cbChan,
	}
	return c
}
func (c *Pipeline)GetMaxPipelineRequest()int {
	return c.maxPipelineRequest
}
func (c *Pipeline)run() {
	for{
		select {
		case cr:=<-c.pipelineRequestChan:
			for {
				if !c.stop{
					c.writeChan<-cr.data
					if cr.noResponse==false{
						if len(c.actionPipelineRequestChan)<=c.maxPipelineRequest{
							c.actionPipelineRequestChan<-cr
						}
					}else {
						if len(c.noResponsePipelineRequestChan)>=c.maxPipelineRequest{
							<-c.noResponsePipelineRequestChan
						}
						c.noResponsePipelineRequestChan<-cr
						cr.cbChan<-[]byte("0")
					}
					break
				}
			}
		case b:=<-c.readChan:
			if !c.stop{
				cr:=<-c.actionPipelineRequestChan
				cr.cbChan<-b
			}
		}
	}
}
func (c *Pipeline)retry() {
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
	if len(c.noResponsePipelineRequestChan)>0{
		for i:=0;i<len(c.noResponsePipelineRequestChan);i++{
			cr:=<-c.noResponsePipelineRequestChan
			c.writeChan<-cr.data
			c.noResponsePipelineRequestChan<-cr
		}
	}
	if len(c.actionPipelineRequestChan)>0{
		for i:=0;i<len(c.actionPipelineRequestChan);i++{
			cr:=<-c.actionPipelineRequestChan
			c.writeChan<-cr.data
			if cr.noResponse==false{
				c.actionPipelineRequestChan<-cr
			}
		}
	}
	c.stop=false
}
