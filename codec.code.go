// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/code"
)

type request struct {
	Seq           uint64
	Upgrade       []byte
	ServiceMethod string
	Args          []byte
}

//Marshal marshals the Request into buf and returns the bytes.
func (req *request) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 10
	size += 10 + uint64(len(req.Upgrade))
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
	//n = code.EncodeBytes(buf[offset:], req.Upgrade)
	if len(req.Upgrade) > 127 {
		var length = uint64(len(req.Upgrade))
		var lengthSize = code.SizeofVarint(length)
		var s = lengthSize + length
		t := length
		for i := uint64(0); i < lengthSize-1; i++ {
			buf[offset+i] = byte(t) | 0x80
			t >>= 7
		}
		buf[offset+lengthSize-1] = byte(t)
		copy(buf[offset+lengthSize:], req.Upgrade)
		n = s
	} else if len(req.Upgrade) > 0 {
		var length = uint64(len(req.Upgrade))
		buf[offset] = byte(length)
		copy(buf[offset+1:], req.Upgrade)
		n = 1 + length
	} else {
		buf[offset] = 0
		n = 1
	}
	offset += n

	//n = code.EncodeString(buf[offset:], req.ServiceMethod)
	if len(req.ServiceMethod) > 127 {
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
	} else if len(req.ServiceMethod) > 0 {
		var length = uint64(len(req.ServiceMethod))
		buf[offset] = byte(length)
		copy(buf[offset+1:], req.ServiceMethod)
		n = 1 + length
	} else {
		buf[offset] = 0
		n = 1
	}
	offset += n
	//n = code.EncodeBytes(buf[offset:], req.Args)
	if len(req.Args) > 127 {
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
	} else if len(req.Args) > 0 {
		var length = uint64(len(req.Args))
		buf[offset] = byte(length)
		copy(buf[offset+1:], req.Args)
		n = 1 + length
	} else {
		buf[offset] = 0
		n = 1
	}
	offset += n
	return buf[:offset], nil
}

//Unmarshal unmarshals the Request from buf and returns the number of bytes read (> 0).
func (req *request) Unmarshal(data []byte) (uint64, error) {
	var offset uint64
	var n uint64
	n = code.DecodeVarint(data[offset:], &req.Seq)
	offset += n
	if data[offset] > 127 {
		n = code.DecodeBytes(data[offset:], &req.Upgrade)
	} else if data[offset] > 0 {
		var s = 1 + uint64(data[offset])
		req.Upgrade = data[offset+1 : offset+s]
		n = s
	} else {
		n = 1
	}
	offset += n
	if data[offset] > 0 {
		n = code.DecodeString(data[offset:], &req.ServiceMethod)
	} else {
		n = 1
	}
	offset += n
	if data[offset] > 127 {
		n = code.DecodeBytes(data[offset:], &req.Args)
	} else if data[offset] > 0 {
		var s = 1 + uint64(data[offset])
		req.Args = data[offset+1 : offset+s]
		n = s
	} else {
		n = 1
	}
	offset += n
	return offset, nil
}

type response struct {
	Seq   uint64
	Error string
	Reply []byte
}

//Marshal marshals the Response into buf and returns the bytes.
func (res *response) Marshal(buf []byte) ([]byte, error) {
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
	if len(res.Error) > 127 {
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
	} else if len(res.Error) > 0 {
		var length = uint64(len(res.Error))
		buf[offset] = byte(length)
		copy(buf[offset+1:], res.Error)
		n = 1 + length
	} else {
		buf[offset] = 0
		n = 1
	}
	offset += n
	//n = code.EncodeBytes(buf[offset:], res.Reply)
	if len(res.Reply) > 127 {
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
	} else if len(res.Reply) > 0 {
		var length = uint64(len(res.Reply))
		buf[offset] = byte(length)
		copy(buf[offset+1:], res.Reply)
		n = 1 + length
	} else {
		buf[offset] = 0
		n = 1
	}
	offset += n
	return buf[:offset], nil
}

//Unmarshal unmarshals the Response from buf and returns the number of bytes read (> 0).
func (res *response) Unmarshal(data []byte) (uint64, error) {
	var offset uint64
	var n uint64
	n = code.DecodeVarint(data[offset:], &res.Seq)
	offset += n
	if data[offset] > 0 {
		n = code.DecodeString(data[offset:], &res.Error)
	} else {
		n = 1
	}
	offset += n
	if data[offset] > 127 {
		n = code.DecodeBytes(data[offset:], &res.Reply)
	} else if data[offset] > 0 {
		var s = 1 + uint64(data[offset])
		res.Reply = data[offset+1 : offset+s]
		n = s
	} else {
		n = 1
	}
	offset += n
	return offset, nil
}
