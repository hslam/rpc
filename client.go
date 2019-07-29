package rpc
import (
	"sync"
	"errors"
	"hslam.com/mgit/Mort/idgenerator"
	"time"
)
type Client struct {
	mu 				sync.Mutex
	transporter		Transporter
	closed			bool
	batchEnabled	bool
	batch			*Batch
	concurrent 		*Concurrent
	readChan		chan []byte
	writeChan		chan []byte
	stopChan		chan bool
	funcsCodecType 	CodecType
	idgenerator		*idgenerator.IDGen
	client_id		int64
}
func NewClient(transporter Transporter,codec string)  (*Client, error) {
	return NewClientnWithConcurrent(transporter,codec,DefultMaxConcurrentRequest)
}
func NewClientnWithConcurrent(transporter	Transporter,codec string,maxConcurrentRequest int)  (*Client, error)  {
	funcsCodecType,err:=FuncsCodecType(codec)
	if err!=nil{
		return nil,err
	}
	readChan := make(chan []byte,maxConcurrentRequest)
	writeChan := make(chan []byte,maxConcurrentRequest)
	stopChan := make(chan bool)
	var client_id int64=0
	idgenerator:=idgenerator.NewSnowFlake(client_id)
	conn :=  &Client{idgenerator:idgenerator,client_id:client_id,transporter:transporter,stopChan:stopChan,readChan:readChan,writeChan:writeChan,funcsCodecType:funcsCodecType}
	transporter.Handle(readChan,writeChan,stopChan)
	go func() {
		select {
		case <-conn.stopChan:
			defer conn.Close()
			close(writeChan)
			close(readChan)
			close(stopChan)
		}
	}()
	conn.concurrent=NewConcurrent(maxConcurrentRequest,readChan,writeChan)
	return conn, nil
}
func (c *Client)EnabledBatch(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batchEnabled = true
	c.batch=NewBatch(c,DefultMaxDelayNanoSecond*c.transporter.TickerFactor())
	c.batch.SetMaxBatchRequest(DefultMaxBatchRequest*c.transporter.BatchFactor())

}
func (c *Client)SetID(id int64)error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if id<0{
		return errors.New("0<=ClientID<=1023")
	}else if id>1023{
		return errors.New("0<=ClientID<=1023")
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
func (c *Client)GetMaxBatchRequest()int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.batch.GetMaxBatchRequest()
}
func (c *Client)GetMaxConcurrentRequest()(int){
	return c.concurrent.GetMaxConcurrentRequest()
}
func (c *Client)SetMaxBatchRequest(maxBatchRequest int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batch.SetMaxBatchRequest(maxBatchRequest)
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
			Panicln("Call failed:", err)
		}
	}()
	if c.batchEnabled{
		return c.batchCall(name,args,reply)
	}
	return c.call(name,args,reply)
}
func (c *Client)CallNoResponse(name string, args interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			Errorln("Call failed:", err)
		}
	}()
	if c.batchEnabled{
		return c.batchCallNoResponse(name,args)
	}
	return c.callNoResponse(name,args)
}
func (c *Client)RemoteCall(b []byte)([]byte,error){
	cbChan := make(chan []byte,1)
	c.concurrent.concurrentChan<-NewConcurrentRequest(b,cbChan)
	data,ok := <-cbChan
	if ok{
		return data,nil
	}
	return nil,errors.New("c.readChan is close")
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
func (c *Client)call(name string, args interface{}, reply interface{}) ( err error) {
	clientCodec:=&ClientCodec{}
	clientCodec.client_id=c.client_id
	clientCodec.req_id=uint64(c.idgenerator.GenUniqueIDInt64())
	clientCodec.name=name
	clientCodec.args=args
	clientCodec.noResponse=false
	clientCodec.funcsCodecType=c.funcsCodecType
	rpc_req_bytes, _ :=clientCodec.Encode()
	data,err:=c.RemoteCall(rpc_req_bytes)
	if err != nil {
		Errorln("Write error: ", err)
		return err
	}
	clientCodec.reply=reply
	return clientCodec.Decode(data)
}

func (c *Client)batchCall(name string, args interface{}, reply interface{}) ( err error) {
	reply_bytes := make(chan []byte, 1)
	reply_error := make(chan error, 1)
	args_bytes,err:=ArgsEncode(args,c.funcsCodecType)
	if err!=nil{
		Errorln("ArgsEncode error: ", err)
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
	case <-time.After(time.Second * 60):
		close(reply_bytes)
		close(reply_error)
		return errors.New("time out")
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
		Errorln("Write error: ", err)
		return err
	}
	return nil
}
func (c *Client)batchCallNoResponse(name string, args interface{}) ( err error) {
	args_bytes,err:=ArgsEncode(args,c.funcsCodecType)
	if err!=nil{
		Errorln("ArgsEncode error: ", err)
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

