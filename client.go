package rpc
import (
	"sync"
	"hslam.com/mgit/Mort/idgenerator"
	"time"
	"math/rand"
	"hslam.com/mgit/Mort/rpc/log"
)
type Client struct {
	mu 				sync.Mutex
	transporter		Transporter
	closed			bool
	hystrix			bool
	batchEnabled	bool
	batch			*Batch
	concurrent 		*Concurrent
	readChan		chan []byte
	writeChan		chan []byte
	stopChan		chan bool
	funcsCodecType 	CodecType
	idgenerator		*idgenerator.IDGen
	client_id		int64
	timeout 		int64
	heartbeatTimeout 		int64
	errCntChan 		chan int
	errCnt 			int
	maxErrPerSecond	int
	maxErrHeartbeat	int
	errCntHeartbeat	int
}
func NewClient(transporter Transporter,codec string)  (*Client, error) {
	return NewClientnWithConcurrent(transporter,codec,DefaultMaxConcurrentRequest)
}
func NewClientnWithConcurrent(transporter	Transporter,codec string,maxConcurrentRequest int)  (*Client, error)  {
	funcsCodecType,err:=FuncsCodecType(codec)
	if err!=nil{
		return nil,err
	}

	stopChan := make(chan bool)
	errCntChan:=make(chan int,1000000)
	var client_id int64=0
	idgenerator:=idgenerator.NewSnowFlake(client_id)
	conn :=  &Client{client_id:client_id,transporter:transporter,stopChan:stopChan,funcsCodecType:funcsCodecType}
	conn.readChan=make(chan []byte,maxConcurrentRequest)
	conn.writeChan= make(chan []byte,maxConcurrentRequest)
	conn.idgenerator=idgenerator
	conn.errCntChan=errCntChan
	conn.timeout=DefaultClientTimeout
	conn.maxErrPerSecond=DefaultClientMaxErrPerSecond
	conn.heartbeatTimeout=DefaultClientHearbeatTimeout
	conn.maxErrHeartbeat=DefaultClientMaxErrHearbeat
	conn.transporter.Handle(conn.readChan,conn.writeChan,conn.stopChan)
	go func() {
		defer func() {
			defer conn.Close()
			close(conn.writeChan)
			close(conn.readChan)
			close(stopChan)
			close(errCntChan)
		}()
		for{
			select {
			case <-conn.stopChan:
				conn.Close()
			}
		}

	}()
	go func() {
		for c :=range conn.errCntChan{
			conn.errCnt+=c
		}
	}()
	go func() {
		time.Sleep(time.Millisecond*time.Duration(rand.Int63n(800)))
		ticker:=time.NewTicker(time.Second)
		heartbeatTicker:=time.NewTicker(time.Second*DefaultClientHearbeatTicker)
		retryTicker:=time.NewTicker(time.Second*DefaultClientRetryTicker)

		for{
			select {
			case <-heartbeatTicker.C:
				if conn.closed==false{
					err:=conn.heartbeat()
					if err!=nil{
						conn.errCntHeartbeat+=1
					}else {
					}
					if conn.errCntHeartbeat>=conn.maxErrHeartbeat{
						conn.Close()
						conn.errCntHeartbeat=0
					}
				}
			case <-retryTicker.C:
				if conn.closed==true{
					conn.Close()
					err:=conn.transporter.Retry()
					if err!=nil{
						conn.hystrix=true
						log.Errorln(conn.client_id,"retry connection err ",err)
					}else {
						conn.closed=false
						conn.hystrix=false
						conn.concurrent.retry()
					}
				}
			case <-ticker.C:
				if conn.errCnt>=DefaultClientMaxErrPerSecond{
					conn.hystrix=true
					conn.Close()
				}else {
					conn.hystrix=false
				}
				conn.errCnt-=conn.errCnt
			}
		}
	}()
	conn.concurrent=NewConcurrent(maxConcurrentRequest,conn.readChan,conn.writeChan,conn)
	return conn, nil
}
func (c *Client)EnabledBatch(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batchEnabled = true
	c.batch=NewBatch(c,DefaultMaxDelayNanoSecond*c.transporter.TickerFactor())
	c.batch.SetMaxBatchRequest(DefaultMaxBatchRequest*c.transporter.BatchFactor())

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
func (c *Client)GetMaxBatchRequest()int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.batch.GetMaxBatchRequest()
}
func (c *Client)GetMaxConcurrentRequest()(int){
	return c.concurrent.GetMaxConcurrentRequest()-1
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
		err= c.batchCallNoResponse(name,args)
		if err!=nil{
			c.errCntChan<-1
		}
		return
	}
	err= c.callNoResponse(name,args)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}
