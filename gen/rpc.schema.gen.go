package gen

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

type Msg struct {
	Version       float32
	Id            int64
	MsgType       int32
	Batch         bool
	CodecType     int32
	CompressType  int32
	CompressLevel int32
	Data          []byte
}

func (d *Msg) Size() (s uint64) {

	{
		l := uint64(len(d.Data))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 29
	return
}
func (d *Msg) Marshal(buf []byte) ([]byte, error) {
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

		v := *(*uint32)(unsafe.Pointer(&(d.Version)))

		buf[0+0] = byte(v >> 0)

		buf[1+0] = byte(v >> 8)

		buf[2+0] = byte(v >> 16)

		buf[3+0] = byte(v >> 24)

	}
	{

		buf[0+4] = byte(d.Id >> 0)

		buf[1+4] = byte(d.Id >> 8)

		buf[2+4] = byte(d.Id >> 16)

		buf[3+4] = byte(d.Id >> 24)

		buf[4+4] = byte(d.Id >> 32)

		buf[5+4] = byte(d.Id >> 40)

		buf[6+4] = byte(d.Id >> 48)

		buf[7+4] = byte(d.Id >> 56)

	}
	{

		buf[0+12] = byte(d.MsgType >> 0)

		buf[1+12] = byte(d.MsgType >> 8)

		buf[2+12] = byte(d.MsgType >> 16)

		buf[3+12] = byte(d.MsgType >> 24)

	}
	{
		if d.Batch {
			buf[16] = 1
		} else {
			buf[16] = 0
		}
	}
	{

		buf[0+17] = byte(d.CodecType >> 0)

		buf[1+17] = byte(d.CodecType >> 8)

		buf[2+17] = byte(d.CodecType >> 16)

		buf[3+17] = byte(d.CodecType >> 24)

	}
	{

		buf[0+21] = byte(d.CompressType >> 0)

		buf[1+21] = byte(d.CompressType >> 8)

		buf[2+21] = byte(d.CompressType >> 16)

		buf[3+21] = byte(d.CompressType >> 24)

	}
	{

		buf[0+25] = byte(d.CompressLevel >> 0)

		buf[1+25] = byte(d.CompressLevel >> 8)

		buf[2+25] = byte(d.CompressLevel >> 16)

		buf[3+25] = byte(d.CompressLevel >> 24)

	}
	{
		l := uint64(len(d.Data))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+29] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+29] = byte(t)
			i++

		}
		copy(buf[i+29:], d.Data)
		i += l
	}
	return buf[:i+29], nil
}

func (d *Msg) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		v := 0 | (uint32(buf[i+0+0]) << 0) | (uint32(buf[i+1+0]) << 8) | (uint32(buf[i+2+0]) << 16) | (uint32(buf[i+3+0]) << 24)
		d.Version = *(*float32)(unsafe.Pointer(&v))

	}
	{

		d.Id = 0 | (int64(buf[i+0+4]) << 0) | (int64(buf[i+1+4]) << 8) | (int64(buf[i+2+4]) << 16) | (int64(buf[i+3+4]) << 24) | (int64(buf[i+4+4]) << 32) | (int64(buf[i+5+4]) << 40) | (int64(buf[i+6+4]) << 48) | (int64(buf[i+7+4]) << 56)

	}
	{

		d.MsgType = 0 | (int32(buf[i+0+12]) << 0) | (int32(buf[i+1+12]) << 8) | (int32(buf[i+2+12]) << 16) | (int32(buf[i+3+12]) << 24)

	}
	{
		d.Batch = buf[i+16] == 1
	}
	{

		d.CodecType = 0 | (int32(buf[i+0+17]) << 0) | (int32(buf[i+1+17]) << 8) | (int32(buf[i+2+17]) << 16) | (int32(buf[i+3+17]) << 24)

	}
	{

		d.CompressType = 0 | (int32(buf[i+0+21]) << 0) | (int32(buf[i+1+21]) << 8) | (int32(buf[i+2+21]) << 16) | (int32(buf[i+3+21]) << 24)

	}
	{

		d.CompressLevel = 0 | (int32(buf[i+0+25]) << 0) | (int32(buf[i+1+25]) << 8) | (int32(buf[i+2+25]) << 16) | (int32(buf[i+3+25]) << 24)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+29] & 0x7F)
			for buf[i+29]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+29]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Data)) >= l {
			d.Data = d.Data[:l]
		} else {
			d.Data = make([]byte, l)
		}
		copy(d.Data, buf[i+29:])
		i += l
	}
	return i + 29, nil
}

type Batch struct {
	Async bool
	Data  [][]byte
}

func (d *Batch) Size() (s uint64) {

	{
		l := uint64(len(d.Data))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}

		for k0 := range d.Data {

			{
				l := uint64(len(d.Data[k0]))

				{

					t := l
					for t >= 0x80 {
						t >>= 7
						s++
					}
					s++

				}
				s += l
			}

		}

	}
	s += 1
	return
}
func (d *Batch) Marshal(buf []byte) ([]byte, error) {
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
		if d.Async {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
	}
	{
		l := uint64(len(d.Data))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+1] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+1] = byte(t)
			i++

		}
		for k0 := range d.Data {

			{
				l := uint64(len(d.Data[k0]))

				{

					t := uint64(l)

					for t >= 0x80 {
						buf[i+1] = byte(t) | 0x80
						t >>= 7
						i++
					}
					buf[i+1] = byte(t)
					i++

				}
				copy(buf[i+1:], d.Data[k0])
				i += l
			}

		}
	}
	return buf[:i+1], nil
}

