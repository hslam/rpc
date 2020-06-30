package rpc

import (
	"errors"
	"github.com/hslam/autowriter"
	"io"
	"sync"
)

type serverCodec struct {
	headerCodec    *Codec
	bodyCodec      *Codec
	replyBuffer    []byte
	responseBuffer []byte
	stream         Stream
	writer         io.WriteCloser
	mutex          sync.Mutex
	seq            uint64
	pending        map[uint64]uint64
}

// NewServerCodec returns a new rpc.ServerCodec using CODE-RPC on conn.
func NewServerCodec(conn io.ReadWriteCloser, codec *Codec) ServerCodec {
	var codecClone = codec.Clone()
	c := &serverCodec{
		headerCodec:    codecClone,
		bodyCodec:      codecClone,
		replyBuffer:    make([]byte, 1024),
		responseBuffer: make([]byte, 1024),
		pending:        make(map[uint64]uint64),
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

func (c *serverCodec) NumConcurrency() int {
	c.mutex.Lock()
	concurrency := len(c.pending)
	c.mutex.Unlock()
	return concurrency
}

func (c *serverCodec) ReadRequestHeader(ctx *Context) error {
	c.headerCodec.Request.Reset()
	var data []byte
	var err error
	data, err = c.stream.ReadMessage()
	if err != nil {
		return err
	}
	c.headerCodec.Codec.Unmarshal(data, c.headerCodec.Request)
	ctx.ServiceMethod = c.headerCodec.Request.GetServiceMethod()
	c.mutex.Lock()
	c.seq++
	ctx.Seq = c.seq
	c.pending[ctx.Seq] = c.headerCodec.Request.GetSeq()
	c.headerCodec.Request.SetSeq(0)
	c.mutex.Unlock()
	return nil
}

func (c *serverCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	return c.bodyCodec.Codec.Unmarshal(c.headerCodec.Request.GetArgs(), x)
}

func (c *serverCodec) WriteResponse(ctx *Context, x interface{}) error {
	c.mutex.Lock()
	reqSeq, ok := c.pending[ctx.Seq]
	if !ok {
		c.mutex.Unlock()
		return errors.New("invalid sequence number in response")
	}
	delete(c.pending, ctx.Seq)
	c.mutex.Unlock()

	var reply []byte
	var data []byte
	var err error
	if len(ctx.Error) == 0 {
		reply, err = c.bodyCodec.Codec.Marshal(c.replyBuffer, x)
		if err != nil {
			return err
		}
	}
	c.headerCodec.Response.SetSeq(reqSeq)
	c.headerCodec.Response.SetError(ctx.Error)
	c.headerCodec.Response.SetReply(reply)
	data, err = c.headerCodec.Codec.Marshal(c.responseBuffer, c.headerCodec.Response)
	if err != nil {
		return err
	}
	return c.stream.WriteMessage(data)
}

func (c *serverCodec) Close() error {
	if c.writer != nil {
		c.writer.Close()
	}
	return c.stream.Close()
}
