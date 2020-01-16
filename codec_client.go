package rpc

import (
	"errors"
	"github.com/hslam/rpc/pb"
)

type clientCodec struct {
	clientID       uint64
	reqID          uint64
	name           string
	args           interface{}
	funcsCodecType CodecType
	noRequest      bool
	noResponse     bool
	reply          interface{}
	res            *response
	batch          bool
	batchingAsync  bool
	compressType   CompressType
	compressLevel  CompressLevel
	msg            *msg
	batchCodec     *batchCodec
	requests       []*batchRequest
	responses      []*response
}

func (c *clientCodec) Encode() ([]byte, error) {
	var err error
	c.msg = &msg{}
	c.msg.version = Version
	c.msg.id = c.clientID
	c.msg.msgType = MsgTypeReq
	c.msg.batch = c.batch
	c.msg.codecType = c.funcsCodecType
	c.msg.compressType = c.compressType
	c.msg.compressLevel = c.compressLevel
	if c.batch == false {
		var req *request
		if c.noRequest == false {
			argsBytes, err := argsEncode(c.args, c.funcsCodecType)
			if err != nil {
				logger.Errorln("ArgsEncode error: ", err)
				return nil, err
			}
			req = &request{c.reqID, c.name, c.noRequest, c.noResponse, argsBytes}
		} else {
			req = &request{c.reqID, c.name, c.noRequest, c.noResponse, nil}
		}
		c.msg.data, err = req.Encode()
		if err != nil {
			logger.Errorln("RequestEncode error: ", err)
			return nil, err
		}
	} else {
		reqBytesArr := make([][]byte, len(c.requests))
		c.responses = make([]*response, len(c.requests))
		for i, v := range c.requests {
			req := &request{v.id, v.name, v.noRequest, v.noResponse, v.argsBytes}
			reqBytes, _ := req.Encode()
			reqBytesArr[i] = reqBytes
			c.responses[i] = &response{}
		}
		batchCodec := &batchCodec{async: c.batchingAsync, data: reqBytesArr}
		c.batchCodec = batchCodec
		c.msg.data, _ = batchCodec.Encode()
	}
	return c.msg.Encode()
}

func (c *clientCodec) Decode(b []byte) error {
	var reqID = c.reqID
	if c.msg == nil {
		c.msg = &msg{}
	} else {
		c.msg.Reset()
	}
	err := c.msg.Decode(b)
	if c.msg.msgType != MsgType(pb.MsgType_res) {
		return errors.New("not be MsgType_res")
	}
	if c.msg.batch == false {
		res := &response{}
		err = res.Decode(c.msg.data)
		if err != nil {
			logger.Errorln("ResponseDecode error: ", err)
			return err
		}
		c.res = res
		if reqID == c.res.id {
			if res.err != nil {
				return res.err
			}
			return replyDecode(res.data, c.reply, c.msg.codecType)
		}
		return ErrReqID
	}
	c.batchCodec.Reset()
	err = c.batchCodec.Decode(c.msg.data)
	if err != nil {
		return err
	}
	if len(c.batchCodec.data) == len(c.requests) {
		for i, res := range c.responses {
			res.Decode(c.batchCodec.data[i])
		}
	}
	return nil
}
