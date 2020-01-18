package rpc

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

//Client defines the interface of client.
type Client interface {
	GetMaxRequests() int
	GetMaxBatchRequest() int
	GetID() uint64
	GetTimeout() int64
	GetHeartbeatTimeout() int64
	GetMaxErrPerSecond() int
	GetMaxErrHeartbeat() int
	CodecName() string
	CodecType() CodecType
	Go(name string, args interface{}, reply interface{}, done chan *Call) *Call
	Call(name string, args interface{}, reply interface{}) (err error)
	CallNoRequest(name string, reply interface{}) (err error)
	CallNoResponse(name string, args interface{}) (err error)
	OnlyCall(name string) (err error)
	Ping() bool
	Close() (err error)
	Closed() bool
}

// Dial connects to an RPC server at the specified network address codec
// and returns a client.
func Dial(network, address, codec string) (Client, error) {
	transporter, err := dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(transporter, codec)
}

// DialWithOptions is like Dial but uses the specified options
// and returns a client.
func DialWithOptions(network, address, codec string, opts *Options) (Client, error) {
	transporter, err := dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(transporter, codec, opts)
}

type client struct {
	mu               sync.RWMutex
	mutex            sync.RWMutex
	conn             Conn
	seq              uint64
	pending          map[int64]*Call
	unordered        bool
	closed           bool
	closing          bool
	closeChan        chan bool
	disconnected     bool
	hystrix          bool
	batchEnabled     bool
	batchingAsync    bool
	batch            *batch
	io               IO
	maxRequests      int
	requestChan      chan bool
	readChan         chan []byte
	writeChan        chan []byte
	finishChan       chan bool
	stopChan         chan bool
	funcsCodecType   CodecType
	compressType     CompressType
	compressLevel    CompressLevel
	clientID         uint64
	timeout          int64
	heartbeatTimeout int64
	errCntChan       chan int
	errCnt           int
	maxErrPerSecond  int
	maxErrHeartbeat  int
	errCntHeartbeat  int
	retry            bool
	noDelay          bool
}

// NewClient returns a new Client to handle requests.
func NewClient(conn Conn, codec string) (Client, error) {
	return NewClientWithOptions(conn, codec, DefaultOptions())
}

