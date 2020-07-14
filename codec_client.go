package rpc

import (
	"fmt"
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/socket"
	"io"
	"sync"
)

type clientCodec struct {
	headerEncoder *encoder.Encoder
	bodyCodec     codec.Codec
	req           *request
	res           *response
	argsBuffer    []byte
	requestBuffer []byte
	messages      socket.Messages
	writer        io.WriteCloser
	mutex         sync.Mutex
	pending       map[uint64]bool
}

func NewClientCodec(bodyCodec codec.Codec, headerEncoder *encoder.Encoder, messages socket.Messages) ClientCodec {
	if messages == nil {
		return nil
	}
	if headerEncoder != nil {
		if bodyCodec == nil {
			bodyCodec = headerEncoder.Codec
		}
	}
	if bodyCodec == nil {
		return nil
	}
	c := &clientCodec{
		headerEncoder: headerEncoder,
		bodyCodec:     bodyCodec,
		argsBuffer:    make([]byte, 1024),
		requestBuffer: make([]byte, 1024),
		pending:       make(map[uint64]bool),
	}
	c.messages = messages
	c.messages.SetBatch(c)
	if headerEncoder == nil {
		c.req = &request{}
		c.res = &response{}
	}
	return c
}
func (c *clientCodec) Concurrency() int {
	c.mutex.Lock()
	concurrency := len(c.pending)
	c.mutex.Unlock()
	return concurrency
}

func (c *clientCodec) WriteRequest(ctx *Context, param interface{}) error {
	c.mutex.Lock()
	c.pending[ctx.Seq] = false
	c.mutex.Unlock()
	var args []byte
	var data []byte
	var err error
	args, err = c.bodyCodec.Marshal(c.argsBuffer, param)
	if err != nil {
		return err
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Request.SetSeq(ctx.Seq)
		c.headerEncoder.Request.SetServiceMethod(ctx.ServiceMethod)
		c.headerEncoder.Request.SetArgs(args)
		data, err = c.headerEncoder.Codec.Marshal(c.requestBuffer, c.headerEncoder.Request)
		if err != nil {
			return err
		}
	} else {
		c.req.SetSeq(ctx.Seq)
		c.req.SetServiceMethod(ctx.ServiceMethod)
		c.req.SetArgs(args)
		data, err = c.req.Marshal(c.requestBuffer)
		if err != nil {
			return err
		}
	}
	return c.messages.WriteMessage(data)
}

func (c *clientCodec) ReadResponseHeader(ctx *Context) error {
	var data []byte
	var err error
	data, err = c.messages.ReadMessage()
	if err != nil {
		return err
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Response.Reset()
		c.headerEncoder.Codec.Unmarshal(data, c.headerEncoder.Response)
		ctx.Error = ""
		ctx.Seq = c.headerEncoder.Response.GetSeq()
		c.mutex.Lock()
		delete(c.pending, ctx.Seq)
		c.mutex.Unlock()
		if c.headerEncoder.Response.GetError() != "" || len(c.headerEncoder.Response.GetReply()) == 0 {
			if len(c.headerEncoder.Response.GetError()) > 0 {
				return fmt.Errorf("invalid error %v", c.headerEncoder.Response.GetError())
			}
			ctx.Error = c.headerEncoder.Response.GetError()
		}
	} else {
		c.res.Reset()
		c.res.Unmarshal(data)
		ctx.Error = ""
		ctx.Seq = c.res.GetSeq()
		c.mutex.Lock()
		delete(c.pending, ctx.Seq)
		c.mutex.Unlock()
		if c.res.GetError() != "" || len(c.res.GetReply()) == 0 {
			if len(c.res.GetError()) > 0 {
				return fmt.Errorf("invalid error %v", c.res.GetError())
			}
			ctx.Error = c.res.GetError()
		}
	}

	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	if x == nil {
		return nil
	}
	if c.headerEncoder != nil {
		return c.bodyCodec.Unmarshal(c.headerEncoder.Response.GetReply(), x)
	}
	return c.bodyCodec.Unmarshal(c.res.GetReply(), x)
}

func (c *clientCodec) Close() error {
	if c.writer != nil {
		c.writer.Close()
	}
	return c.messages.Close()
}
