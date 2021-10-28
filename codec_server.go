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

type serverCodec struct {
	headerEncoder Encoder
	bodyCodec     Codec
	pool          *buffer.Pool
	messages      socket.Messages
	count         int64
	closed        uint32
}

// NewServerCodec returns a new ServerCodec.
func NewServerCodec(bodyCodec Codec, headerEncoder Encoder, messages socket.Messages, noBatch bool) ServerCodec {
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
	c := &serverCodec{
		headerEncoder: headerEncoder,
		bodyCodec:     bodyCodec,
		pool:          buffer.AssignPool(bufferSize),
	}
	c.messages = messages
	if !noBatch {
		if batch, ok := c.messages.(socket.Batch); ok {
			batch.SetConcurrency(c.Concurrency)
		}
	}
	return c
}

func (c *serverCodec) Concurrency() int {
	return int(atomic.LoadInt64(&c.count))
}

func (c *serverCodec) ReadRequestHeader(ctx *Context) error {
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
		req := c.headerEncoder.NewRequest()
		req.Reset()
		codec := c.headerEncoder.NewCodec()
		err = codec.Unmarshal(data, req)
		if err == nil {
			ctx.ServiceMethod = req.GetServiceMethod()
			ctx.Upgrade = req.GetUpgrade()
			ctx.Seq = req.GetSeq()
			ctx.value = req.GetArgs()
		}
	} else {
		req := &pbRequest{}
		req.Reset()
		err = req.Unmarshal(data)
		if err == nil {
			ctx.ServiceMethod = req.GetServiceMethod()
			ctx.Upgrade = req.GetUpgrade()
			ctx.Seq = req.GetSeq()
			ctx.value = req.GetArgs()
		}
	}
	atomic.AddInt64(&c.count, 1)
	return err
}

func (c *serverCodec) ReadRequestBody(args []byte, x interface{}) error {
	if x == nil {
		return errors.New("x is nil")
	}
	return c.bodyCodec.Unmarshal(args, x)
}

func (c *serverCodec) WriteResponse(ctx *Context, x interface{}) error {
	reqSeq := ctx.Seq
	var reply []byte
	var data []byte
	var err error
	hasResponse := len(ctx.Error) == 0 && ctx.upgrade.NoResponse != noResponse
	var replyBuffer []byte
	if hasResponse {
		replyBuffer = c.pool.GetBuffer(bufferSize)
		reply, err = c.bodyCodec.Marshal(replyBuffer, x)
		if err != nil {
			ctx.Error = err.Error()
		}
	} else if len(ctx.value) > 0 {
		reply = ctx.value
	}
	var responseBuffer = c.pool.GetBuffer(bufferSize)
	if c.headerEncoder != nil {
		res := c.headerEncoder.NewResponse()
		res.SetSeq(reqSeq)
		res.SetError(ctx.Error)
		res.SetReply(reply)
		codec := c.headerEncoder.NewCodec()
		data, err = codec.Marshal(responseBuffer, res)
	} else {
		var res = &pbResponse{}
		res.SetSeq(reqSeq)
		res.SetError(ctx.Error)
		res.SetReply(reply)
		size := res.Size()
		var buf = checkBuffer(responseBuffer, size)
		var n int
		n, err = res.MarshalTo(buf)
		data = buf[:n]
	}
	defer func() {
		if hasResponse {
			c.pool.PutBuffer(replyBuffer)
		}
		c.pool.PutBuffer(responseBuffer)
	}()
	if err == nil {
		if atomic.LoadUint32(&c.closed) > 0 {
			atomic.AddInt64(&c.count, -1)
			return io.EOF
		}
		err = c.messages.WriteMessage(data)
		atomic.AddInt64(&c.count, -1)
	}
	return err

}

func (c *serverCodec) Close() error {
	atomic.StoreUint32(&c.closed, 1)
	return c.messages.Close()
}