// NewClientWithOptions is like NewClient but uses the specified options.
func NewClientWithOptions(conn Conn, codec string, opts *Options) (Client, error) {
	funcsCodecType, err := funcsCodecType(codec)
	if err != nil {
		return nil, err
	}
	var clientID uint64 = 1
	c := &client{
		clientID:         clientID,
		conn:             conn,
		pending:          make(map[int64]*Call),
		finishChan:       make(chan bool, 1),
		stopChan:         make(chan bool, 1),
		closeChan:        make(chan bool, 1),
		funcsCodecType:   funcsCodecType,
		compressLevel:    NoCompression,
		compressType:     CompressTypeNo,
		retry:            true,
		errCntChan:       make(chan int, 1000000),
		timeout:          DefaultClientTimeout,
		maxErrPerSecond:  DefaultClientMaxErrPerSecond,
		heartbeatTimeout: DefaultClientHearbeatTimeout,
		maxErrHeartbeat:  DefaultClientMaxErrHeartbeat,
	}
	c.loadOptions(opts)
	c.conn.NoDelay(c.noDelay)
	c.conn.Multiplexing(c.unordered)
	c.conn.Handle(c.readChan, c.writeChan, c.stopChan, c.finishChan)
	go c.run()
	return c, nil
}
func (c *client) run() {
	go func() {
		for i := range c.errCntChan {
			c.errCnt += i
		}
	}()
	time.Sleep(time.Millisecond * time.Duration(rand.Int63n(800)))
	ticker := time.NewTicker(time.Second)
	heartbeatTicker := time.NewTicker(time.Millisecond * DefaultClientHearbeatTicker)
	retryTicker := time.NewTicker(time.Millisecond * DefaultClientRetryTicker)
	for {
		select {
		case <-c.finishChan:
			//logger.Traceln(c.client_id,"client.run finishChan")
			c.retryConnect()
		case <-heartbeatTicker.C:
			if c.disconnected == false && c.retry {
				err := c.heartbeat()
				if err != nil {
					c.errCntHeartbeat++
				} else {
					c.errCntHeartbeat = 0
				}
				if c.errCntHeartbeat >= c.maxErrHeartbeat {
					c.retryConnect()
					c.errCntHeartbeat = 0
				}
			}
		case <-retryTicker.C:
			if c.disconnected == true && c.retry {
				c.retryConnect()
			}
		case <-c.closeChan:
			func() {
				defer func() {
					if err := recover(); err != nil {
					}
				}()
				c.stopChan <- true
			}()
			close(c.closeChan)
			ticker.Stop()
			ticker = nil
			retryTicker.Stop()
			retryTicker = nil
			heartbeatTicker.Stop()
			heartbeatTicker = nil
			goto endfor

		case <-ticker.C:
			if !c.retry {
				return
			}
			if c.errCnt >= DefaultClientMaxErrPerSecond {
				c.hystrix = true
			} else {
				c.hystrix = false
			}
			c.errCnt = 0
		}
	}
endfor:
}
func (c *client) retryConnect() {
	if !c.retry {
		return
	}
	c.disconnect()
	err := c.conn.Retry()
	if err != nil {
		c.hystrix = true
		logger.Traceln(c.clientID, "retry connection err ", err)
	} else {
		logger.Traceln(c.clientID, "retry connection success")
		c.disconnected = false
		c.hystrix = false
		c.io.Retry()
	}
}
func (c *client) loadOptions(opts *Options) {
	opts.Check()
	c.setMaxRequests(opts.MaxRequests)
	if opts.Pipelining {
		c.enablePipelining()
	} else if opts.Multiplexing {
		c.enableMultiplexing()
	} else {
		c.enablePipelining()
	}
	if opts.Batching {
		c.enableBatching()
		if opts.Multiplexing {
			c.enableBatchingAsync()
		}
		if opts.MaxBatchRequest > c.GetMaxBatchRequest() {
			c.setMaxBatchRequest(opts.MaxBatchRequest)
		}
	}
	c.setCompressType(opts.CompressType)
	c.setCompressLevel(opts.CompressLevel)
	c.setID(opts.ID)
	c.setTimeout(opts.Timeout)
	c.setHeartbeatTimeout(opts.HeartbeatTimeout)
	c.setMaxErrPerSecond(opts.MaxErrPerSecond)
	c.setMaxErrHeartbeat(opts.MaxErrHeartbeat)
	if !opts.Retry {
		c.disableRetry()
	}
	c.setNoDelay(opts.noDelay)
}
func (c *client) setNoDelay(enable bool) {
	c.noDelay = enable
}
func (c *client) setMaxRequests(max int) {
	c.maxRequests = max + 1
	c.readChan = make(chan []byte, c.maxRequests)
	c.writeChan = make(chan []byte, c.maxRequests)
	c.requestChan = make(chan bool, c.maxRequests)
	if c.io != nil {
		c.io.ResetMaxRequests(c.maxRequests)
		c.io.Reset(c.readChan, c.writeChan)
	}
}

//GetMaxRequests returns the number of max requests.
func (c *client) GetMaxRequests() int {
	return c.maxRequests - 1
}

func (c *client) enablePipelining() {
	c.unordered = false
	if c.maxRequests == 0 || c.readChan == nil || c.writeChan == nil {
		c.setMaxRequests(DefaultMaxRequests)
	}
	if c.io != nil {
		c.io.Close()
	}
	c.io = newPipeline(c.maxRequests, c.readChan, c.writeChan)
}

func (c *client) enableMultiplexing() {
	c.unordered = true
	if c.maxRequests == 0 || c.readChan == nil || c.writeChan == nil {
		c.setMaxRequests(DefaultMaxRequests)
	}
	if c.io != nil {
		c.io.Close()
	}
	c.io = newMultiplex(c, c.maxRequests, c.readChan, c.writeChan)
}
func (c *client) enableBatching() {
	c.batchEnabled = true
	c.batch = newBatch(c, DefaultMaxDelayNanoSecond*c.conn.TickerFactor(), DefaultMaxBatchRequest*c.conn.BatchFactor())
	c.batch.SetMaxBatchRequest(DefaultMaxBatchRequest * c.conn.BatchFactor())
}
func (c *client) enableBatchingAsync() {
	c.batchingAsync = true
}
func (c *client) setMaxBatchRequest(max int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if max <= 0 {
		return ErrSetMaxBatchRequest
	}
	c.batch.SetMaxBatchRequest(max)
	return nil
}

//GetMaxBatchRequest returns the number of max batch requests.
func (c *client) GetMaxBatchRequest() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.batch.GetMaxBatchRequest()
}

func (c *client) setCompressType(compressType CompressType) {
	c.compressType = compressType
	c.compressLevel = DefaultCompression
}
func (c *client) setCompressLevel(level CompressLevel) {
	c.compressLevel = level
}
func (c *client) setID(id uint64) error {
	c.clientID = id
	return nil
}

