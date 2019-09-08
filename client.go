package rpc
import (
	"sync"
	"hslam.com/mgit/Mort/idgenerator"
	"time"
	"math/rand"
	"hslam.com/mgit/Mort/rpc/log"
)

func Dial(network,address,codec string) (*Client, error) {
	transporter,err:=dial(network,address)
	if err!=nil{
		return nil,err
	}
	return NewClient(transporter,codec)
}
func DialWithPipeline(network,address,codec string,MaxPipelineRequest int) (*Client, error) {
	transporter,err:=dial(network,address)
	if err!=nil{
		return nil,err
	}
	return NewClientnWithConcurrent(transporter,codec,MaxPipelineRequest+1)
}
type Client struct {
	mu 					sync.Mutex
	conn				Conn
	closed				bool
	hystrix				bool
	batchEnabled		bool
	batchAsync			bool
	batch				*Batch
	pipeline 			*Pipeline
	pipelineChan		chan bool
	readChan			chan []byte
	writeChan			chan []byte
	finishChan			chan bool
	stopChan			chan bool
	funcsCodecType 		CodecType
	compressType 		CompressType
	compressLevel 		CompressLevel
	idgenerator			*idgenerator.IDGen
	client_id			int64
	timeout 			int64
	heartbeatTimeout	int64
	errCntChan 			chan int
	errCnt 				int
	maxErrPerSecond		int
	maxErrHeartbeat		int
	errCntHeartbeat		int
}
func NewClient(conn	Conn,codec string)  (*Client, error) {
	return NewClientnWithConcurrent(conn,codec,DefaultMaxPipelineRequest+1)
}
func NewClientnWithConcurrent(conn	Conn,codec string,maxPipeliningRequest int)  (*Client, error)  {
	funcsCodecType,err:=FuncsCodecType(codec)
	if err!=nil{
		return nil,err
	}
	finishChan := make(chan bool)
	stopChan := make(chan bool)
	errCntChan:=make(chan int,1000000)
	var client_id int64=0
	idgenerator:=idgenerator.NewSnowFlake(client_id)
	client :=  &Client{
		client_id:client_id,
		conn:conn,
		finishChan:finishChan,
		stopChan:stopChan,
		funcsCodecType:funcsCodecType,
		compressLevel:NoCompression,
		compressType:CompressTypeNocom,
	}
	client.readChan=make(chan []byte,maxPipeliningRequest)
	client.writeChan= make(chan []byte,maxPipeliningRequest)
	client.idgenerator=idgenerator
	client.errCntChan=errCntChan
	client.timeout=DefaultClientTimeout
	client.maxErrPerSecond=DefaultClientMaxErrPerSecond
	client.heartbeatTimeout=DefaultClientHearbeatTimeout
	client.maxErrHeartbeat=DefaultClientMaxErrHearbeat
	client.conn.Handle(client.readChan,client.writeChan,client.stopChan,client.finishChan)
	go func() {
		defer func() {
			defer client.Close()
			close(client.writeChan)
			close(client.readChan)
			close(stopChan)
			close(errCntChan)
		}()
		for{
			select {
			case <-client.stopChan:
				client.Close()
			}
		}

	}()
	go func() {
		for c :=range client.errCntChan{
			client.errCnt+=c
		}
	}()
	go func() {
		time.Sleep(time.Millisecond*time.Duration(rand.Int63n(800)))
		ticker:=time.NewTicker(time.Second)
		heartbeatTicker:=time.NewTicker(time.Millisecond*DefaultClientHearbeatTicker)
		retryTicker:=time.NewTicker(time.Millisecond*DefaultClientRetryTicker)
		for{
			select {
			case <-heartbeatTicker.C:
				if client.closed==false{
					err:=client.heartbeat()
					if err!=nil{
						client.errCntHeartbeat+=1
					}else {
						client.errCntHeartbeat-=1
						if client.errCntHeartbeat<0{
							client.errCntHeartbeat=0
						}
					}
					if client.errCntHeartbeat>=client.maxErrHeartbeat{
						client.Close()
						client.errCntHeartbeat=0
					}
				}
			case <-retryTicker.C:
				if client.closed==true{
					client.Close()
					err:=client.conn.Retry()
					if err!=nil{
						client.hystrix=true
						log.Errorln(client.client_id,"retry connection err ",err)
					}else {
						client.closed=false
						client.hystrix=false
						client.pipeline.retry()
					}
				}
			case <-ticker.C:
				if client.errCnt>=DefaultClientMaxErrPerSecond{
					client.hystrix=true
					client.Close()
				}else {
					client.hystrix=false
				}
				client.errCnt-=client.errCnt
			}
		}
	}()
	client.pipeline=NewPipeline(maxPipeliningRequest,client.readChan,client.writeChan)
	client.pipelineChan = make(chan bool,maxPipeliningRequest)
	return client, nil
}
func (c *Client)EnableBatch(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batchEnabled = true
	c.batch=NewBatch(c,DefaultMaxDelayNanoSecond*c.conn.TickerFactor())
	c.batch.SetMaxBatchRequest(DefaultMaxBatchRequest*c.conn.BatchFactor())
}
func (c *Client)EnableBatchAsync(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batchAsync = true
}
func (c *Client)SetMaxBatchRequest(maxBatchRequest int)error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxBatchRequest<=0{
		return ErrSetMaxBatchRequest
	}
	c.batch.SetMaxBatchRequest(maxBatchRequest)
	return nil
}
func (c *Client)GetMaxBatchRequest()int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.batch.GetMaxBatchRequest()
}
func (c *Client)GetMaxPipelineRequest()(int){
	return c.pipeline.GetMaxPipelineRequest()-1
}
func (c *Client)SetCompressType(compress string){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.compressType=getCompressType(compress)
	c.compressLevel=DefaultCompression
}
func (c *Client)SetCompressLevel(compress,level string){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.compressType=getCompressType(compress)
	c.compressLevel=getCompressLevel(level)
}
func (c *Client)SetID(id int64)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if id<0{
		return ErrSetClientID
	}else if id>1023{
		return ErrSetClientID
	}
	c.client_id=id
	c.idgenerator=idgenerator.NewSnowFlake(id)
	return nil
}
func (c *Client)GetID()int64{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.client_id
}
func (c *Client)SetTimeout(timeout int64)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout<=0{
		return ErrSetTimeout
	}
	c.timeout=timeout
	return nil
}
func (c *Client)GetTimeout()int64{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeout
}

