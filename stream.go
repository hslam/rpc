package rpc

import (
	"io"
)

type Stream interface {
	SetReader(reader io.Reader) Stream
	SetWriter(writer io.Writer) Stream
	SetCloser(closer io.Closer) Stream
	ReadMessage() (p []byte, err error)
	WriteMessage(b []byte) (err error)
	Close() error
}

type stream struct {
	Reader io.Reader
	Writer io.Writer
	Closer io.Closer
	Send   []byte
	Read   []byte
	buffer []byte
}

func NewStream(r io.Reader, w io.Writer, c io.Closer, bufferSize int) Stream {
	if bufferSize < 1 {
		bufferSize = 1024
	}
	return &stream{
		Reader: r,
		Writer: w,
		Closer: c,
		Send:   make([]byte, bufferSize+8),
		Read:   make([]byte, bufferSize),
	}
}
func (s *stream) SetReader(reader io.Reader) Stream {
	s.Reader = reader
	return s
}

func (s *stream) SetWriter(writer io.Writer) Stream {
	s.Writer = writer
	return s
}

func (s *stream) SetCloser(closer io.Closer) Stream {
	s.Closer = closer
	return s
}

func (s *stream) ReadMessage() (p []byte, err error) {
	for {
		length := uint64(len(s.buffer))
		var i uint64 = 0
		for i < length {
			if length < i+8 {
				break
			}
			var msgLength uint64
			buf := s.buffer[i : i+8]
			var t uint64
			t = uint64(buf[0])
			t |= uint64(buf[1]) << 8
			t |= uint64(buf[2]) << 16
			t |= uint64(buf[3]) << 24
			t |= uint64(buf[4]) << 32
			t |= uint64(buf[5]) << 40
			t |= uint64(buf[6]) << 48
			t |= uint64(buf[7]) << 56
			msgLength = t
			if length < i+8+msgLength {
				break
			}
			p = s.buffer[i+8 : i+8+msgLength]
			i += 8 + msgLength
			break
		}
		s.buffer = s.buffer[i:]
		if i > 0 {
			break
		}
		n, err := s.Reader.Read(s.Read)
		if err != nil {
			return p, err
		}
		if n > 0 {
			s.buffer = append(s.buffer, s.Read[:n]...)
		}
	}
	return
}

func (s *stream) WriteMessage(b []byte) error {
	var length = uint64(len(b))
	var size = 8 + length
	if uint64(cap(s.Send)) >= size {
		s.Send = s.Send[:size]
	} else {
		s.Send = make([]byte, size)
	}
	var t = length
	var buf = s.Send[0:8]
	buf[0] = uint8(t)
	buf[1] = uint8(t >> 8)
	buf[2] = uint8(t >> 16)
	buf[3] = uint8(t >> 24)
	buf[4] = uint8(t >> 32)
	buf[5] = uint8(t >> 40)
	buf[6] = uint8(t >> 48)
	buf[7] = uint8(t >> 56)
	copy(s.Send[8:], b)
	_, err := s.Writer.Write(s.Send[:size])
	return err
}

func (s *stream) Close() error {
	return s.Closer.Close()
}
