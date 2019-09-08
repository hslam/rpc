package rpc
import (
	"sync"
	"time"
	"hslam.com/mgit/Mort/rpc/pb"
	"hslam.com/mgit/Mort/rpc/log"
)

type RequestChan chan *BatchRequest

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
	reqChan RequestChan
	client *Client
	readyRequests []*BatchRequest
	maxBatchRequest	int
	maxDelayNanoSecond	int
}
func NewBatch(client *Client,maxDelayNanoSecond int) *Batch {
	c:= &Batch{
		reqChan:make(chan *BatchRequest,DefaultMaxCacheRequest),
		client:client,
		readyRequests:make([]*BatchRequest,0),
		maxBatchRequest:DefaultMaxBatchRequest,
		maxDelayNanoSecond:maxDelayNanoSecond,
	}
	go c.run()
	return c
}

func (c *Batch)GetMaxBatchRequest()int {
	return c.maxBatchRequest
}
func (c *Batch)SetMaxBatchRequest(maxBatchRequest int) {
	c.maxBatchRequest=maxBatchRequest
}
func (c *Batch)run() {
	go func() {
		for cr := range c.reqChan {
			c.mut.Lock()
			c.readyRequests=append(c.readyRequests, cr)
			if len(c.readyRequests)>=c.maxBatchRequest{
				crs:=c.readyRequests[:]
				c.readyRequests=nil
				c.readyRequests=make([]*BatchRequest,0)
				c.Ticker(crs)
			}
			c.mut.Unlock()
		}
	}()
	tick := time.NewTicker(1 * time.Nanosecond*time.Duration(c.maxDelayNanoSecond))
	for {
		select {
		case <-tick.C:
			c.mut.Lock()
			if len(c.readyRequests)>c.maxBatchRequest{
				crs:=c.readyRequests[:c.maxBatchRequest]
				c.readyRequests=c.readyRequests[c.maxBatchRequest:]
				c.Ticker(crs)
			}else  if len(c.readyRequests)>0{
				crs:=c.readyRequests[:]
				c.readyRequests=c.readyRequests[len(c.readyRequests):]
				c.Ticker(crs)
			}else {
			}
			c.mut.Unlock()
		}
	}
}
func (c *Batch)Ticker(crs []*BatchRequest){
	req_bytes_s:=make([][]byte,len(crs))
	NoResponseCnt:=0
	var noResponse bool
	for i,v :=range crs{
		req:=&Request{v.id,v.name,v.noRequest,v.noResponse,v.args_bytes}
		req_bytes,_:=req.Encode()
		req_bytes_s[i]=req_bytes
		if v.noResponse==true{
			NoResponseCnt++
		}
	}
	if NoResponseCnt==len(crs){
		noResponse=true
	}else {
		noResponse=false
	}
	batch:=&BatchCodec{async:c.client.batchAsync,data:req_bytes_s}
	batch_bytes,err:=batch.Encode()
	msg:=&Msg{}
	msg.id=c.client.GetID()
	msg.data=batch_bytes
	msg.batch=true
	msg.msgType=MsgType(pb.MsgType_req)
	msg.codecType=c.client.CodecType()
	msg_bytes,err:=msg.Encode()
	if err==nil{
		if noResponse==false{
			data,err:=c.client.RemoteCall(msg_bytes)
			if err == nil {
				msg:=&Msg{}
				err=msg.Decode(data)
				if err!=nil{
					return
				}
				if msg.msgType!=MsgType(pb.MsgType_res){
					return
				}
				batch:=&BatchCodec{}
				err=batch.Decode(msg.data)
				if err!=nil{
					return
				}
				if len(batch.data)==len(crs){
					for i,v:=range crs{
						res:=Response{}
						res.Decode(batch.data[i])
						if crs[i].id==res.id&& v.noResponse==false{
							func() {
								defer func() {
									if err := recover(); err != nil {
										log.Errorln("v.reply err", err)
									}
								}()
								if res.err!=nil{
									v.reply_error<-res.err
								}else {
									v.reply_bytes<-res.data
								}
							}()

						}
					}
				}
			}
		}else {
			_=c.client.RemoteCallNoResponse(msg_bytes)
		}

	}
}