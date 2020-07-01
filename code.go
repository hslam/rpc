package rpc

import (
	"github.com/hslam/code"
)

type request struct {
	Seq           uint64
	ServiceMethod string
	Args          []byte
}

//Reset resets the Request.
func (req *request) Reset() {
	*req = request{}
}

//SetSeq sets the value of Seq.
func (req *request) SetSeq(seq uint64) {
	req.Seq = seq
}

//GetSeq returns the value of Seq.
func (req *request) GetSeq() uint64 {
	return req.Seq
}

//SetServiceMethod sets the value of ServiceMethod.
func (req *request) SetServiceMethod(serviceMethod string) {
	req.ServiceMethod = serviceMethod
}

//GetServiceMethod returns the value of ServiceMethod.
func (req *request) GetServiceMethod() string {
	return req.ServiceMethod
}

//SetArgs sets the value of Args.
func (req *request) SetArgs(args []byte) {
	req.Args = args
}

//GetArgs returns the value of Args.
func (req *request) GetArgs() []byte {
	return req.Args
}

//Marshal marshals the Request into buf and returns the bytes.
func (req *request) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 8
	size += 10 + uint64(len(req.ServiceMethod))
	size += 10 + uint64(len(req.Args))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	{
		var t = req.Seq
		buf[0] = uint8(t)
		buf[1] = uint8(t >> 8)
		buf[2] = uint8(t >> 16)
		buf[3] = uint8(t >> 24)
		buf[4] = uint8(t >> 32)
		buf[5] = uint8(t >> 40)
		buf[6] = uint8(t >> 48)
		buf[7] = uint8(t >> 56)
		n = 8
	}
	offset += n
	{
		var length = uint64(len(req.ServiceMethod))
		var lengthSize = code.SizeofVarint(length)
		var s = lengthSize + length
		t := length
		for i := uint64(0); i < lengthSize-1; i++ {
			buf[offset+i] = byte(t) | 0x80
			t >>= 7
		}
		buf[offset+lengthSize-1] = byte(t)
		copy(buf[offset+lengthSize:], req.ServiceMethod)
		n = s
	}
	offset += n
	{
		var length = uint64(len(req.Args))
		var lengthSize = code.SizeofVarint(length)
		var s = lengthSize + length
		t := length
		for i := uint64(0); i < lengthSize-1; i++ {
			buf[offset+i] = byte(t) | 0x80
			t >>= 7
		}
		buf[offset+lengthSize-1] = byte(t)
		copy(buf[offset+lengthSize:], req.Args)
		n = s
	}
	offset += n
	return buf[:offset], nil
}

//Unmarshal unmarshals the Request from buf and returns the number of bytes read (> 0).
func (req *request) Unmarshal(data []byte) (uint64, error) {
	var offset uint64
	var n uint64
	n = code.DecodeUint64(data[offset:], &req.Seq)
	offset += n
	n = code.DecodeString(data[offset:], &req.ServiceMethod)
	offset += n
	n = code.DecodeBytes(data[offset:], &req.Args)
	offset += n
	return offset, nil
}

type response struct {
	Seq   uint64
	Error string
	Reply []byte
}

//Reset resets the Response.
func (res *response) Reset() {
	*res = response{}
}

//SetSeq sets the value of Seq.
func (res *response) SetSeq(seq uint64) {
	res.Seq = seq
}

//GetSeq returns the value of Seq.
func (res *response) GetSeq() uint64 {
	return res.Seq
}

//SetError sets the value of Error.
func (res *response) SetError(errorMsg string) {
	res.Error = errorMsg
}

//GetError returns the value of Error.
func (res *response) GetError() string {
	return res.Error
}

//SetReply sets the value of Reply.
func (res *response) SetReply(reply []byte) {
	res.Reply = reply
}

//GetReply returns the value of Reply.
func (res *response) GetReply() []byte {
	return res.Reply
}

//Marshal marshals the Response into buf and returns the bytes.
func (res *response) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 8
	size += 10 + uint64(len(res.Error))
	size += 10 + uint64(len(res.Reply))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	{
		var t = res.Seq
		buf[0] = uint8(t)
		buf[1] = uint8(t >> 8)
		buf[2] = uint8(t >> 16)
		buf[3] = uint8(t >> 24)
		buf[4] = uint8(t >> 32)
		buf[5] = uint8(t >> 40)
		buf[6] = uint8(t >> 48)
		buf[7] = uint8(t >> 56)
		n = 8
	}
	offset += n
	{
		var length = uint64(len(res.Error))
		var lengthSize = code.SizeofVarint(length)
		var s = lengthSize + length
		t := length
		for i := uint64(0); i < lengthSize-1; i++ {
			buf[offset+i] = byte(t) | 0x80
			t >>= 7
		}
		buf[offset+lengthSize-1] = byte(t)
		copy(buf[offset+lengthSize:], res.Error)
		n = s
	}
	offset += n
	{
		var length = uint64(len(res.Reply))
		var lengthSize = code.SizeofVarint(length)
		var s = lengthSize + length
		t := length
		for i := uint64(0); i < lengthSize-1; i++ {
			buf[offset+i] = byte(t) | 0x80
			t >>= 7
		}
		buf[offset+lengthSize-1] = byte(t)
		copy(buf[offset+lengthSize:], res.Reply)
		n = s
	}
	offset += n
	return buf[:offset], nil
}

//Unmarshal unmarshals the Response from buf and returns the number of bytes read (> 0).
func (res *response) Unmarshal(data []byte) (uint64, error) {
	var offset uint64
	var n uint64
	n = code.DecodeUint64(data[offset:], &res.Seq)
	offset += n
	n = code.DecodeString(data[offset:], &res.Error)
	offset += n
	n = code.DecodeBytes(data[offset:], &res.Reply)
	offset += n
	return offset, nil
}
