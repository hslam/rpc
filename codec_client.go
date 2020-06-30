package rpc

import (
	"fmt"
	"github.com/hslam/autowriter"
	"io"
	"sync"
)

type clientCodec struct {
	headerCodec   *Codec
	bodyCodec     *Codec
	argsBuffer    []byte
	requestBuffer []byte
	stream        Stream
	writer        io.WriteCloser
	mutex         sync.Mutex
	pending       map[uint64]string
}

func NewClientCodec(conn io.ReadWriteCloser, codec *Codec) ClientCodec {
	var codecClone = codec.Clone()
	c := &clientCodec{
		headerCodec:   codecClone,
		bodyCodec:     codecClone,
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
	if fast {
		c.headerCodec = newFastCodec()
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
	args, err = c.bodyCodec.Codec.Marshal(c.argsBuffer, param)
	if err != nil {
		return err
	}
	c.headerCodec.Request.SetSeq(ctx.Seq)
	c.headerCodec.Request.SetServiceMethod(ctx.ServiceMethod)
	c.headerCodec.Request.SetArgs(args)
	data, err = c.headerCodec.Codec.Marshal(c.requestBuffer, c.headerCodec.Request)
	if err != nil {
		return err
	}
	return c.stream.WriteMessage(data)
}

func (c *clientCodec) ReadResponseHeader(ctx *Context) error {
	c.headerCodec.Response.Reset()
	var data []byte
	var err error
	data, err = c.stream.ReadMessage()
	if err != nil {
		return err
	}
	c.headerCodec.Codec.Unmarshal(data, c.headerCodec.Response)
	ctx.Error = ""
	ctx.Seq = c.headerCodec.Response.GetSeq()
	c.mutex.Lock()
	delete(c.pending, ctx.Seq)
	c.mutex.Unlock()
	if c.headerCodec.Response.GetError() != "" || len(c.headerCodec.Response.GetReply()) == 0 {
		if len(c.headerCodec.Response.GetError()) > 0 {
			return fmt.Errorf("invalid error %v", c.headerCodec.Response.GetError())
		}
		ctx.Error = c.headerCodec.Response.GetError()
	}
	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	if x == nil {
		return nil
	}
	return c.bodyCodec.Codec.Unmarshal(c.headerCodec.Response.GetReply(), x)
}

func (c *clientCodec) Close() error {
	if c.writer != nil {
		c.writer.Close()
	}
	return c.stream.Close()
}
