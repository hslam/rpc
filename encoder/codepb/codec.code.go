package codepb

import (
	"fmt"
	"github.com/hslam/code"
)

type Request struct {
	Seq           uint64
	Upgrade       []byte
	ServiceMethod string
	Args          []byte
}

//Marshal marshals the Request into buf and returns the bytes.
func (req *Request) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 11
	size += 11 + uint64(len(req.Upgrade))
	size += 11 + uint64(len(req.ServiceMethod))
	size += 11 + uint64(len(req.Args))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
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
	return buf[:offset], nil
}

//Unmarshal unmarshals the Request from buf and returns the number of bytes read (> 0).
func (req *Request) Unmarshal(data []byte) (uint64, error) {
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
				return 0, fmt.Errorf("proto: wrong wireType = %d for field Seq", wireType)
			}
			n = code.DecodeVarint(data[offset:], &req.Seq)
			offset += n
		case 2:
			if wireType != 2 {
				return 0, fmt.Errorf("proto: wrong wireType = %d for field Upgrade", wireType)
			}
			n = code.DecodeBytes(data[offset:], &req.Upgrade)
			offset += n
		case 3:
			if wireType != 2 {
				return 0, fmt.Errorf("proto: wrong wireType = %d for field ServiceMethod", wireType)
			}
			n = code.DecodeString(data[offset:], &req.ServiceMethod)
			offset += n
		case 4:
			if wireType != 2 {
				return 0, fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			n = code.DecodeBytes(data[offset:], &req.Args)
			offset += n
		}
	}
	return offset, nil
}

type Response struct {
	Seq     uint64
	Upgrade []byte
	Error   string
	Reply   []byte
}

//Marshal marshals the Response into buf and returns the bytes.
func (res *Response) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 11
	size += 11 + uint64(len(res.Upgrade))
	size += 11 + uint64(len(res.Error))
	size += 11 + uint64(len(res.Reply))
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
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
	if len(res.Upgrade) > 0 {
		buf[offset] = 2<<3 | 2
		offset++
		//n = code.EncodeBytes(buf[offset:], res.Upgrade)
		{
			var length = uint64(len(res.Upgrade))
			var lengthSize = code.SizeofVarint(length)
			var s = lengthSize + length
			t := length
			for i := uint64(0); i < lengthSize-1; i++ {
				buf[offset+i] = byte(t) | 0x80
				t >>= 7
			}
			buf[offset+lengthSize-1] = byte(t)
			copy(buf[offset+lengthSize:], res.Upgrade)
			n = s
		}
		offset += n
	}
	if len(res.Error) > 0 {
		buf[offset] = 3<<3 | 2
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
		buf[offset] = 4<<3 | 2
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
	return buf[:offset], nil
}

//Unmarshal unmarshals the Response from buf and returns the number of bytes read (> 0).
func (res *Response) Unmarshal(data []byte) (uint64, error) {
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
				return 0, fmt.Errorf("proto: wrong wireType = %d for field Seq", wireType)
			}
			n = code.DecodeVarint(data[offset:], &res.Seq)
			offset += n
		case 2:
			if wireType != 2 {
				return 0, fmt.Errorf("proto: wrong wireType = %d for field Upgrade", wireType)
			}
			n = code.DecodeBytes(data[offset:], &res.Upgrade)
			offset += n
		case 3:
			if wireType != 2 {
				return 0, fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			if data[offset] > 0 {
				n = code.DecodeString(data[offset:], &res.Error)
			} else {
				n = 1
			}
			offset += n
		case 4:
			if wireType != 2 {
				return 0, fmt.Errorf("proto: wrong wireType = %d for field Reply", wireType)
			}
			n = code.DecodeBytes(data[offset:], &res.Reply)
			offset += n
		}
	}
	return offset, nil
}
