package code

import (
	"github.com/hslam/code"
)

type Request struct {
	Seq           uint64
	ServiceMethod string
	Args          []byte
}

//Marshal marshals the Request into buf and returns the bytes.
func (req *Request) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 10
	size += 10 + uint64(len(req.ServiceMethod))
	size += 10 + uint64(len(req.Args))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	//n = code.EncodeVarint(buf[offset:], req.Seq)
	{
		var t = req.Seq
		var size = code.SizeofVarint(t)
		for i := uint64(0); i < size-1; i++ {
			buf[offset+i] = byte(t) | 0x80
			t >>= 7
		}
		buf[offset+size-1] = byte(t)
		n = size
	}
	offset += n
	//n = code.EncodeString(buf[offset:], req.ServiceMethod)
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
	//n = code.EncodeBytes(buf[offset:], req.Args)
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
func (req *Request) Unmarshal(data []byte) (uint64, error) {
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

type Response struct {
	Seq   uint64
	Error string
	Reply []byte
}

//Marshal marshals the Response into buf and returns the bytes.
func (res *Response) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 10
	size += 10 + uint64(len(res.Error))
	size += 10 + uint64(len(res.Reply))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	//n = code.EncodeVarint(buf[offset:], res.Seq)
	{
		var t = res.Seq
		var size = code.SizeofVarint(t)
		for i := uint64(0); i < size-1; i++ {
			buf[offset+i] = byte(t) | 0x80
			t >>= 7
		}
		buf[offset+size-1] = byte(t)
		n = size
	}
	offset += n
	//n = code.EncodeString(buf[offset:], res.Error)
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
	//n = code.EncodeBytes(buf[offset:], res.Reply)
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
func (res *Response) Unmarshal(data []byte) (uint64, error) {
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
