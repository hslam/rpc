package rpc

import (
	"math/rand"
	"time"
	"sync"
	"fmt"
)

const (
	FrameHeaderLength		= 5
)

func PacketFrame(priority uint8,id uint32,body []byte) []byte {
	var buffer=make([]byte,FrameHeaderLength+len(body))
	copy(buffer[0:], []byte{priority})
	copy(buffer[1:], uint32ToBytes(id))
	copy(buffer[FrameHeaderLength:],body)
	return buffer
}

func UnpackFrame(buffer []byte) (priority uint8,id uint32,body []byte,err error) {
	if len(buffer)<FrameHeaderLength{
		err=fmt.Errorf("buffer length %d",len(buffer))
		return
	}
	priority=buffer[:1][0]
	id = bytesToUint32(buffer[1:FrameHeaderLength])
	body=buffer[FrameHeaderLength:]
	return
}

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
		idChan:make(chan uint32,maxRequests*2),
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
	writeStop:=false
	go func() {
		for mr:=range c.requestChan{
			c.mu.Lock()
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
				c.idChan<-id
				mr.id=id
				frameBytes:=PacketFrame(mr.priority,id,mr.data)
				c.writeChan<-frameBytes
				if mr.noResponse==false{
					if c.Length()<=c.maxRequests*2{
						c.Set(id,mr)
					}
				}else {
					if len(c.noResponseChan)>=c.maxRequests{
						<-c.noResponseChan
					}
					c.noResponseChan<-mr
					mr.cbChan<-[]byte("0")
				}
			}()
			c.mu.Unlock()
			if writeStop{
				goto endfor
			}
		}
	endfor:
	}()
	go func() {
		for{
			select {
			case b:=<-c.readChan:
				c.mu.Lock()
				_,ID,body,err:=UnpackFrame(b)
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
				c.mu.Unlock()

			}
		}
	}()
	ticker:=time.NewTicker(time.Second)
	for{
		select {
		case <-c.closeChan:
			writeStop=true
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
				if mr.startTime.Add(time.Millisecond*time.Duration(m.client.timeout)).Before(time.Now()){
					m.Delete(mr.id)
					if len(m.idChan)>0{
						<-m.idChan
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
			frameBytes:=PacketFrame(mr.priority,mr.id,mr.data)
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
				frameBytes:=PacketFrame(mr.priority,mr.id,mr.data)
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
	c.closeChan<-true
}