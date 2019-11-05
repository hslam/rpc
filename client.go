package rpc
import (
	"sync"
	"hslam.com/git/x/idgenerator"
	"time"
	"math/rand"
	"hslam.com/git/x/rpc/log"
	"sync/atomic"
)
type Client interface {
	SetMaxRequests(max int)
	GetMaxRequests()(int)
	EnablePipelining()
	EnableMultiplexing()
	EnableBatch()
	EnableBatchAsync()
	SetMaxBatchRequest(max int)error
	GetMaxBatchRequest()int
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
	Go(name string, args interface{}, reply interface{}, done chan *Call) *Call
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
func DialWithMaxRequests(network,address,codec string,max int) (Client, error) {
	transporter,err:=dial(network,address)
	if err!=nil{
		return nil,err
	}
	return NewClientWithMaxRequests(transporter,codec,max)
}
type client struct {
	mu 					sync.RWMutex
	mutex				sync.RWMutex
	reqMutex 			sync.Mutex // protects following
	conn				Conn
	seq 				int64
	pending  			map[int64]*Call
	unordered 			bool
	closed				bool
	closing 			bool
	closeChan 			chan bool
	disconnect			bool
	hystrix				bool
	batchEnabled		bool
	batchAsync			bool
	batch				*Batch
	io					IO
	maxRequests			int
	requestChan			chan bool
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
	return NewClientWithMaxRequests(conn,codec,DefaultMaxRequests)
}
func NewClientWithMaxRequests(conn	Conn,codec string,max int)  (*client, error)  {
	funcsCodecType,err:=FuncsCodecType(codec)
	if err!=nil{
		return nil,err
	}
	var client_id int64=1
	c :=  &client{
		client_id:client_id,
		idgenerator:idgenerator.NewSnowFlake(client_id),
		conn:conn,
		pending: make(map[int64]*Call),
		finishChan:make(chan bool,1),
		stopChan:make(chan bool,1),
		closeChan:make(chan bool,1),
		funcsCodecType:funcsCodecType,
		compressLevel:NoCompression,
		compressType:CompressTypeNocom,
		retry:true,
		errCntChan:make(chan int,1000000),
		timeout:DefaultClientTimeout,
		maxErrPerSecond:DefaultClientMaxErrPerSecond,
		heartbeatTimeout:DefaultClientHearbeatTimeout,
		maxErrHeartbeat:DefaultClientMaxErrHearbeat,
	}
	c.setMaxRequests(max)
	c.conn.Handle(c.readChan,c.writeChan,c.stopChan,c.finishChan)
	c.enablePipelining()
	go c.run()
	return c, nil
}
func (c *client)run(){
	go func() {
		for i :=range c.errCntChan{
			c.errCnt+=i
		}
	}()
	time.Sleep(time.Millisecond*time.Duration(rand.Int63n(800)))
	ticker:=time.NewTicker(time.Second)
	heartbeatTicker:=time.NewTicker(time.Millisecond*DefaultClientHearbeatTicker)
	retryTicker:=time.NewTicker(time.Millisecond*DefaultClientRetryTicker)
	for{
		select {
		case <-c.finishChan:
			//log.Traceln(c.client_id,"client.run finishChan")
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
		case <-c.closeChan:
			func() {
				defer func() {if err := recover(); err != nil {}}()
				c.stopChan<-true
			}()
			close(c.closeChan)
			ticker.Stop()
			ticker=nil
			retryTicker.Stop()
			retryTicker=nil
			heartbeatTicker.Stop()
			heartbeatTicker=nil
			goto endfor

		case <-ticker.C:
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
		c.io.Retry()
	}
}
func (c *client)SetMaxRequests(max int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setMaxRequests(max)
}

func (c *client)setMaxRequests(max int) {
	c.maxRequests=max+1
	c.readChan=make(chan []byte,c.maxRequests)
	c.writeChan= make(chan []byte,c.maxRequests)
	c.requestChan=make(chan bool,c.maxRequests)
	if c.io!=nil{
		c.io.ResetMaxRequests(c.maxRequests)
		c.io.Reset(c.readChan,c.writeChan)
	}
}
func (c *client)GetMaxRequests()(int){
	return c.maxRequests-1
}
func (c *client)EnablePipelining(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enablePipelining()
}
func (c *client)enablePipelining(){
	c.unordered=false
	if c.maxRequests==0||c.readChan==nil||c.writeChan==nil{
		c.setMaxRequests(DefaultMaxRequests)
	}
	if c.io!=nil{
		c.io.Close()
	}
	c.io=NewPipeline(c.maxRequests,c.readChan,c.writeChan)
}
func (c *client)EnableMultiplexing(){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enableMultiplexing()
}
func (c *client)enableMultiplexing(){
	c.unordered=true
	if c.maxRequests==0||c.readChan==nil||c.writeChan==nil{
		c.setMaxRequests(DefaultMaxRequests)
	}
	if c.io!=nil{
		c.io.Close()
	}
	c.io=NewMultiplex(c,c.maxRequests,c.readChan,c.writeChan)
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
func (c *client)SetMaxBatchRequest(max int)error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if max<=0{
		return ErrSetMaxBatchRequest
	}
	c.batch.SetMaxBatchRequest(max)
	return nil
}
func (c *client)GetMaxBatchRequest()int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.batch.GetMaxBatchRequest()
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
func (c *client) Go(name string, args interface{}, reply interface{}, done chan *Call) *Call {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Call failed:", err)
		}
	}()
	call := new(Call)
	call.start=time.Now()
	call.ServiceMethod = name
	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *Call, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	if c.hystrix{
		call.Error=ErrHystrix
		call.done()
		return call
	}
	if c.unordered{
		go func(call *Call) {
			err:=c.Call(call.ServiceMethod,call.Args,call.Reply)
			if err!=nil{
				call.Error=err
			}
			call.done()
		}(call)
	}else {
		c.send(call)
	}
	return call
}
func (c *client)Call(name string, args interface{}, reply interface{}) ( err error) {
	if !c.unordered{
		call := <-c.Go(name, args, reply, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Call failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
	}
	err= c.call(name,args,reply)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}
func (c *client)CallNoRequest(name string, reply interface{}) ( err error) {
	if !c.unordered{
		call := <-c.Go(name, nil, reply, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Call failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
	}
	err= c.call(name,nil,reply)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}
func (c *client)CallNoResponse(name string, args interface{}) ( err error) {
	if !c.unordered{
		call := <-c.Go(name, args, nil, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("CallNoResponse failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
	}
	err= c.call(name,args,nil)
	if err!=nil{
		c.errCntChan<-1
	}
	return
}

func (c *client)OnlyCall(name string) ( err error) {
	if !c.unordered{
		call := <-c.Go(name, nil, nil, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("OnlyCall failed:", err)
		}
	}()
	if c.hystrix{
		return ErrHystrix
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
			log.Errorln("client.Ping failed:", err)
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
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("client.RemoteCall failed:", err)
		}
	}()
	c.requestChan<-true
	cbChan := make(chan []byte,1)
	c.io.RequestChan()<-c.io.NewRequest(0,b,false,cbChan)
	data,ok := <-cbChan
	<-c.requestChan
	if ok{
		return data,nil
	}
	return nil,ErrRemoteCall
}
func (c *client)RemoteCallNoResponse(b []byte)(error){
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("client.RemoteCallNoResponse failed:", err)
		}
	}()
	c.requestChan<-true
	cbChan := make(chan []byte,1)
	c.io.RequestChan()<-c.io.NewRequest(0,b,true,cbChan)
	<-cbChan
	<-c.requestChan
	return nil
}
func (c *client)Disconnect() ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("client.Disconnect failed:", err)
		}
	}()
	if !c.conn.Closed(){
		log.Traceln(c.client_id,"client conn Closed",c.conn.Closed())
		c.stopChan<-true
		time.Sleep(time.Millisecond*200)
	}
	c.disconnect = true
	return c.conn.Close()
}
func (c *client)Close() ( err error) {
	c.closing=true
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("client.Close failed:", err)
		}
		c.closing=false
	}()
	if c.closed {
		return nil
	}
	c.closed = true
	c.retry=false
	if c.batch!=nil{
		c.batch.Close()
	}
	if c.io!=nil{
		c.io.Close()
	}
	c.closeChan<-true
	close(c.writeChan)
	close(c.readChan)
	close(c.stopChan)
	close(c.finishChan)
	close(c.errCntChan)
	close(c.requestChan)
	c.idgenerator=nil
	return c.conn.Close()
}
func (c *client)Closed()bool{
	return c.closed||c.closing
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
	if c.batchEnabled{
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
			if c.timeout>0{
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
			}else {
				select {
				case bytes ,ok:= <- cr.reply_bytes:
					if ok{
						return ReplyDecode(bytes,reply,c.funcsCodecType)
					}
				case err ,ok:= <- cr.reply_error:
					if ok{
						return err
					}
				}
				return nil
			}
		}
	}else {
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
				data, err = c.RemoteCall(rpc_req_bytes)
				if err != nil {
					log.Errorln("Write error: ", err)
					ch <- 1
					return
				}
				clientCodec.reply = reply
				err = clientCodec.Decode(data)
				ch <- 1
			}()
			if c.timeout > 0 {
				select {
				case <-ch:
				case <-time.After(time.Millisecond * time.Duration(c.timeout)):
					err = ErrTimeOut
				}
				return err
			} else {
				select {
				case <-ch:
				}
				return nil
			}
		}
	}

}
func (c *client)Seq()int64{
	seq := atomic.AddInt64(&c.seq,1)
	return seq
}
func (c *client) send(call *Call) {
	c.reqMutex.Lock()
	defer c.reqMutex.Unlock()
	// Register this call.
	if c.closed || c.closing {
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := c.Seq()
	// Encode and send the request.
	if !c.batchEnabled{
		clientCodec:=&ClientCodec{}
		clientCodec.client_id=c.client_id
		clientCodec.req_id=uint64(seq)
		clientCodec.name=call.ServiceMethod
		if call.Args!=nil{
			clientCodec.args=call.Args
			clientCodec.noRequest=false
		}else {
			clientCodec.args=nil
			clientCodec.noRequest=true
		}
		clientCodec.funcsCodecType=c.funcsCodecType
		if call.Reply!=nil{
			clientCodec.noResponse=false
		}else {
			clientCodec.noResponse=true
		}
		rpc_req_bytes, err :=clientCodec.Encode()
		if err != nil {
			if call != nil {
				call.Error = err
				call.done()
				return
			}
		}
		if clientCodec.noResponse{
			err=c.RemoteCallNoResponse(rpc_req_bytes)
			if call != nil {
				if err != nil {
					call.Error = err
				}
				call.done()
			}
			return
		}else {
			c.mutex.Lock()
			c.pending[seq] = call
			c.mutex.Unlock()
			go func(c *client,rpc_req_bytes []byte,call *Call) {
				var data []byte
				data, err = c.RemoteCall(rpc_req_bytes)
				if err != nil {
					log.Errorln("Write error: ", err)
					return
				}
				clientCodec.reply = call.Reply
				err = clientCodec.Decode(data)
				c.mutex.Lock()
				var ok bool
				if call ,ok= c.pending[seq];ok{
					delete(c.pending, seq)
				}
				c.mutex.Unlock()
				if err != nil {
					if call != nil {
						call.Error = err
						call.done()
						return
					}
				}
				if call != nil {
					call.Reply=clientCodec.reply
					call.Done<-call
					return
				}
			}(c,rpc_req_bytes,call)
		}
	}else {
		cr:=&BatchRequest{
			id:uint64(c.idgenerator.GenUniqueIDInt64()),
			name:call.ServiceMethod,
			noResponse:false,
		}
		if call.Args!=nil{
			args_bytes,err:=ArgsEncode(call.Args,c.funcsCodecType)
			if err!=nil{
				log.Errorln("ArgsEncode error: ", err)
			}
			cr.args_bytes=args_bytes
			cr.noRequest=false
		}else {
			cr.args_bytes=nil
			cr.noRequest=true
		}
		if call.Reply!=nil{
			cr.reply_bytes= make(chan []byte, 1)
			cr.reply_error=make(chan error, 1)
			cr.noResponse=false
		}else {
			cr.noResponse=true
		}
		c.batch.reqChan<-cr
		if cr.noResponse{
			if call != nil {
				call.done()
			}
			return
		}else{
			c.mutex.Lock()
			c.pending[seq] = call
			c.mutex.Unlock()
			go func(cr *BatchRequest,call *Call) {
				select {
				case bytes ,ok:= <- cr.reply_bytes:
					if ok{
						err:= ReplyDecode(bytes,call.Reply,c.funcsCodecType)
						c.mutex.Lock()
						var ok bool
						if call ,ok= c.pending[seq];ok{
							delete(c.pending, seq)
						}
						c.mutex.Unlock()
						if err != nil {
							if call != nil {
								call.Error = err
								call.done()
								return
							}
						}
						call.Done<-call
						return
					}
				case err ,ok:= <- cr.reply_error:
					if ok{
						if err != nil {
							c.mutex.Lock()
							var ok bool
							if call ,ok= c.pending[seq];ok{
								delete(c.pending, seq)
							}
							c.mutex.Unlock()
							if call != nil {
								call.Error = err
								call.done()
								return
							}
						}
					}
				}
			}(cr,call)
		}
	}
}
