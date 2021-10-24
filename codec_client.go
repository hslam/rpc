// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"github.com/hslam/buffer"
	"github.com/hslam/socket"
	"io"
	"sync/atomic"
)

type clientCodec struct {
	headerEncoder Encoder
	bodyCodec     Codec
	pool          *buffer.Pool
	messages      socket.Messages
	count         int64
	closed        uint32
}

// NewClientCodec returns a new ClientCodec.
func NewClientCodec(bodyCodec Codec, headerEncoder Encoder, messages socket.Messages) ClientCodec {
	if messages == nil {
		return nil
	}
	if headerEncoder != nil {
		if bodyCodec == nil {
			bodyCodec = headerEncoder.NewCodec()
		}
	}
	if bodyCodec == nil {
		return nil
	}
	c := &clientCodec{
		headerEncoder: headerEncoder,
		bodyCodec:     bodyCodec,
		pool:          buffers.AssignPool(bufferSize),
	}
	c.messages = messages
	if batch, ok := c.messages.(socket.Batch); ok {
		batch.SetConcurrency(c.Concurrency)
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
	var argsBuffer []byte
	if ctx.upgrade.NoRequest != noRequest {
		argsBuffer = c.pool.GetBuffer(bufferSize)
		args, err = c.bodyCodec.Marshal(argsBuffer, param)
		if err != nil {
			return err
		}
	}
	var requestBuffer = c.pool.GetBuffer(bufferSize)
	if c.headerEncoder != nil {
		req := c.headerEncoder.NewRequest()
		req.SetSeq(ctx.Seq)
		req.SetUpgrade(ctx.Upgrade)
		req.SetServiceMethod(ctx.ServiceMethod)
		req.SetArgs(args)
		codec := c.headerEncoder.NewCodec()
		data, err = codec.Marshal(requestBuffer, req)
	} else {
		req := &pbRequest{}
		req.SetSeq(ctx.Seq)
		req.SetUpgrade(ctx.Upgrade)
		req.SetServiceMethod(ctx.ServiceMethod)
		req.SetArgs(args)
		size := req.Size()
		var buf = checkBuffer(requestBuffer, size)
		var n int
		n, err = req.MarshalTo(buf)
		data = buf[:n]
	}
	defer func() {
		if ctx.upgrade.NoRequest != noRequest {
			c.pool.PutBuffer(argsBuffer)
		}
		c.pool.PutBuffer(requestBuffer)
	}()
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
	var buffer = ctx.buffer
	data, err = c.messages.ReadMessage(buffer)
	if err != nil {
		if err == io.EOF {
			atomic.StoreUint32(&c.closed, 1)
		}
		return err
	}
	if c.headerEncoder != nil {
		res := c.headerEncoder.NewResponse()
		res.Reset()
		codec := c.headerEncoder.NewCodec()
		err = codec.Unmarshal(data, res)
		if err == nil {
			ctx.Seq = res.GetSeq()
			ctx.Error = res.GetError()
			ctx.value = res.GetReply()
		}
	} else {
		res := &pbResponse{}
		res.Reset()
		err = res.Unmarshal(data)
		if err == nil {
			ctx.Seq = res.GetSeq()
			ctx.Error = res.GetError()
			ctx.value = res.GetReply()
		}
	}
	atomic.AddInt64(&c.count, -1)
	return err
}

func (c *clientCodec) ReadResponseBody(reply []byte, x interface{}) error {
	if x == nil {
		return errors.New("x is nil")
	}
	return c.bodyCodec.Unmarshal(reply, x)
}

func (c *clientCodec) Close() error {
	atomic.StoreUint32(&c.closed, 1)
	return c.messages.Close()
}
