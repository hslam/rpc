package rpc

import (
	"math/rand"
	"time"
	"sync"
	"hslam.com/git/x/protocol"
)

type Multiplex struct {
	mu 					sync.RWMutex
	client *client
	requestChan 		chan *IORequest
	cache 				map[uint32]*IORequest
	noResponseChan 		chan *IORequest
	idChan				chan uint32
	readChan			chan []byte
	writeChan			chan []byte
	maxRequests			int
	stop				bool
	closeChan 			chan bool
}

func NewMultiplex(c *client,maxRequests int,readChan  chan []byte,writeChan  chan []byte) *Multiplex {
	m:= &Multiplex{
		client:c,
		requestChan:make(chan *IORequest,maxRequests),
		cache:make(map[uint32]*IORequest,maxRequests*2),
		noResponseChan:make(chan *IORequest,maxRequests*2),
		idChan:make(chan uint32,maxRequests),
		readChan :readChan,
		writeChan :writeChan,
		maxRequests:maxRequests,
		closeChan:make(chan bool,1),
	}
	go m.run()
	return m
}
func (c *Multiplex)NewRequest(priority uint8,data []byte,noResponse  bool,cbChan  chan []byte) *IORequest {
	r:= &IORequest{
		priority:priority,
		data:data,
		noResponse:noResponse,
		cbChan:cbChan,
		startTime:time.Now(),
	}
	return r
}
func (c *Multiplex)RequestChan()RequestChan {
	return c.requestChan
}
func (c *Multiplex)ResetMaxRequests(max int) {
	c.maxRequests=max
}
func (c *Multiplex)Reset(readChan chan []byte,writeChan chan []byte) {
	c.readChan=readChan
	c.writeChan=writeChan
}
func (c *Multiplex)run() {
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	var startbit =uint(r.Intn(13))
	id:=uint32(r.Int31n(int32(1<<startbit)))
	max_id:=uint32(1<<32-1)
	go func() {
		for mr:=range c.requestChan{
			func(){
				c.mu.Lock()
				defer c.mu.Unlock()
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					for{
						id=(id+1)%max_id
						if !c.IsExisted(id){
							break
						}
					}
					mr.id=id
					frameBytes:=protocol.PacketFrame(mr.priority,id,mr.data)
					c.writeChan<-frameBytes
					if mr.noResponse==false{
						c.idChan<-id
						c.Set(id,mr)
					}else {
						if len(c.noResponseChan)>=c.maxRequests{
							<-c.noResponseChan
						}
						c.noResponseChan<-mr
						mr.cbChan<-[]byte("0")
					}
				}()
			}()
		}
	}()
	go func() {
		for b:=range c.readChan{
			func(){
				c.mu.Lock()
				defer c.mu.Unlock()
				_,ID,body,err:=protocol.UnpackFrame(b)
				if err!=nil{
					return
				}
				mr:=c.Get(ID)
				if mr!=nil{
					mr.cbChan<-body
				}
				c.Delete(ID)
				if len(c.idChan)>0{
					<-c.idChan
				}
			}()
		}
	}()
	ticker:=time.NewTicker(time.Second)
	for{
		select {
		case <-c.closeChan:
			close(c.closeChan)
			ticker.Stop()
			ticker=nil
			goto endfor
		case <-ticker.C:
			c.deleteOld()
		}
	}
endfor:
}
func (m *Multiplex)deleteOld() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.cache)>0{
		for _,mr:=range m.cache{
			func() {
				defer func() {
					if err := recover(); err != nil {
					}
				}()
				if m.client.timeout>0{
					if mr.startTime.Add(time.Millisecond*time.Duration(m.client.timeout)).Before(time.Now()){
						m.Delete(mr.id)
						if len(m.idChan)>0{
							<-m.idChan
						}
					}
				}else {
					if mr.startTime.Add(time.Minute*5).Before(time.Now()){
						m.Delete(mr.id)
						if len(m.idChan)>0{
							<-m.idChan
						}
					}
				}
			}()
		}
	}
}
func (c *Multiplex)Retry() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.idChan=make(chan uint32,c.maxRequests*2)
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
			mr:=<-c.noResponseChan
			c.idChan<-mr.id
			frameBytes:=protocol.PacketFrame(mr.priority,mr.id,mr.data)
			c.writeChan<-frameBytes
			c.noResponseChan<-mr
		}
	}
	if len(c.cache)>0{
		for _,mr:=range c.cache{
			func() {
				defer func() {
					if err := recover(); err != nil {
					}
				}()
				c.idChan<-mr.id
				frameBytes:=protocol.PacketFrame(mr.priority,mr.id,mr.data)
				c.writeChan<-frameBytes
			}()
		}
	}
}

func(c *Multiplex)IsExisted(id uint32) bool{
	if _,ok:=c.cache[id];!ok{
		return false
	}
	return true
}
func(c *Multiplex)Set(id uint32,multiplexRequest *IORequest) {
	c.cache[id]=multiplexRequest
}
func(c *Multiplex)Get(id uint32)*IORequest {
	if _,ok:=c.cache[id];!ok{
		return nil
	}
	return c.cache[id]
}
func(c *Multiplex)Delete(id uint32) {
	if _,ok:=c.cache[id];!ok{
		return
	}
	delete(c.cache,id)
}
func(c *Multiplex)Length() int{
	return len(c.cache)
}
func (c *Multiplex)Close() {
	close(c.requestChan)
	close(c.noResponseChan)
	close(c.idChan)
	c.readChan=nil
	c.writeChan=nil
	c.client=nil
	c.cache=nil
	c.closeChan<-true
}
