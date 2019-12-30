package code

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type ArithRequest struct {
	A int32
	B int32
}

func (d *ArithRequest) Size() (s uint64) {

	s += 8
	return
}
func (d *ArithRequest) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.A >> 0)

		buf[1+0] = byte(d.A >> 8)

		buf[2+0] = byte(d.A >> 16)

		buf[3+0] = byte(d.A >> 24)

	}
	{

		buf[0+4] = byte(d.B >> 0)

		buf[1+4] = byte(d.B >> 8)

		buf[2+4] = byte(d.B >> 16)

		buf[3+4] = byte(d.B >> 24)

	}
	return buf[:i+8], nil
}

func (d *ArithRequest) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.A = 0 | (int32(buf[0+0]) << 0) | (int32(buf[1+0]) << 8) | (int32(buf[2+0]) << 16) | (int32(buf[3+0]) << 24)

	}
	{

		d.B = 0 | (int32(buf[0+4]) << 0) | (int32(buf[1+4]) << 8) | (int32(buf[2+4]) << 16) | (int32(buf[3+4]) << 24)

	}
	return i + 8, nil
}

type ArithResponse struct {
	Pro int32
	Quo int32
	Rem int32
}

func (d *ArithResponse) Size() (s uint64) {

	s += 12
	return
}
func (d *ArithResponse) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.Pro >> 0)

		buf[1+0] = byte(d.Pro >> 8)

		buf[2+0] = byte(d.Pro >> 16)

		buf[3+0] = byte(d.Pro >> 24)

	}
	{

		buf[0+4] = byte(d.Quo >> 0)

		buf[1+4] = byte(d.Quo >> 8)

		buf[2+4] = byte(d.Quo >> 16)

		buf[3+4] = byte(d.Quo >> 24)

	}
	{

		buf[0+8] = byte(d.Rem >> 0)

		buf[1+8] = byte(d.Rem >> 8)

		buf[2+8] = byte(d.Rem >> 16)

		buf[3+8] = byte(d.Rem >> 24)

	}
	return buf[:i+12], nil
}

func (d *ArithResponse) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Pro = 0 | (int32(buf[0+0]) << 0) | (int32(buf[1+0]) << 8) | (int32(buf[2+0]) << 16) | (int32(buf[3+0]) << 24)

	}
	{

		d.Quo = 0 | (int32(buf[0+4]) << 0) | (int32(buf[1+4]) << 8) | (int32(buf[2+4]) << 16) | (int32(buf[3+4]) << 24)

	}
	{

		d.Rem = 0 | (int32(buf[0+8]) << 0) | (int32(buf[1+8]) << 8) | (int32(buf[2+8]) << 16) | (int32(buf[3+8]) << 24)

	}
	return i + 12, nil
}
