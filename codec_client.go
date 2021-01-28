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
	req           *pbRequest
	res           *pbResponse
	buffer        []byte
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
		buffer:        make([]byte, bufferSize),
		argsBuffer:    make([]byte, bufferSize),
		requestBuffer: make([]byte, bufferSize),
	}
	c.messages = messages
	if batch, ok := c.messages.(socket.Batch); ok {
		batch.SetConcurrency(c.Concurrency)
	}
	if headerEncoder == nil {
		c.req = &pbRequest{}
		c.res = &pbResponse{}
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
	} else {
		c.req.SetSeq(ctx.Seq)
		c.req.SetUpgrade(ctx.Upgrade)
		c.req.SetServiceMethod(ctx.ServiceMethod)
		c.req.SetArgs(args)
		size := c.req.Size()
		var buf = checkBuffer(c.requestBuffer, size)
		var n int
		n, err = c.req.MarshalTo(buf)
		data = buf[:n]
	}
	if err == nil {
		atomic.AddInt64(&c.count, 1)
		if atomic.LoadUint32(&c.closed) > 0 {
			return io.EOF
		}
		err = c.messages.WriteMessage(data)
	}
	return err
}

func (c *clientCodec) ReadResponseHeader(ctx *Context) error {
	if atomic.LoadUint32(&c.closed) > 0 {
		return io.EOF
	}
	var data []byte
	var err error
	data, err = c.messages.ReadMessage(c.buffer)
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
		if err == nil {
			ctx.Seq = c.headerEncoder.Response.GetSeq()
			ctx.Error = c.headerEncoder.Response.GetError()
			ctx.value = c.headerEncoder.Response.GetReply()
		}
	} else {
		c.res.Reset()
		err = c.res.Unmarshal(data)
		if err == nil {
			ctx.Seq = c.res.GetSeq()
			ctx.Error = c.res.GetError()
			ctx.value = c.res.GetReply()
		}
	}
	return err
}

func (c *clientCodec) ReadResponseBody(reply []byte, x interface{}) error {
	if x == nil {
		return errors.New("x is nil")
	}
	if c.headerEncoder != nil {
		return c.bodyCodec.Unmarshal(reply, x)
	}
	return c.bodyCodec.Unmarshal(reply, x)
}

func (c *clientCodec) Close() error {
	atomic.StoreUint32(&c.closed, 1)
	return c.messages.Close()
}
