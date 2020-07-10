package rpc

import (
	"errors"
	"github.com/hslam/socket"
	"io"
	"sync"
)

var ErrShutdown = errors.New("connection is shut down")

type Call struct {
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
	codec    ClientCodec
	reqMutex sync.Mutex
	ctx      Context
	mutex    sync.Mutex
	seq      uint64
	pending  map[uint64]*Call
	closing  bool
	shutdown bool
}
type NewClientCodecFunc func(messages socket.Messages) ClientCodec

func NewClient() *Client {
	return &Client{
		pending: make(map[uint64]*Call),
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
	c := &Client{
		codec:   codec,
		pending: make(map[uint64]*Call),
	}
	go c.read()
	return c
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
	client.ctx.ServiceMethod = call.ServiceMethod
	err := client.codec.WriteRequest(&client.ctx, call.Args)
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
