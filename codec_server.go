// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/socket"
	"sync/atomic"
)

type serverCodec struct {
	headerEncoder  *Encoder
	bodyCodec      codec.Codec
	req            *request
	res            *response
	replyBuffer    []byte
	responseBuffer []byte
	messages       socket.Messages
	count          int64
}

// NewServerCodec returns a new ServerCodec.
func NewServerCodec(bodyCodec codec.Codec, headerEncoder *Encoder, messages socket.Messages) ServerCodec {
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
		replyBuffer:    make([]byte, 1024),
		responseBuffer: make([]byte, 1024),
	}
	c.messages = messages
	c.messages.SetConcurrency(c.Concurrency)
	if headerEncoder == nil {
		c.req = &request{}
		c.res = &response{}
	}
	return c
}

func (c *serverCodec) Concurrency() int {
	return int(atomic.LoadInt64(&c.count))
}

func (c *serverCodec) ReadRequestHeader(ctx *Context) error {
	var data []byte
	var err error
	data, err = c.messages.ReadMessage()
	if err != nil {
		return err
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Request.Reset()
		c.headerEncoder.Codec.Unmarshal(data, c.headerEncoder.Request)
		ctx.ServiceMethod = c.headerEncoder.Request.GetServiceMethod()
		ctx.Upgrade = c.headerEncoder.Request.GetUpgrade()
		ctx.Seq = c.headerEncoder.Request.GetSeq()
	} else {
		c.req.Reset()
		c.req.Unmarshal(data)
		ctx.ServiceMethod = c.req.GetServiceMethod()
		ctx.Upgrade = c.req.GetUpgrade()
		ctx.Seq = c.req.GetSeq()
	}
	atomic.AddInt64(&c.count, 1)
	return nil
}

func (c *serverCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return errors.New("x is nil")
	}
	if c.headerEncoder != nil {
		return c.bodyCodec.Unmarshal(c.headerEncoder.Request.GetArgs(), x)
	}
	return c.bodyCodec.Unmarshal(c.req.GetArgs(), x)
}

func (c *serverCodec) WriteResponse(ctx *Context, x interface{}) error {
	defer atomic.AddInt64(&c.count, -1)
	reqSeq := ctx.Seq
	var reply []byte
	var data []byte
	var err error
	if len(ctx.Error) == 0 && ctx.upgrade.NoResponse != noResponse {
		reply, err = c.bodyCodec.Marshal(c.replyBuffer, x)
		if err != nil {
			return err
		}
	} else if len(ctx.value) > 0 {
		reply = ctx.value
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Response.SetSeq(reqSeq)
		c.headerEncoder.Response.SetError(ctx.Error)
		c.headerEncoder.Response.SetReply(reply)
		data, err = c.headerEncoder.Codec.Marshal(c.responseBuffer, c.headerEncoder.Response)
		if err != nil {
			return err
		}
	} else {
		c.res.SetSeq(reqSeq)
		c.res.SetError(ctx.Error)
		c.res.SetReply(reply)
		data, err = c.res.Marshal(c.responseBuffer)
		if err != nil {
			return err
		}
	}
	return c.messages.WriteMessage(data)
}

func (c *serverCodec) Close() error {
	return c.messages.Close()
}