//GetID returns the client id.
func (c *client) GetID() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.clientID
}
func (c *client) setTimeout(timeout int64) error {
	c.timeout = timeout
	return nil
}

//GetTimeout returns the request timeout.
func (c *client) GetTimeout() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeout
}

func (c *client) setHeartbeatTimeout(timeout int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout <= 0 {
		return ErrSetTimeout
	}
	c.heartbeatTimeout = timeout
	return nil
}

//GetHeartbeatTimeout returns the heartbeat timeout.
func (c *client) GetHeartbeatTimeout() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.heartbeatTimeout
}
func (c *client) setMaxErrPerSecond(maxErrPerSecond int) error {
	if maxErrPerSecond <= 0 {
		return ErrSetMaxErrPerSecond
	}
	c.maxErrPerSecond = maxErrPerSecond
	return nil
}

//GetMaxErrPerSecond returns the number of max errors per second.
func (c *client) GetMaxErrPerSecond() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.maxErrPerSecond
}
func (c *client) setMaxErrHeartbeat(maxErrHeartbeat int) error {
	if maxErrHeartbeat <= 0 {
		return ErrSetMaxErrHeartbeat
	}
	c.maxErrHeartbeat = maxErrHeartbeat
	return nil
}

//GetMaxErrHeartbeat returns the number of max heartbeat errors.
func (c *client) GetMaxErrHeartbeat() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.maxErrHeartbeat
}

//CodecName returns the codec name.
func (c *client) CodecName() string {
	return funcsCodecName(c.funcsCodecType)
}

//CodecType returns the codec type.
func (c *client) CodecType() CodecType {
	return c.funcsCodecType
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (c *client) Go(name string, args interface{}, reply interface{}, done chan *Call) *Call {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("Go failed:", err)
		}
	}()
	call := new(Call)
	//call.start=time.Now()
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
			logger.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	if c.hystrix {
		call.Error = ErrHystrix
		call.done()
		return call
	}
	if c.unordered {
		go func(call *Call) {
			err := c.Call(call.ServiceMethod, call.Args, call.Reply)
			if err != nil {
				call.Error = err
			}
			call.done()
		}(call)
	} else {
		c.send(call)
	}
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *client) Call(name string, args interface{}, reply interface{}) (err error) {
	if !c.unordered {
		call := <-c.Go(name, args, reply, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("Call failed:", err)
		}
	}()
	if c.hystrix {
		return ErrHystrix
	}
	err = c.call(name, args, reply)
	if err != nil {
		c.errCntChan <- 1
	}
	return
}

// CallNoRequest invokes the named function but doesn't use args, waits for it to complete, and returns its error status.
func (c *client) CallNoRequest(name string, reply interface{}) (err error) {
	if !c.unordered {
		call := <-c.Go(name, nil, reply, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("Call failed:", err)
		}
	}()
	if c.hystrix {
		return ErrHystrix
	}
	err = c.call(name, nil, reply)
	if err != nil {
		c.errCntChan <- 1
	}
	return
}

// CallNoResponse invokes the named function but doesn't return reply, waits for it to complete, and returns its error status.
func (c *client) CallNoResponse(name string, args interface{}) (err error) {
	if !c.unordered {
		call := <-c.Go(name, args, nil, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("CallNoResponse failed:", err)
		}
	}()
	if c.hystrix {
		return ErrHystrix
	}
	err = c.call(name, args, nil)
	if err != nil {
		c.errCntChan <- 1
	}
	return
}

// OnlyCall invokes the named function but doesn't use args and doesn't return reply, waits for it to complete, and returns its error status.
func (c *client) OnlyCall(name string) (err error) {
	if !c.unordered {
		call := <-c.Go(name, nil, nil, make(chan *Call, 1)).Done
		return call.Error
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("OnlyCall failed:", err)
		}
	}()
	if c.hystrix {
		return ErrHystrix
	}
	err = c.call(name, nil, nil)
	if err != nil {
		c.errCntChan <- 1
	}
	return
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (c *client) Ping() bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("client.Ping failed:", err)
		}
	}()
	err := c.heartbeat()
	if err != nil {
		return false
	}
	return true
}
func (c *client) disableRetry() {
	c.retry = false
}
func (c *client) remoteCall(b []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("client.RemoteCall failed:", err)
		}
	}()
	c.requestChan <- true
	cbChan := make(chan []byte, 1)
	c.io.RequestChan() <- c.io.NewRequest(0, b, false, cbChan)
	data, ok := <-cbChan
	<-c.requestChan
	if ok {
		return data, nil
	}
	return nil, ErrRemoteCall
}
func (c *client) remoteGo(b []byte, cbChan chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("client.RemoteCall failed:", err)
		}
	}()
	c.requestChan <- true
	c.io.RequestChan() <- c.io.NewRequest(0, b, false, cbChan)
}
func (c *client) remoteCallNoResponse(b []byte) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("client.RemoteCallNoResponse failed:", err)
		}
	}()
	c.requestChan <- true
	cbChan := make(chan []byte, 1)
	c.io.RequestChan() <- c.io.NewRequest(0, b, true, cbChan)
	<-cbChan
	<-c.requestChan
	return nil
}
func (c *client) disconnect() (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("client.Disconnect failed:", err)
		}
	}()
	if !c.conn.Closed() {
		logger.Traceln(c.clientID, "client conn Closed", c.conn.Closed())
		c.stopChan <- true
		time.Sleep(time.Millisecond * 200)
	}
	c.disconnected = true
	return c.conn.Close()
}

