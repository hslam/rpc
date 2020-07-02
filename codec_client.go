package rpc

import (
	"fmt"
	"github.com/hslam/autowriter"
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
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
	stream        Stream
	writer        io.WriteCloser
	mutex         sync.Mutex
	pending       map[uint64]string
}

func NewClientCodec(conn io.ReadWriteCloser, bodyCodec codec.Codec, headerEncoder *encoder.Encoder) ClientCodec {
	var encoderClone *encoder.Encoder
	if headerEncoder != nil {
		encoderClone = headerEncoder.Clone()
		if bodyCodec == nil {
			bodyCodec = encoderClone.Codec
		}
	}
	if bodyCodec == nil {
		return nil
	}
	c := &clientCodec{
		headerEncoder: encoderClone,
		bodyCodec:     bodyCodec,
		argsBuffer:    make([]byte, 1024),
		requestBuffer: make([]byte, 1024),
		pending:       make(map[uint64]string),
	}
	c.writer = autowriter.NewAutoWriter(conn, false, 65536, 4, c)
	c.stream = &stream{
		Reader: conn,
		Writer: c.writer,
		Closer: conn,
		Send:   make([]byte, 1032),
		Read:   make([]byte, 1024),
	}
	if headerEncoder == nil {
		c.req = &request{}
		c.res = &response{}
	}
	return c
}

func (c *clientCodec) NumConcurrency() int {
	c.mutex.Lock()
	concurrency := len(c.pending)
	c.mutex.Unlock()
	return concurrency
}

func (c *clientCodec) WriteRequest(ctx *Context, param interface{}) error {
	c.mutex.Lock()
	c.pending[ctx.Seq] = ctx.ServiceMethod
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
	return c.stream.WriteMessage(data)
}

func (c *clientCodec) ReadResponseHeader(ctx *Context) error {
	var data []byte
	var err error
	data, err = c.stream.ReadMessage()
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
	return c.stream.Close()
}
