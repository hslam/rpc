// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"github.com/hslam/socket"
	"io"
	"sync/atomic"
)

type serverCodec struct {
	headerEncoder  *Encoder
	bodyCodec      Codec
	req            *pbRequest
	res            *pbResponse
	replyBuffer    []byte
	responseBuffer []byte
	messages       socket.Messages
	count          int64
	closed         uint32
}

// NewServerCodec returns a new ServerCodec.
func NewServerCodec(bodyCodec Codec, headerEncoder *Encoder, messages socket.Messages, noBatch bool) ServerCodec {
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
	c := &serverCodec{
		headerEncoder:  headerEncoder,
		bodyCodec:      bodyCodec,
		replyBuffer:    make([]byte, bufferSize),
		responseBuffer: make([]byte, bufferSize),
	}
	c.messages = messages
	if !noBatch {
		if batch, ok := c.messages.(socket.Batch); ok {
			batch.SetConcurrency(c.Concurrency)
		}
	}
	if headerEncoder == nil {
		c.req = &pbRequest{}
		c.res = &pbResponse{}
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
	var buffer []byte
	if ctx != nil && len(ctx.buffer) > 0 {
		buffer = ctx.buffer
	}
	data, err = c.messages.ReadMessage(buffer)
	if err != nil {
		if err == io.EOF {
			atomic.StoreUint32(&c.closed, 1)
		}
		return err
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Request.Reset()
		err = c.headerEncoder.Codec.Unmarshal(data, c.headerEncoder.Request)
		if err == nil {
			ctx.ServiceMethod = c.headerEncoder.Request.GetServiceMethod()
			ctx.Upgrade = c.headerEncoder.Request.GetUpgrade()
			ctx.Seq = c.headerEncoder.Request.GetSeq()
			ctx.value = c.headerEncoder.Request.GetArgs()
		}
	} else {
		c.req.Reset()
		err = c.req.Unmarshal(data)
		if err == nil {
			ctx.ServiceMethod = c.req.GetServiceMethod()
			ctx.Upgrade = c.req.GetUpgrade()
			ctx.Seq = c.req.GetSeq()
			ctx.value = c.req.GetArgs()
		}
	}
	atomic.AddInt64(&c.count, 1)
	return err
}

func (c *serverCodec) ReadRequestBody(args []byte, x interface{}) error {
	if x == nil {
		return errors.New("x is nil")
	}
	if c.headerEncoder != nil {
		return c.bodyCodec.Unmarshal(args, x)
	}
	return c.bodyCodec.Unmarshal(args, x)
}

func (c *serverCodec) WriteResponse(ctx *Context, x interface{}) error {
	reqSeq := ctx.Seq
	var reply []byte
	var data []byte
	var err error
	if len(ctx.Error) == 0 && ctx.upgrade.NoResponse != noResponse {
		reply, err = c.bodyCodec.Marshal(c.replyBuffer, x)
		if err != nil {
			ctx.Error = err.Error()
		}
	} else if len(ctx.value) > 0 {
		reply = ctx.value
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Response.SetSeq(reqSeq)
		c.headerEncoder.Response.SetError(ctx.Error)
		c.headerEncoder.Response.SetReply(reply)
		data, err = c.headerEncoder.Codec.Marshal(c.responseBuffer, c.headerEncoder.Response)
	} else {
		c.res.SetSeq(reqSeq)
		c.res.SetError(ctx.Error)
		c.res.SetReply(reply)
		size := c.res.Size()
		var buf = checkBuffer(c.responseBuffer, size)
		var n int
		n, err = c.res.MarshalTo(buf)
		data = buf[:n]
	}
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