func (d *Batch) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		d.Async = buf[i+0] == 1
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+1] & 0x7F)
			for buf[i+1]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+1]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Data)) >= l {
			d.Data = d.Data[:l]
		} else {
			d.Data = make([][]byte, l)
		}
		for k0 := range d.Data {

			{
				l := uint64(0)

				{

					bs := uint8(7)
					t := uint64(buf[i+1] & 0x7F)
					for buf[i+1]&0x80 == 0x80 {
						i++
						t |= uint64(buf[i+1]&0x7F) << bs
						bs += 7
					}
					i++

					l = t

				}
				if uint64(cap(d.Data[k0])) >= l {
					d.Data[k0] = d.Data[k0][:l]
				} else {
					d.Data[k0] = make([]byte, l)
				}
				copy(d.Data[k0], buf[i+1:])
				i += l
			}

		}
	}
	return i + 1, nil
}

type Request struct {
	Id         uint64
	Method     string
	NoRequest  bool
	NoResponse bool
	Data       []byte
}

func (d *Request) Size() (s uint64) {

	{
		l := uint64(len(d.Method))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Data))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 10
	return
}
func (d *Request) Marshal(buf []byte) ([]byte, error) {
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

		buf[0+0] = byte(d.Id >> 0)

		buf[1+0] = byte(d.Id >> 8)

		buf[2+0] = byte(d.Id >> 16)

		buf[3+0] = byte(d.Id >> 24)

		buf[4+0] = byte(d.Id >> 32)

		buf[5+0] = byte(d.Id >> 40)

		buf[6+0] = byte(d.Id >> 48)

		buf[7+0] = byte(d.Id >> 56)

	}
	{
		l := uint64(len(d.Method))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Method)
		i += l
	}
	{
		if d.NoRequest {
			buf[i+8] = 1
		} else {
			buf[i+8] = 0
		}
	}
	{
		if d.NoResponse {
			buf[i+9] = 1
		} else {
			buf[i+9] = 0
		}
	}
	{
		l := uint64(len(d.Data))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+10] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+10] = byte(t)
			i++

		}
		copy(buf[i+10:], d.Data)
		i += l
	}
	return buf[:i+10], nil
}

func (d *Request) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Id = 0 | (uint64(buf[i+0+0]) << 0) | (uint64(buf[i+1+0]) << 8) | (uint64(buf[i+2+0]) << 16) | (uint64(buf[i+3+0]) << 24) | (uint64(buf[i+4+0]) << 32) | (uint64(buf[i+5+0]) << 40) | (uint64(buf[i+6+0]) << 48) | (uint64(buf[i+7+0]) << 56)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Method = string(buf[i+8 : i+8+l])
		i += l
	}
	{
		d.NoRequest = buf[i+8] == 1
	}
	{
		d.NoResponse = buf[i+9] == 1
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+10] & 0x7F)
			for buf[i+10]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+10]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Data)) >= l {
			d.Data = d.Data[:l]
		} else {
			d.Data = make([]byte, l)
		}
		copy(d.Data, buf[i+10:])
		i += l
	}
	return i + 10, nil
}

type Response struct {
	Id     uint64
	Data   []byte
	ErrMsg string
}

func (d *Response) Size() (s uint64) {

	{
		l := uint64(len(d.Data))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.ErrMsg))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 8
	return
}
func (d *Response) Marshal(buf []byte) ([]byte, error) {
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

		buf[0+0] = byte(d.Id >> 0)

		buf[1+0] = byte(d.Id >> 8)

		buf[2+0] = byte(d.Id >> 16)

		buf[3+0] = byte(d.Id >> 24)

		buf[4+0] = byte(d.Id >> 32)

		buf[5+0] = byte(d.Id >> 40)

		buf[6+0] = byte(d.Id >> 48)

		buf[7+0] = byte(d.Id >> 56)

	}
	{
		l := uint64(len(d.Data))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Data)
		i += l
	}
	{
		l := uint64(len(d.ErrMsg))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.ErrMsg)
		i += l
	}
	return buf[:i+8], nil
}

func (d *Response) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Id = 0 | (uint64(buf[i+0+0]) << 0) | (uint64(buf[i+1+0]) << 8) | (uint64(buf[i+2+0]) << 16) | (uint64(buf[i+3+0]) << 24) | (uint64(buf[i+4+0]) << 32) | (uint64(buf[i+5+0]) << 40) | (uint64(buf[i+6+0]) << 48) | (uint64(buf[i+7+0]) << 56)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Data)) >= l {
			d.Data = d.Data[:l]
		} else {
			d.Data = make([]byte, l)
		}
		copy(d.Data, buf[i+8:])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.ErrMsg = string(buf[i+8 : i+8+l])
		i += l
	}
	return i + 8, nil
}
