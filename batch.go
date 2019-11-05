package rpc
import (
	"sync"
	"time"
	"hslam.com/git/x/rpc/pb"
	"hslam.com/git/x/rpc/log"
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
func NewBatch(c *client,maxDelayNanoSecond int) *Batch {
	b:= &Batch{
		reqChan:make(chan *BatchRequest,DefaultMaxCacheRequest),
		client:c,
		readyRequests:make([]*BatchRequest,0),
		maxBatchRequest:DefaultMaxBatchRequest,
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
func (b *Batch)Ticker(crs []*BatchRequest){
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
	batch:=&BatchCodec{async:b.client.batchAsync,data:req_bytes_s}
	batch_bytes,err:=batch.Encode()
	msg:=&Msg{}
	msg.id=b.client.GetID()
	msg.data=batch_bytes
	msg.batch=true
	msg.msgType=MsgType(pb.MsgType_req)
	msg.codecType=b.client.CodecType()
	msg_bytes,err:=msg.Encode()
	if err==nil{
		if noResponse==false{
			data,err:=b.client.RemoteCall(msg_bytes)
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
				batch=nil
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