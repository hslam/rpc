package service

import (
	"github.com/hslam/code"
)

//ArithRequest defines the request of arith.
type ArithRequest struct {
	A int32
	B int32
}

//Marshal takes a buffer and encodes the ArithRequest to bytes
func (a *ArithRequest) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 4
	size += 4
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	{
		var t = a.A
		buf[0] = uint8(t)
		buf[1] = uint8(t >> 8)
		buf[2] = uint8(t >> 16)
		buf[3] = uint8(t >> 24)
		n = 4
	}
	offset += n
	{
		var t = a.B
		buf[offset+0] = uint8(t)
		buf[offset+1] = uint8(t >> 8)
		buf[offset+2] = uint8(t >> 16)
		buf[offset+3] = uint8(t >> 24)
		n = 4
	}
	offset += n
	return buf[:offset], nil
}

//Unmarshal parses the encoded data.
func (a *ArithRequest) Unmarshal(b []byte) (uint64, error) {
	var offset uint64
	var n uint64
	var A uint32
	n = code.DecodeUint32(b[offset:], &A)
	a.A = int32(A)
	offset += n
	var B uint32
	n = code.DecodeUint32(b[offset:], &B)
	a.B = int32(B)
	offset += n
	return offset, nil
}

//ArithResponse defines the response of arith.
type ArithResponse struct {
	Pro int32 //product
}

//Marshal takes a buffer and encodes the ArithResponse to bytes
func (a *ArithResponse) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 4
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	{
		var t = a.Pro
		buf[0] = uint8(t)
		buf[1] = uint8(t >> 8)
		buf[2] = uint8(t >> 16)
		buf[3] = uint8(t >> 24)
		n = 4
	}
	offset += n
	return buf[:offset], nil
}

//Unmarshal parses the encoded data.
func (a *ArithResponse) Unmarshal(b []byte) (uint64, error) {
	var offset uint64
	var n uint64
	var Pro uint32
	n = code.DecodeUint32(b[offset:], &Pro)
	a.Pro = int32(Pro)
	offset += n
	return offset, nil
}
