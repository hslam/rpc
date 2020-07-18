package rpc

import (
	"errors"
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/socket"
	"io"
	"sync"
)

type serverCodec struct {
	headerEncoder  *encoder.Encoder
	bodyCodec      codec.Codec
	req            *request
	res            *response
	replyBuffer    []byte
	responseBuffer []byte
	messages       socket.Messages
	writer         io.WriteCloser
	mutex          sync.Mutex
	seq            uint64
	pending        map[uint64]uint64
}

func NewServerCodec(bodyCodec codec.Codec, headerEncoder *encoder.Encoder, messages socket.Messages) ServerCodec {
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
		pending:        make(map[uint64]uint64),
	}
	c.messages = messages
	c.messages.SetBatch(c)
	if headerEncoder == nil {
		c.req = &request{}
		c.res = &response{}
	}
	return c
}

func (c *serverCodec) Concurrency() int {
	c.mutex.Lock()
	concurrency := len(c.pending)
	c.mutex.Unlock()
	return concurrency
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
		c.mutex.Lock()
		c.seq++
		ctx.Seq = c.seq
		c.pending[ctx.Seq] = c.headerEncoder.Request.GetSeq()
		c.mutex.Unlock()
	} else {
		c.req.Reset()
		c.req.Unmarshal(data)
		ctx.ServiceMethod = c.req.GetServiceMethod()
		ctx.Upgrade = c.req.GetUpgrade()
		c.mutex.Lock()
		c.seq++
		ctx.Seq = c.seq
		c.pending[ctx.Seq] = c.req.GetSeq()
		c.mutex.Unlock()
	}
	return nil
}

func (c *serverCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	if c.headerEncoder != nil {
		return c.bodyCodec.Unmarshal(c.headerEncoder.Request.GetArgs(), x)
	}
	return c.bodyCodec.Unmarshal(c.req.GetArgs(), x)
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
	if len(ctx.Error) == 0 && !ctx.noResponse {
		reply, err = c.bodyCodec.Marshal(c.replyBuffer, x)
		if err != nil {
			return err
		}
	}
	if c.headerEncoder != nil {
		c.headerEncoder.Response.SetSeq(reqSeq)
		c.headerEncoder.Response.SetUpgrade(ctx.Upgrade)
		c.headerEncoder.Response.SetError(ctx.Error)
		c.headerEncoder.Response.SetReply(reply)
		data, err = c.headerEncoder.Codec.Marshal(c.responseBuffer, c.headerEncoder.Response)
		if err != nil {
			return err
		}
	} else {
		c.res.SetSeq(reqSeq)
		c.res.SetUpgrade(ctx.Upgrade)
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
	if c.writer != nil {
		c.writer.Close()
	}
	return c.messages.Close()
}