// Close closes the connection
func (c *client) Close() (err error) {
	c.closing = true
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("client.Close failed:", err)
		}
		c.closing = false
	}()
	if c.closed {
		return nil
	}
	c.closed = true
	c.retry = false
	if c.batch != nil {
		c.batch.Close()
	}
	if c.io != nil {
		c.io.Close()
	}
	c.closeChan <- true
	close(c.writeChan)
	close(c.readChan)
	close(c.stopChan)
	close(c.finishChan)
	close(c.errCntChan)
	close(c.requestChan)
	return c.conn.Close()
}

// Closed returns the closed
func (c *client) Closed() bool {
	return c.closed || c.closing
}

func (c *client) heartbeat() (err error) {
	uid := c.getSeq()
	msg := &msg{}
	msg.id = uid
	msg.msgType = MsgType(MsgTypeHea)
	msgBytes, _ := msg.Encode()
	ch := make(chan int)
	go func() {
		var data []byte
		data, err = c.remoteCall(msgBytes)
		if err != nil {
			logger.Errorln("Write error: ", err)
			ch <- 1
			return
		}
		err = msg.Decode(data)
		if err == nil {
			if msg.id != uid {
				err = ErrClientID
			}
		}
		ch <- 1
		return
	}()
	select {
	case <-ch:
	case <-time.After(time.Millisecond * time.Duration(c.heartbeatTimeout)):
		err = ErrTimeOut
	}
	return err
}