func (c *Client)SetHeartbeatTimeout(timeout int64)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout<=0{
		return ErrSetTimeout
	}
	c.heartbeatTimeout=timeout
	return nil
}
func (c *Client)GetHeartbeatTimeout()int64{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.heartbeatTimeout
}
func (c *Client)SetMaxErrPerSecond(maxErrPerSecond int)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxErrPerSecond<=0{
		return ErrSetMaxErrPerSecond
	}
	c.maxErrPerSecond=maxErrPerSecond
	return nil
}
func (c *Client)GetMaxErrPerSecond()int{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.maxErrPerSecond
}
func (c *Client)SetMaxErrHeartbeat(maxErrHeartbeat int)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxErrHeartbeat<=0{
		return ErrSetMaxErrHeartbeat
	}
	c.maxErrHeartbeat=maxErrHeartbeat
	return nil
}
func (c *Client)GetMaxErrHeartbeat()int{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.maxErrHeartbeat
}
func (c *Client)CodecName()string {
	return FuncsCodecName(c.funcsCodecType)
}
func (c *Client)CodecType()CodecType {
	return c.funcsCodecType
}
func (c *Client)Call(name string, args interface{}, reply interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Call failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
	}
	if c.batchEnabled{
		err=c.batchCall(name,args,reply)
		if err!=nil{
			c.errCntChan<-1
		}
		return
	}
	err= c.call(name,args,reply)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}
func (c *Client)CallNoRequest(name string, reply interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Call failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
	}
	if c.batchEnabled{
		err=c.batchCall(name,nil,reply)
		if err!=nil{
			c.errCntChan<-1
		}
		return
	}
	err= c.call(name,nil,reply)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}
func (c *Client)CallNoResponse(name string, args interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("CallNoResponse failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
	}
	if c.batchEnabled{
		err= c.batchCall(name,args,nil)
		if err!=nil{
			c.errCntChan<-1
		}
		return
	}
	err= c.call(name,args,nil)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}

func (c *Client)OnlyCall(name string) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("OnlyCall failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
	}
	if c.batchEnabled{
		err= c.batchCall(name,nil,nil)
		if err!=nil{
			c.errCntChan<-1
		}
		return
	}
	err= c.call(name,nil,nil)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}

