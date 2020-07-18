package rpc

import (
	"errors"
	"github.com/hslam/socket"
	"io"
	"sync"
)

var ErrShutdown = errors.New("connection is shut down")

type Call struct {
	heartbeat     bool
	noRequest     bool
	noResponse    bool
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

func (call *Call) done() {
	select {
	case call.Done <- call:
	default:
	}
}

type Client struct {
	codec         ClientCodec
	reqMutex      sync.Mutex
	ctx           Context
	mutex         sync.Mutex
	seq           uint64
	pending       map[uint64]*Call
	upgradePool   *sync.Pool
	upgradeBuffer []byte
	closing       bool
	shutdown      bool
}

type NewClientCodecFunc func(messages socket.Messages) ClientCodec

func NewClient() *Client {
	return &Client{
		pending:       make(map[uint64]*Call),
		upgradePool:   &sync.Pool{New: func() interface{} { return &upgrade{} }},
		upgradeBuffer: make([]byte, 1024),
	}
}

func (client *Client) Dial(s socket.Socket, address string, New NewClientCodecFunc) (*Client, error) {
	conn, err := s.Dial(address)
	if err != nil {
		return nil, err
	}
	client.codec = New(conn.Messages())
	go client.read()
	return client, nil
}

func NewClientWithCodec(codec ClientCodec) *Client {
	c := NewClient()
	c.codec = codec
	go c.read()
	return c
}

func (client *Client) getUpgrade() *upgrade {
	return client.upgradePool.Get().(*upgrade)
}

func (client *Client) putUpgrade(u *upgrade) {
	u.Reset()
	client.upgradePool.Put(u)
}

func (client *Client) write(call *Call) {
	client.reqMutex.Lock()
	defer client.reqMutex.Unlock()
	client.mutex.Lock()
	if client.shutdown || client.closing {
		client.mutex.Unlock()
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := client.seq
	client.seq++
	client.pending[seq] = call
	client.mutex.Unlock()
	client.ctx.Seq = seq
	client.ctx.heartbeat = call.heartbeat
	client.ctx.noRequest = call.noRequest
	client.ctx.noResponse = call.noResponse
	u := client.getUpgrade()
	if call.heartbeat || call.noRequest || call.noResponse {
		if call.heartbeat {
			u.Heartbeat = Heartbeat
		}
		if call.noRequest {
			u.NoRequest = NoRequest
		}
		if call.noResponse {
			u.NoResponse = NoResponse
		}
		client.ctx.Upgrade, _ = u.Marshal(client.upgradeBuffer)
	}
	client.ctx.ServiceMethod = call.ServiceMethod
	err := client.codec.WriteRequest(&client.ctx, call.Args)
	client.putUpgrade(u)
	if err != nil {
		client.mutex.Lock()
		delete(client.pending, seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (client *Client) read() {
	var err error
	var ctx Context
	for err == nil {
		ctx = Context{}
		err = client.codec.ReadResponseHeader(&ctx)
		if err != nil {
			break
		}
		seq := ctx.Seq
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		switch {
		case call == nil:
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		case ctx.Error != "":
			call.Error = errors.New(ctx.Error)
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
			call.done()
		default:
			if len(ctx.Upgrade) > 0 {
				u := client.getUpgrade()
				u.Unmarshal(ctx.Upgrade)
				if u.Heartbeat == Heartbeat || u.NoResponse == NoResponse {
					client.putUpgrade(u)
					call.done()
					continue
				}
				client.putUpgrade(u)
			}
			err = client.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	client.reqMutex.Lock()
	client.mutex.Lock()
	client.shutdown = true
	closing := client.closing
	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
	client.mutex.Unlock()
	client.reqMutex.Unlock()
	if err != io.EOF && !closing {
		logger.Allln("rpc: client protocol error:", err)
	}
}

func (client *Client) NumCalls() (n uint64) {
	client.mutex.Lock()
	n = uint64(len(client.pending))
	client.mutex.Unlock()
	return
}

func (client *Client) Close() error {
	client.mutex.Lock()
	if client.closing {
		client.mutex.Unlock()
		return ErrShutdown
	}
	client.closing = true
	client.mutex.Unlock()
	return client.codec.Close()
}

func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *Call, 10)
	} else {
		if cap(done) == 0 {
			logger.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	client.write(call)
	return call
}

func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

func (client *Client) Ping() error {
	call := new(Call)
	call.heartbeat = true
	call.noRequest = true
	call.noResponse = true
	call.Done = make(chan *Call, 10)
	client.write(call)
	c := <-call.Done
	return c.Error
}
