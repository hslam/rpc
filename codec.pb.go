// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"fmt"
	"github.com/hslam/code"
)

type pbRequest struct {
	Seq           uint64
	Upgrade       []byte
	ServiceMethod string
	Args          []byte
}

// Size return the buf size.
func (req *pbRequest) Size() (n int) {
	var size uint64
	size += 11
	size += 11 + uint64(len(req.Upgrade))
	size += 11 + uint64(len(req.ServiceMethod))
	size += 11 + uint64(len(req.Args))
	return int(size)
}

//Marshal marshals the pbRequest and returns the bytes.
func (req *pbRequest) Marshal() ([]byte, error) {
	size := req.Size()
	buf := make([]byte, size)
	n, err := req.MarshalTo(buf[:size])
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

//MarshalTo marshals the pbRequest into buf and returns the bytes.
func (req *pbRequest) MarshalTo(buf []byte) (int, error) {
	var size uint64
	size += 11
	size += 11 + uint64(len(req.Upgrade))
	size += 11 + uint64(len(req.ServiceMethod))
	size += 11 + uint64(len(req.Args))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		return 0, fmt.Errorf("proto: pbRequest: buf is too short")
	}
	var offset uint64
	var n uint64
	if req.Seq != 0 {
		buf[offset] = 1<<3 | 0
		offset++
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
	}
	if len(req.Upgrade) > 0 {
		buf[offset] = 2<<3 | 2
		offset++
		//n = code.EncodeBytes(buf[offset:], req.Upgrade)
		{
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
		}
		offset += n
	}
	if len(req.ServiceMethod) > 0 {
		buf[offset] = 3<<3 | 2
		offset++
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
	}
	if len(req.Args) > 0 {
		buf[offset] = 4<<3 | 2
		offset++
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
	}
	return int(offset), nil
}

//Unmarshal unmarshals the pbRequest from buf and returns the number of bytes read (> 0).
func (req *pbRequest) Unmarshal(data []byte) error {
	var length = uint64(len(data))
	var offset uint64
	var n uint64
	var tag uint64
	var fieldNumber int
	var wireType uint8
	for {
		if offset < length {
			tag = uint64(data[offset])
			offset++
		} else {
			break
		}
		fieldNumber = int(tag >> 3)
		wireType = uint8(tag & 0x7)
		switch fieldNumber {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Seq", wireType)
			}
			n = code.DecodeVarint(data[offset:], &req.Seq)
			offset += n
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Upgrade", wireType)
			}
			n = code.DecodeBytes(data[offset:], &req.Upgrade)
			offset += n
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServiceMethod", wireType)
			}
			n = code.DecodeString(data[offset:], &req.ServiceMethod)
			offset += n
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			n = code.DecodeBytes(data[offset:], &req.Args)
			offset += n
		}
	}
	return nil
}

type pbResponse struct {
	Seq   uint64
	Error string
	Reply []byte
}

// Size return the buf size.
func (res *pbResponse) Size() (n int) {
	var size uint64
	size += 11
	size += 11 + uint64(len(res.Error))
	size += 11 + uint64(len(res.Reply))
	return int(size)
}

//Marshal marshals the pbResponse and returns the bytes.
func (res *pbResponse) Marshal() ([]byte, error) {
	size := res.Size()
	buf := make([]byte, size)
	n, err := res.MarshalTo(buf[:size])
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

//MarshalTo marshals the pbResponse into buf and returns the bytes.
func (res *pbResponse) MarshalTo(buf []byte) (int, error) {
	var size uint64
	size += 11
	size += 11 + uint64(len(res.Error))
	size += 11 + uint64(len(res.Reply))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		return 0, fmt.Errorf("proto: pbResponse: buf is too short")
	}
	var offset uint64
	var n uint64
	if res.Seq != 0 {
		buf[offset] = 1<<3 | 0
		offset++
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
	}
	if len(res.Error) > 0 {
		buf[offset] = 2<<3 | 2
		offset++
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
	}
	if len(res.Reply) > 0 {
		buf[offset] = 3<<3 | 2
		offset++
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
	}
	return int(offset), nil
}

//Unmarshal unmarshals the pbResponse from buf and returns the number of bytes read (> 0).
func (res *pbResponse) Unmarshal(data []byte) error {
	var length = uint64(len(data))
	var offset uint64
	var n uint64
	var tag uint64
	var fieldNumber int
	var wireType uint8
	for {
		if offset < length {
			tag = uint64(data[offset])
			offset++
		} else {
			break
		}
		fieldNumber = int(tag >> 3)
		wireType = uint8(tag & 0x7)
		switch fieldNumber {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Seq", wireType)
			}
			n = code.DecodeVarint(data[offset:], &res.Seq)
			offset += n
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			if data[offset] > 0 {
				n = code.DecodeString(data[offset:], &res.Error)
			} else {
				n = 1
			}
			offset += n
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reply", wireType)
			}
			n = code.DecodeBytes(data[offset:], &res.Reply)
			offset += n
		}
	}
	return nil
}