func (c *Client)RemoteCall(b []byte)([]byte,error){
	cbChan := make(chan []byte,1)
	c.concurrent.concurrentChan<-NewConcurrentRequest(b,cbChan)
	data,ok := <-cbChan
	if ok{
		return data,nil
	}
	return nil,ErrRemoteCall
}
func (c *Client)RemoteCallNoResponse(b []byte)(error){
	c.concurrent.concurrentChan<-NewConcurrentRequest(b,nil)
	return nil
}
func (c *Client)Close() ( err error) {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.transporter.Close()
}
func (c *Client)Closed()bool{
	return c.closed
}
func (c *Client)call(name string, args interface{}, reply interface{}) ( err error) {
	clientCodec:=&ClientCodec{}
	clientCodec.client_id=c.client_id
	clientCodec.req_id=uint64(c.idgenerator.GenUniqueIDInt64())
	clientCodec.name=name
	clientCodec.args=args
	clientCodec.noResponse=false
	clientCodec.funcsCodecType=c.funcsCodecType
	rpc_req_bytes, _ :=clientCodec.Encode()
	ch := make(chan int)
	go func() {
		data,err:=c.RemoteCall(rpc_req_bytes)
		if err != nil {
			log.Errorln("Write error: ", err)
			ch<-1
			return
		}
		clientCodec.reply=reply
		err=clientCodec.Decode(data)
		if clientCodec.res!=nil&&err==nil{
			if clientCodec.res.err!=nil{
				err=clientCodec.res.err
			}
		}
		ch<-1
	}()
	select {
	case <-ch:
	case <-time.After(time.Second * time.Duration(c.timeout)):
		err=ErrTimeOut
	}
	return err
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
	case <-time.After(time.Second * time.Duration(c.heartbeatTimeout)):
		err=ErrTimeOut
	}
	return err
}
func (c *Client)batchCall(name string, args interface{}, reply interface{}) ( err error) {
	reply_bytes := make(chan []byte, 1)
	reply_error := make(chan error, 1)
	args_bytes,err:=ArgsEncode(args,c.funcsCodecType)
	if err!=nil{
		log.Errorln("ArgsEncode error: ", err)
	}
	cr:=&BatchRequest{
		id:uint64(c.idgenerator.GenUniqueIDInt64()),
		name:name,
		args_bytes:args_bytes,
		reply_bytes:reply_bytes,
		reply_error:reply_error,
		noResponse:false,
	}
	c.batch.reqChan<-cr
	select {
	case bytes ,ok:= <- reply_bytes:
		if ok{
			return ReplyDecode(bytes,reply,c.funcsCodecType)
		}
	case err ,ok:= <- reply_error:
		if ok{
			return err
		}
	case <-time.After(time.Second * time.Duration(c.timeout)):
		close(reply_bytes)
		close(reply_error)
		return ErrTimeOut
	}
	return nil
}
func (c *Client)callNoResponse(name string, args interface{}) ( err error) {
	clientCodec:=&ClientCodec{}
	clientCodec.client_id=c.client_id
	clientCodec.req_id=uint64(c.idgenerator.GenUniqueIDInt64())
	clientCodec.name=name
	clientCodec.args=args
	clientCodec.noResponse=true
	clientCodec.funcsCodecType=c.funcsCodecType
	rpc_req_bytes, _ :=clientCodec.Encode()
	err=c.RemoteCallNoResponse(rpc_req_bytes)
	if err != nil {
		log.Errorln("Write error: ", err)
		return err
	}
	return nil
}
func (c *Client)batchCallNoResponse(name string, args interface{}) ( err error) {
	args_bytes,err:=ArgsEncode(args,c.funcsCodecType)
	if err!=nil{
		log.Errorln("ArgsEncode error: ", err)
	}
	cr:=&BatchRequest{
		id:uint64(c.idgenerator.GenUniqueIDInt64()),
		name:name,
		args_bytes:args_bytes,
		reply_bytes:nil,
		reply_error:nil,
		noResponse:true,
	}
	c.batch.reqChan<-cr
	return nil
}

