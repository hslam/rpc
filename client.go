package rpc
import (
	"sync"
	"hslam.com/mgit/Mort/idgenerator"
	"time"
	"math/rand"
	"hslam.com/mgit/Mort/rpc/log"
)
type Client interface {
	EnableBatch()
	EnableBatchAsync()
	SetMaxBatchRequest(maxBatchRequest int)error
	GetMaxBatchRequest()int
	GetMaxPipelineRequest()(int)
	SetCompressType(compress string)
	SetCompressLevel(compress,level string)
	SetID(id int64)error
	GetID()int64
	SetTimeout(timeout int64)error
	GetTimeout()int64
	SetHeartbeatTimeout(timeout int64)error
	GetHeartbeatTimeout()int64
	SetMaxErrPerSecond(maxErrPerSecond int)error
	GetMaxErrPerSecond()int
	SetMaxErrHeartbeat(maxErrHeartbeat int)error
	GetMaxErrHeartbeat()int
	CodecName()string
	CodecType()CodecType
	Call(name string, args interface{}, reply interface{}) ( err error)
	CallNoRequest(name string, reply interface{}) ( err error)
	CallNoResponse(name string, args interface{}) ( err error)
	OnlyCall(name string) ( err error)
	Ping() bool
	DisableRetry()
	Close() ( err error)
	Closed()bool
}
func Dial(network,address,codec string) (Client, error) {
	transporter,err:=dial(network,address)
	if err!=nil{
		return nil,err
	}
	return NewClient(transporter,codec)
}
func DialWithPipelining(network,address,codec string,MaxPipelineRequest int) (Client, error) {
	transporter,err:=dial(network,address)
	if err!=nil{
		return nil,err
	}
	return NewClientnWithConcurrent(transporter,codec,MaxPipelineRequest*2)
}


type client struct {
	mu 					sync.Mutex
	conn				Conn
	closed				bool
	disconnect			bool
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
	retry 				bool
}
func NewClient(conn	Conn,codec string)  (*client, error) {
	return NewClientnWithConcurrent(conn,codec,DefaultMaxPipelineRequest*2)
}
func NewClientnWithConcurrent(conn	Conn,codec string,maxPipeliningRequest int)  (*client, error)  {
	funcsCodecType,err:=FuncsCodecType(codec)
	if err!=nil{
		return nil,err
	}
	var client_id int64=0
	idgenerator:=idgenerator.NewSnowFlake(client_id)
	c :=  &client{
		client_id:client_id,
		conn:conn,
		finishChan:make(chan bool,1),
		stopChan:make(chan bool,1),
		funcsCodecType:funcsCodecType,
		compressLevel:NoCompression,
		compressType:CompressTypeNocom,
		retry:true,
	}
	c.readChan=make(chan []byte,maxPipeliningRequest)
	c.writeChan= make(chan []byte,maxPipeliningRequest)
	c.idgenerator=idgenerator
	c.errCntChan=make(chan int,1000000)
	c.timeout=DefaultClientTimeout
	c.maxErrPerSecond=DefaultClientMaxErrPerSecond
	c.heartbeatTimeout=DefaultClientHearbeatTimeout
	c.maxErrHeartbeat=DefaultClientMaxErrHearbeat
	c.conn.Handle(c.readChan,c.writeChan,c.stopChan,c.finishChan)
	go func() {
		for i :=range c.errCntChan{
			c.errCnt+=i
			if c.closed{
				goto endfor
			}
		}
		endfor:
	}()
	go c.run()
	c.pipeline=NewPipeline(maxPipeliningRequest,c.readChan,c.writeChan)
	c.pipelineChan = make(chan bool,maxPipeliningRequest)
	return c, nil
}
func (c *client)run(){
		time.Sleep(time.Millisecond*time.Duration(rand.Int63n(800)))
		ticker:=time.NewTicker(time.Second)
		heartbeatTicker:=time.NewTicker(time.Millisecond*DefaultClientHearbeatTicker)
		retryTicker:=time.NewTicker(time.Millisecond*DefaultClientRetryTicker)
		for{
			select {
			case <-c.finishChan:
				log.Traceln(c.client_id,"client.run finishChan")
				c.retryConnect()
			case <-heartbeatTicker.C:
				if c.disconnect==false&&c.retry{
					err:=c.heartbeat()
					if err!=nil{
						c.errCntHeartbeat+=1
					}else {
						c.errCntHeartbeat-=1
						if c.errCntHeartbeat<0{
							c.errCntHeartbeat=0
						}
					}
					if c.errCntHeartbeat>=c.maxErrHeartbeat{
						c.Disconnect()
						c.errCntHeartbeat=0
					}
				}
			case <-retryTicker.C:
				if c.disconnect==true&&c.retry{
					c.retryConnect()
				}
			case <-ticker.C:
				if c.closed{
					goto endfor
				}
				if !c.retry{
					return
				}
				if c.errCnt>=DefaultClientMaxErrPerSecond{
					c.hystrix=true
				}else {
					c.hystrix=false
				}
				c.errCnt=0
			}

		}
	endfor:
}
func (c *client)retryConnect(){
	c.Disconnect()
	err:=c.conn.Retry()
	if err!=nil{
		c.hystrix=true
		log.Traceln(c.client_id,"retry connection err ",err)
	}else {
		log.Traceln(c.client_id,"retry connection success")
		c.disconnect=false
		c.hystrix=false
		c.pipeline.retry()
	}
}
func (c *client)EnableBatch(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batchEnabled = true
	c.batch=NewBatch(c,DefaultMaxDelayNanoSecond*c.conn.TickerFactor())
	c.batch.SetMaxBatchRequest(DefaultMaxBatchRequest*c.conn.BatchFactor())
}
func (c *client)EnableBatchAsync(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batchAsync = true
}
func (c *client)SetMaxBatchRequest(maxBatchRequest int)error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxBatchRequest<=0{
		return ErrSetMaxBatchRequest
	}
	c.batch.SetMaxBatchRequest(maxBatchRequest)
	return nil
}
func (c *client)GetMaxBatchRequest()int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.batch.GetMaxBatchRequest()
}
func (c *client)GetMaxPipelineRequest()(int){
	return c.pipeline.GetMaxPipelineRequest()-1
}
func (c *client)SetCompressType(compress string){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.compressType=getCompressType(compress)
	c.compressLevel=DefaultCompression
}
func (c *client)SetCompressLevel(compress,level string){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.compressType=getCompressType(compress)
	c.compressLevel=getCompressLevel(level)
}
func (c *client)SetID(id int64)error{
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
func (c *client)GetID()int64{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.client_id
}
func (c *client)SetTimeout(timeout int64)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout<=0{
		return ErrSetTimeout
	}
	c.timeout=timeout
	return nil
}
func (c *client)GetTimeout()int64{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeout
}

