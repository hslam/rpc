// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"fmt"
	"github.com/hslam/codec"
	"github.com/hslam/socket"
	"sync/atomic"
)

type clientCodec struct {
	headerEncoder *Encoder
	bodyCodec     codec.Codec
	req           *request
	res           *response
	argsBuffer    []byte
	requestBuffer []byte
	messages      socket.Messages
	count         int64
}

// NewClientCodec returns a new ClientCodec.
func NewClientCodec(bodyCodec codec.Codec, headerEncoder *Encoder, messages socket.Messages) ClientCodec {
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
	}
	c.messages = messages
	c.messages.SetConcurrency(c.Concurrency)
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
	return c.messages.WriteMessage(data)
}

func (c *clientCodec) ReadResponseHeader(ctx *Context) error {
	var data []byte
	var err error
	data, err = c.messages.ReadMessage()
	if err != nil {
		return err
	}
	defer atomic.AddInt64(&c.count, -1)
	if c.headerEncoder != nil {
		c.headerEncoder.Response.Reset()
		c.headerEncoder.Codec.Unmarshal(data, c.headerEncoder.Response)
		ctx.Error = ""
		ctx.Seq = c.headerEncoder.Response.GetSeq()
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
		if c.res.GetError() != "" || len(c.res.GetReply()) == 0 {
			if len(c.res.GetError()) > 0 {
				return fmt.Errorf("invalid error %v", c.res.GetError())
			}
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
	return c.messages.Close()
}