func (c *client) call(name string, args interface{}, reply interface{}) (err error) {
	if c.batchEnabled {
		cr := &batchRequest{
			id:         uint64(c.getSeq()),
			name:       name,
			noResponse: false,
		}
		if args != nil {
			argsBytes, err := argsEncode(args, c.funcsCodecType)
			if err != nil {
				logger.Errorln("ArgsEncode error: ", err)
			}
			cr.argsBytes = argsBytes
			cr.noRequest = false
		} else {
			cr.argsBytes = nil
			cr.noRequest = true
		}
		if reply != nil {
			cr.replyBytes = make(chan []byte, 1)
			cr.replyError = make(chan error, 1)
			cr.noResponse = false
		} else {
			cr.noResponse = true
		}
		c.batch.reqChan <- cr
		if cr.noResponse {
			return nil
		}
		if c.timeout > 0 {
			select {
			case bytes, ok := <-cr.replyBytes:
				if ok {
					return replyDecode(bytes, reply, c.funcsCodecType)
				}
			case err, ok := <-cr.replyError:
				if ok {
					return err
				}
			case <-time.After(time.Millisecond * time.Duration(c.timeout)):
				close(cr.replyBytes)
				close(cr.replyError)
				return ErrTimeOut
			}
			return nil
		}
		select {
		case bytes, ok := <-cr.replyBytes:
			if ok {
				return replyDecode(bytes, reply, c.funcsCodecType)
			}
		case err, ok := <-cr.replyError:
			if ok {
				return err
			}
		}
		return nil
	}

	clientCodec := &clientCodec{}
	clientCodec.clientID = c.clientID
	clientCodec.reqID = uint64(c.getSeq())
	clientCodec.name = name
	clientCodec.compressType = c.compressType
	clientCodec.compressLevel = c.compressLevel
	clientCodec.funcsCodecType = c.funcsCodecType
	if args != nil {
		clientCodec.args = args
		clientCodec.noRequest = false
	} else {
		clientCodec.args = nil
		clientCodec.noRequest = true
	}
	if reply != nil {
		clientCodec.noResponse = false
	} else {
		clientCodec.noResponse = true
	}
	rpcReqBytes, _ := clientCodec.Encode()
	if clientCodec.noResponse {
		err = c.remoteCallNoResponse(rpcReqBytes)
		if err != nil {
			logger.Errorln("Write error: ", err)
			return err
		}
		return nil
	}
	ch := make(chan int)
	go func() {
		var data []byte
		data, err = c.remoteCall(rpcReqBytes)
		if err != nil {
			logger.Errorln("Write error: ", err)
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
	}
	select {
	case <-ch:
	}
	return nil

}
func (c *client) getSeq() uint64 {
	seq := atomic.AddUint64(&c.seq, 1)
	return seq
}
func (c *client) send(call *Call) {
	// Register this call.
	if c.closed || c.closing {
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := c.getSeq()
	// Encode and send the request.
	if !c.batchEnabled {
		clientCodec := &clientCodec{}
		clientCodec.clientID = c.clientID
		clientCodec.reqID = uint64(seq)
		clientCodec.name = call.ServiceMethod
		clientCodec.compressType = c.compressType
		clientCodec.compressLevel = c.compressLevel
		clientCodec.funcsCodecType = c.funcsCodecType
		if call.Args != nil {
			clientCodec.args = call.Args
			clientCodec.noRequest = false
		} else {
			clientCodec.args = nil
			clientCodec.noRequest = true
		}
		if call.Reply != nil {
			clientCodec.noResponse = false
		} else {
			clientCodec.noResponse = true
		}
		rpcReqBytes, err := clientCodec.Encode()
		if err != nil {
			if call != nil {
				call.Error = err
				call.done()
				return
			}
		}
		if clientCodec.noResponse {
			err = c.remoteCallNoResponse(rpcReqBytes)
			if call != nil {
				if err != nil {
					call.Error = err
				}
				call.done()
			}
			return
		}
		//c.mutex.Lock()
		//c.pending[seq] = call
		//c.mutex.Unlock()
		cbChan := make(chan []byte, 1)
		c.remoteGo(rpcReqBytes, cbChan)
		go func(c *client, cbChan chan []byte, call *Call) {
			var data []byte
			data, ok := <-cbChan
			<-c.requestChan
			if !ok {
				return
			}
			clientCodec.reply = call.Reply
			err = clientCodec.Decode(data)
			//c.mutex.Lock()
			//var ok bool
			//if call ,ok= c.pending[seq];ok{
			//	delete(c.pending, seq)
			//}
			//c.mutex.Unlock()
			if err != nil {
				if call != nil {
					call.Error = err
					call.done()
					return
				}
			}
			if call != nil {
				call.Reply = clientCodec.reply
				call.Done <- call
				return
			}
		}(c, cbChan, call)
	} else {
		cr := &batchRequest{
			id:         uint64(c.getSeq()),
			name:       call.ServiceMethod,
			noResponse: false,
		}
		if call.Args != nil {
			argsBytes, err := argsEncode(call.Args, c.funcsCodecType)
			if err != nil {
				logger.Errorln("ArgsEncode error: ", err)
			}
			cr.argsBytes = argsBytes
			cr.noRequest = false
		} else {
			cr.argsBytes = nil
			cr.noRequest = true
		}
		if call.Reply != nil {
			cr.replyBytes = make(chan []byte, 1)
			cr.replyError = make(chan error, 1)
			cr.noResponse = false
		} else {
			cr.noResponse = true
		}
		c.batch.reqChan <- cr
		if cr.noResponse {
			if call != nil {
				call.done()
			}
			return
		}
		//c.mutex.Lock()
		//c.pending[seq] = call
		//c.mutex.Unlock()
		go func(cr *batchRequest, call *Call) {
			select {
			case bytes, ok := <-cr.replyBytes:
				if ok {
					err := replyDecode(bytes, call.Reply, c.funcsCodecType)
					//c.mutex.Lock()
					//var ok bool
					//if call ,ok= c.pending[seq];ok{
					//	delete(c.pending, seq)
					//}
					//c.mutex.Unlock()
					if err != nil {
						if call != nil {
							call.Error = err
							call.done()
							return
						}
					}
					call.Done <- call
					return
				}
			case err, ok := <-cr.replyError:
				if ok {
					if err != nil {
						//c.mutex.Lock()
						//var ok bool
						//if call ,ok= c.pending[seq];ok{
						//	delete(c.pending, seq)
						//}
						//c.mutex.Unlock()
						if call != nil {
							call.Error = err
							call.done()
							return
						}
					}
				}
			}
		}(cr, call)
	}
}