func (c *client)SetHeartbeatTimeout(timeout int64)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout<=0{
		return ErrSetTimeout
	}
	c.heartbeatTimeout=timeout
	return nil
}
func (c *client)GetHeartbeatTimeout()int64{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.heartbeatTimeout
}
func (c *client)SetMaxErrPerSecond(maxErrPerSecond int)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxErrPerSecond<=0{
		return ErrSetMaxErrPerSecond
	}
	c.maxErrPerSecond=maxErrPerSecond
	return nil
}
func (c *client)GetMaxErrPerSecond()int{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.maxErrPerSecond
}
func (c *client)SetMaxErrHeartbeat(maxErrHeartbeat int)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxErrHeartbeat<=0{
		return ErrSetMaxErrHeartbeat
	}
	c.maxErrHeartbeat=maxErrHeartbeat
	return nil
}
func (c *client)GetMaxErrHeartbeat()int{
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.maxErrHeartbeat
}
func (c *client)CodecName()string {
	return FuncsCodecName(c.funcsCodecType)
}
func (c *client)CodecType()CodecType {
	return c.funcsCodecType
}
func (c *client)Call(name string, args interface{}, reply interface{}) ( err error) {
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
func (c *client)CallNoRequest(name string, reply interface{}) ( err error) {
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
func (c *client)CallNoResponse(name string, args interface{}) ( err error) {
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

func (c *client)OnlyCall(name string) ( err error) {
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

func (c *client)Ping() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("*client.Ping failed:", err)
		}
	}()
	err:=c.heartbeat()
	if err!=nil{
		return false
	}
	return true
}
func (c *client)DisableRetry() {
	c.retry=false
}
func (c *client)RemoteCall(b []byte)([]byte,error){
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
func (c *client)RemoteCallNoResponse(b []byte)(error){
	c.pipelineChan<-true
	cbChan := make(chan []byte,1)
	c.pipeline.pipelineRequestChan<-NewPipelineRequest(b,true,cbChan)
	<-cbChan
	<-c.pipelineChan
	return nil
}
func (c *client)Disconnect() ( err error) {
	if !c.conn.Closed(){
		log.Traceln(c.client_id,"client conn Closed",c.conn.Closed())
		c.stopChan<-true
		time.Sleep(time.Millisecond*200)
	}
	c.disconnect = true
	return c.conn.Close()
}
func (c *client)Close() ( err error) {
	if c.closed {
		return nil
	}
	c.closed = true
	c.retry=false
	if c.batch!=nil{
		c.batch.Close()
	}
	if c.pipeline!=nil{
		c.pipeline.Close()
	}
	close(c.writeChan)
	close(c.readChan)
	close(c.stopChan)
	close(c.errCntChan)
	return c.conn.Close()
}
func (c *client)Closed()bool{
	return c.closed
}

func (c *client)heartbeat() ( err error) {
	uid:=c.idgenerator.GenUniqueIDInt64()
	msg:=&Msg{}
	msg.id=uid
	msg.msgType=MsgType(MsgTypeHea)
	msg_bytes, _ :=msg.Encode()
	ch := make(chan int)
	go func() {
		var data []byte
		data,err=c.RemoteCall(msg_bytes)
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

func (c *client)call(name string, args interface{}, reply interface{}) ( err error) {
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
func (c *client)batchCall(name string, args interface{}, reply interface{}) ( err error) {
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
