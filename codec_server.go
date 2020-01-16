package rpc

import (
	"fmt"
)

type serverCodec struct {
	msg        *msg
	batchCodec *batchCodec
	request    *request
	response   *response
	requests   []*request
	responses  []*response
}

func (c *serverCodec) Encode() ([]byte, error) {
	return c.msg.Encode()
}

func (c *serverCodec) Decode(b []byte) error {
	msg := &msg{}
	err := msg.Decode(b)
	c.msg = msg
	if err != nil {
		logger.Warnf("ServerCodec.Decode msg error: %s", err)
		return fmt.Errorf("ServerCodec.Decode msg error: %s", err)
	}
	if msg.version != Version {
		logger.Warnf("%.2f %.2f Version is not matched", Version, msg.version)
		return fmt.Errorf("%.2f %.2f Version is not matched", Version, msg.version)
	}
	if msg.msgType == MsgTypeHea {
		return nil
	}
	if msg.batch {
		batchCodec := &batchCodec{}
		batchCodec.Decode(msg.data)
		c.batchCodec = batchCodec
		length := len(batchCodec.data)
		c.requests = make([]*request, length)
		c.responses = make([]*response, length)
		for i, v := range batchCodec.data {
			req := &request{}
			err := req.Decode(v)
			if err != nil {
				c.requests[i] = nil
			}
			c.requests[i] = req
			c.responses[i] = &response{}
		}
		batchCodec.data = make([][]byte, length)
	} else {
		req := &request{}
		err := req.Decode(msg.data)
		if err != nil {
			logger.Warnf("ServerCodec.Decode id:%d req:%d error:%s ", msg.id, req.id, err)
			return fmt.Errorf("ServerCodec.Decode id:%d req:%d error:%s ", msg.id, req.id, err)
		}
		c.request = req
		c.response = &response{}
	}
	return nil
}
