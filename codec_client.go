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
	writeBufSize  int
	bufferPool    *buffer.Pool
	messages      socket.Messages
	closed        uint32
}

// NewClientCodec returns a new ClientCodec.
func NewClientCodec(bodyCodec Codec, headerEncoder Encoder, messages socket.Messages, writeBufSize int) ClientCodec {
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
	if writeBufSize < 1 {
		writeBufSize = bufferSize
	}
	c := &clientCodec{
		headerEncoder: headerEncoder,
		bodyCodec:     bodyCodec,
		writeBufSize:  writeBufSize,
		bufferPool:    buffer.AssignPool(writeBufSize),
	}
	c.messages = messages
	if set, ok := c.messages.(socket.BufferedOutput); ok {
		set.SetBufferedOutput(writeBufSize)
	}
	return c
}

//SetBufferSize sets buffer size.
func (c *clientCodec) SetBufferSize(size int) {
	if size < 1 {
		size = bufferSize
	}
	c.writeBufSize = size
	c.bufferPool = buffer.AssignPool(size)
	if set, ok := c.messages.(socket.BufferedOutput); ok {
		set.SetBufferedOutput(size)
	}
}

// SetDirectIO disables async io.
func (c *clientCodec) SetDirectIO(directIO bool) {
	if directIO {
		if set, ok := c.messages.(socket.BufferedOutput); ok {
			set.SetBufferedOutput(c.writeBufSize)
		}
	} else {
		if set, ok := c.messages.(socket.BufferedOutput); ok {
			set.SetBufferedOutput(0)
		}
	}
}

func (c *clientCodec) Messages() socket.Messages {
	return c.messages
}

func (c *clientCodec) WriteRequest(ctx *Context, param interface{}) error {
	if atomic.LoadUint32(&c.closed) > 0 {
		return io.EOF
	}
	var args []byte
	var data []byte
	var err error
	var argsBuffer []byte
	if ctx.upgrade.NoRequest != noRequest {
		argsBuffer = c.bufferPool.GetBuffer(0)
		args, err = c.bodyCodec.Marshal(argsBuffer, param)
		if err != nil {
			return err
		}
	}
	var requestBuffer = c.bufferPool.GetBuffer(0)
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
	if err == nil {
		err = c.messages.WriteMessage(data)
	}
	if ctx.upgrade.NoRequest != noRequest {
		c.bufferPool.PutBuffer(argsBuffer)
	}
	c.bufferPool.PutBuffer(requestBuffer)
	return err
}

func (c *clientCodec) ReadResponseHeader(ctx *Context) error {
	if atomic.LoadUint32(&c.closed) > 0 {
		return io.EOF
	}
	var data = ctx.data
	var err error
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
