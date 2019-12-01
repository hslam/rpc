package rpc

import (
	"sync"
	"time"
)

type BatchRequestChan chan *BatchRequest

type BatchRequest struct {
	id uint64
	name string
	args_bytes []byte
	reply_bytes chan []byte
	reply_error chan error
	noRequest bool
	noResponse bool
}
type Batch struct {
	mut sync.Mutex
	reqChan BatchRequestChan
	client *client
	readyRequests []*BatchRequest
	maxBatchRequest	int
	maxDelayNanoSecond	int
	closeChan 			chan bool
}
func NewBatch(c *client,maxDelayNanoSecond int,maxBatchRequest int) *Batch {
	b:= &Batch{
		reqChan:make(chan *BatchRequest,DefaultMaxCacheRequest),
		client:c,
		readyRequests:make([]*BatchRequest,0),
		maxBatchRequest:maxBatchRequest,
		maxDelayNanoSecond:maxDelayNanoSecond,
		closeChan:make(chan bool,1),
	}
	go b.run()
	return b
}

func (b *Batch)GetMaxBatchRequest()int {
	return b.maxBatchRequest
}
func (b *Batch)SetMaxBatchRequest(maxBatchRequest int) {
	b.maxBatchRequest=maxBatchRequest
}
func (b *Batch)run() {
	go func() {
		for cr := range b.reqChan {
			func(){
				b.mut.Lock()
				defer b.mut.Unlock()
				b.readyRequests=append(b.readyRequests, cr)
				if len(b.readyRequests)>=b.maxBatchRequest{
					crs:=b.readyRequests[:]
					b.readyRequests=nil
					b.readyRequests=make([]*BatchRequest,0)
					b.Ticker(crs)
				}
			}()
		}
	}()
	tick := time.NewTicker(1 * time.Nanosecond*time.Duration(b.maxDelayNanoSecond))
	for {
		select {
		case <-b.closeChan:
			close(b.closeChan)
			tick.Stop()
			tick=nil
			goto endfor
		case <-tick.C:
			func(){
				b.mut.Lock()
				defer b.mut.Unlock()
				if len(b.readyRequests)>b.maxBatchRequest{
					crs:=b.readyRequests[:b.maxBatchRequest]
					b.readyRequests=b.readyRequests[b.maxBatchRequest:]
					b.Ticker(crs)
				}else  if len(b.readyRequests)>0{
					crs:=b.readyRequests[:]
					b.readyRequests=b.readyRequests[len(b.readyRequests):]
					b.Ticker(crs)
				}
			}()
		}
	}
endfor:
}
func (b *Batch)Ticker(brs []*BatchRequest){
	NoResponseCnt:=0
	var noResponse bool
	clientCodec:=&ClientCodec{}
	clientCodec.client_id=b.client.GetID()
	clientCodec.batch=true
	clientCodec.batchAsync=b.client.batchAsync
	clientCodec.requests=brs
	clientCodec.funcsCodecType=b.client.CodecType()
	for _,v :=range brs{
		if v.noResponse==true{
			NoResponseCnt++
		}
	}
	if NoResponseCnt==len(brs){
		noResponse=true
	}else {
		noResponse=false
	}
	msg_bytes,err:=clientCodec.Encode()
	if err==nil{
		if noResponse==false{
			data,err:=b.client.RemoteCall(msg_bytes)
			if err == nil {
				err:=clientCodec.Decode(data)
				if err!=nil{
					return
				}
				if len(clientCodec.responses)==len(brs){
					for i,v:=range brs{
						if brs[i].id==clientCodec.responses[i].id&& v.noResponse==false{
							func() {
								defer func() {
									if err := recover(); err != nil {
										Errorln("v.reply err", err)
									}
								}()
								if clientCodec.responses[i].err!=nil{
									v.reply_error<-clientCodec.responses[i].err
								}else {
									v.reply_bytes<-clientCodec.responses[i].data
								}
							}()
						}
					}
				}
			}
		}else {
			_=b.client.RemoteCallNoResponse(msg_bytes)
		}

	}
}

func (b *Batch)Close() {
	close(b.reqChan)
	b.client=nil
	b.readyRequests=nil
	b.closeChan<-true
}