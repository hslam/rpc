// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"github.com/hslam/socket"
	"io"
	"sync/atomic"
)

type clientCodec struct {
	headerEncoder *Encoder
	bodyCodec     Codec
	req           *request
	res           *response
	argsBuffer    []byte
	requestBuffer []byte
	messages      socket.Messages
	count         int64
	closed        uint32
}

// NewClientCodec returns a new ClientCodec.
func NewClientCodec(bodyCodec Codec, headerEncoder *Encoder, messages socket.Messages) ClientCodec {
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
		argsBuffer:    make([]byte, bufferSize),
		requestBuffer: make([]byte, bufferSize),
	}
	c.messages = messages
	if batch, ok := c.messages.(socket.Batch); ok {
		batch.SetConcurrency(c.Concurrency)
	}
	if headerEncoder == nil {
		c.req = &request{}
		c.res = &response{}
	}
	return c
}

func (c *clientCodec) Concurrency() int {
	return int(atomic.LoadInt64(&c.count))
}

func (c *clientCodec) WriteRequest(ctx *Context, param interface{}) error {
	var args []byte
	var data []byte
	var err error
	if ctx.upgrade.NoRequest != noRequest {
		args, err = c.bodyCodec.Marshal(c.argsBuffer, param)
		if err != nil {
			return err
		}
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Request.SetSeq(ctx.Seq)
		c.headerEncoder.Request.SetUpgrade(ctx.Upgrade)
		c.headerEncoder.Request.SetServiceMethod(ctx.ServiceMethod)
		c.headerEncoder.Request.SetArgs(args)
		data, err = c.headerEncoder.Codec.Marshal(c.requestBuffer, c.headerEncoder.Request)
		if err != nil {
			return err
		}
	} else {
		c.req.SetSeq(ctx.Seq)
		c.req.SetUpgrade(ctx.Upgrade)
		c.req.SetServiceMethod(ctx.ServiceMethod)
		c.req.SetArgs(args)
		data, err = c.req.Marshal(c.requestBuffer)
		if err != nil {
			return err
		}
	}
	atomic.AddInt64(&c.count, 1)
	if atomic.LoadUint32(&c.closed) > 0 {
		return io.EOF
	}
	return c.messages.WriteMessage(data)
}

func (c *clientCodec) ReadResponseHeader(ctx *Context) error {
	if atomic.LoadUint32(&c.closed) > 0 {
		return io.EOF
	}
	var data []byte
	var err error
	data, err = c.messages.ReadMessage(nil)
	if err != nil {
		if err == io.EOF {
			atomic.StoreUint32(&c.closed, 1)
		}
		return err
	}
	defer atomic.AddInt64(&c.count, -1)
	if c.headerEncoder != nil {
		c.headerEncoder.Response.Reset()
		err = c.headerEncoder.Codec.Unmarshal(data, c.headerEncoder.Response)
		if err != nil {
			return err
		}
		ctx.Error = ""
		ctx.Seq = c.headerEncoder.Response.GetSeq()
		if c.headerEncoder.Response.GetError() != "" || len(c.headerEncoder.Response.GetReply()) == 0 {
			ctx.Error = c.headerEncoder.Response.GetError()
		} else if len(c.headerEncoder.Response.GetReply()) > 0 {
			ctx.value = c.headerEncoder.Response.GetReply()
		}
	} else {
		c.res.Reset()
		_, err = c.res.Unmarshal(data)
		if err != nil {
			return err
		}
		ctx.Error = ""
		ctx.Seq = c.res.GetSeq()
		if c.res.GetError() != "" || len(c.res.GetReply()) == 0 {
			ctx.Error = c.res.GetError()
		} else if len(c.res.GetReply()) > 0 {
			ctx.value = c.res.GetReply()
		}
	}
	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	if x == nil {
		return errors.New("x is nil")
	}
	if c.headerEncoder != nil {
		return c.bodyCodec.Unmarshal(c.headerEncoder.Response.GetReply(), x)
	}
	return c.bodyCodec.Unmarshal(c.res.GetReply(), x)
}

func (c *clientCodec) Close() error {
	atomic.StoreUint32(&c.closed, 1)
	return c.messages.Close()
}
