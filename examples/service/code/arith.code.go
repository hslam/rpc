package code

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
	n = code.EncodeUint32(buf[offset:], uint32(a.A))
	offset += n
	n = code.EncodeUint32(buf[offset:], uint32(a.B))
	offset += n
	return buf, nil
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
	Pro int32
	Quo int32
	Rem int32
}

//Marshal takes a buffer and encodes the ArithResponse to bytes
func (a *ArithResponse) Marshal(buf []byte) ([]byte, error) {
	var size uint64
	size += 4
	size += 4
	size += 4
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	var n uint64
	n = code.EncodeUint32(buf[offset:], uint32(a.Pro))
	offset += n
	n = code.EncodeUint32(buf[offset:], uint32(a.Quo))
	offset += n
	n = code.EncodeUint32(buf[offset:], uint32(a.Rem))
	offset += n
	return buf, nil
}

//Unmarshal parses the encoded data.
func (a *ArithResponse) Unmarshal(b []byte) (uint64, error) {
	var offset uint64
	var n uint64
	var Pro uint32
	n = code.DecodeUint32(b[offset:], &Pro)
	a.Pro = int32(Pro)
	offset += n
	var Quo uint32
	n = code.DecodeUint32(b[offset:], &Quo)
	a.Quo = int32(Quo)
	offset += n
	var Rem uint32
	n = code.DecodeUint32(b[offset:], &Rem)
	a.Rem = int32(Rem)
	offset += n
	return offset, nil
}
