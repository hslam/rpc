package rpc

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"compress/gzip"
	"io"
)
type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Uncompress(data []byte) ([]byte, error)
}

type FlateCompressor struct{
	Level CompressLevel
}

func (c FlateCompressor) Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	var(
		write *flate.Writer
		err error
	)
	switch c.Level {
	case BestSpeed:
		write, err = flate.NewWriter(buf, flate.BestSpeed)
	case BestCompression:
		write, err = flate.NewWriter(buf, flate.BestCompression)
	case DefaultCompression:
		write, err = flate.NewWriter(buf, flate.DefaultCompression)
	default:
		return data,nil
	}
	if err != nil {
		return nil,err
	}
	defer write.Close()
	_,err=write.Write(data)
	if err != nil {
		return nil,err
	}
	err=write.Flush()
	if err != nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func (c FlateCompressor) Uncompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	reader := flate.NewReader(buf)
	defer reader.Close()
	var out bytes.Buffer
	_,err:=io.Copy(&out,reader)
	if err != nil&&err!=io.ErrUnexpectedEOF {
		return nil,err
	}
	return out.Bytes(),nil
}

type ZlibCompressor struct{
	Level CompressLevel
}

func (c ZlibCompressor) Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	var(
		write *zlib.Writer
		err error
	)
	switch c.Level {
	case BestSpeed:
		write, err = zlib.NewWriterLevel(buf, zlib.BestSpeed)
	case BestCompression:
		write, err = zlib.NewWriterLevel(buf, zlib.BestCompression)
	case DefaultCompression:
		write, err = zlib.NewWriterLevel(buf, zlib.DefaultCompression)
	default:
		return data,nil
	}
	if err != nil {
		return nil,err
	}
	defer write.Close()
	_,err=write.Write(data)
	if err != nil {
		return nil,err
	}
	err=write.Flush()
	if err != nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func (c ZlibCompressor) Uncompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	reader, err:= zlib.NewReader(buf)
	defer reader.Close()
	if err != nil {
		return nil,err
	}
	var out bytes.Buffer
	_,err=io.Copy(&out,reader)
	if err != nil&&err!=io.ErrUnexpectedEOF {
		return nil,err
	}
	return out.Bytes(),nil
}



type GzipCompressor struct{
	Level CompressLevel
}

func (c GzipCompressor) Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	var(
		write *gzip.Writer
		err error
	)
	switch c.Level {
	case BestSpeed:
		write, err = gzip.NewWriterLevel(buf, gzip.BestSpeed)
	case BestCompression:
		write, err = gzip.NewWriterLevel(buf, gzip.BestCompression)
	case DefaultCompression:
		write, err = gzip.NewWriterLevel(buf, gzip.DefaultCompression)
	default:
		return data,nil
	}
	if err != nil {
		return nil,err
	}
	defer write.Close()
	_,err=write.Write(data)
	if err != nil {
		return nil,err
	}
	err=write.Flush()
	if err != nil {
		return nil,err
	}
	return buf.Bytes(),nil
}

func (c GzipCompressor) Uncompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	reader, err:= gzip.NewReader(buf)
	defer reader.Close()
	if err != nil {
		return nil,err
	}
	var out bytes.Buffer
	_,err=io.Copy(&out,reader)
	if err != nil&&err!=io.ErrUnexpectedEOF {
		return nil,err
	}
	return out.Bytes(),nil
}