func (c *Client)Ping() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Client.Ping failed:", err)
		}
	}()
	err:=c.heartbeat()
	if err!=nil{
		return false
	}
	return true
}
func (c *Client)RemoteCall(b []byte)([]byte,error){
	c.pipelineChan<-true
	cbChan := make(chan []byte,1)
	c.pipeline.pipelineRequestChan<-NewPipelineRequest(b,false,cbChan)
	data,ok := <-cbChan
	<-c.pipelineChan
	if ok{
		return data,nil
	}
	return nil,ErrRemoteCall
}
func (c *Client)RemoteCallNoResponse(b []byte)(error){
	c.pipelineChan<-true
	cbChan := make(chan []byte,1)
	c.pipeline.pipelineRequestChan<-NewPipelineRequest(b,true,cbChan)
	<-cbChan
	<-c.pipelineChan
	return nil
}
func (c *Client)Close() ( err error) {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.conn.Close()
}
func (c *Client)Closed()bool{
	return c.closed
}

func (c *Client)heartbeat() ( err error) {
	uid:=c.idgenerator.GenUniqueIDInt64()
	msg:=&Msg{}
	msg.id=uid
	msg.msgType=MsgType(MsgTypeHea)
	msg_bytes, _ :=msg.Encode()
	ch := make(chan int)
	go func() {
		data,err:=c.RemoteCall(msg_bytes)
		if err != nil {
			log.Errorln("Write error: ", err)
			ch<-1
			return
		}
		err=msg.Decode(data)
		if err==nil{
			if msg.id!=uid{
				err=ErrClientId
			}
		}
		ch<-1
		return
	}()
	select {
	case <-ch:
	case <-time.After(time.Millisecond * time.Duration(c.heartbeatTimeout)):
		err=ErrTimeOut
	}
	return err
}

func (c *Client)call(name string, args interface{}, reply interface{}) ( err error) {
	clientCodec:=&ClientCodec{}
	clientCodec.client_id=c.client_id
	clientCodec.req_id=uint64(c.idgenerator.GenUniqueIDInt64())
	clientCodec.name=name
	if args!=nil{
		clientCodec.args=args
		clientCodec.noRequest=false
	}else {
		clientCodec.args=nil
		clientCodec.noRequest=true
	}
	clientCodec.funcsCodecType=c.funcsCodecType
	if reply!=nil{
		clientCodec.noResponse=false
	}else {
		clientCodec.noResponse=true
	}
	rpc_req_bytes, _ :=clientCodec.Encode()
	if clientCodec.noResponse{
		err=c.RemoteCallNoResponse(rpc_req_bytes)
		if err != nil {
			log.Errorln("Write error: ", err)
			return err
		}
		return nil
	}else {
		ch := make(chan int)
		go func() {
			var data []byte
			data,err=c.RemoteCall(rpc_req_bytes)
			if err != nil {
				log.Errorln("Write error: ", err)
				ch<-1
				return
			}
			clientCodec.reply=reply
			err=clientCodec.Decode(data)
			ch<-1
		}()
		select {
		case <-ch:
		case <-time.After(time.Millisecond * time.Duration(c.timeout)):
			err=ErrTimeOut
		}
		return err
	}

}
func (c *Client)batchCall(name string, args interface{}, reply interface{}) ( err error) {
	cr:=&BatchRequest{
		id:uint64(c.idgenerator.GenUniqueIDInt64()),
		name:name,
		noResponse:false,
	}
	if args!=nil{
		args_bytes,err:=ArgsEncode(args,c.funcsCodecType)
		if err!=nil{
			log.Errorln("ArgsEncode error: ", err)
		}
		cr.args_bytes=args_bytes
		cr.noRequest=false
	}else {
		cr.args_bytes=nil
		cr.noRequest=true
	}
	if reply!=nil{
		cr.reply_bytes= make(chan []byte, 1)
		cr.reply_error=make(chan error, 1)
		cr.noResponse=false
	}else {
		cr.noResponse=true
	}
	c.batch.reqChan<-cr
	if cr.noResponse{
		return nil
	}else{
		select {
		case bytes ,ok:= <- cr.reply_bytes:
			if ok{
				return ReplyDecode(bytes,reply,c.funcsCodecType)
			}
		case err ,ok:= <- cr.reply_error:
			if ok{
				return err
			}
		case <-time.After(time.Millisecond * time.Duration(c.timeout)):
			close(cr.reply_bytes)
			close(cr.reply_error)
			return ErrTimeOut
		}
		return nil
	}

}
