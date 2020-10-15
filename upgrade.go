// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import "errors"

const (
	// NoRequest represents a no request transaction.
	noRequest = 0x1
	// NoResponse represents a no response transaction.
	noResponse = 0x1
	// FlateCompress represents a flate compress transaction.
	flateCompress = 0x1
	// ZlibCompress represents a zlib compress transaction.
	zlibCompress = 0x2
	// GzipCompress represents a gzip compress transaction.
	gzipCompress = 0x3
	// Heartbeat represents a heartbeat transaction.
	heartbeat = 0x1
	// startWatch represents a start watch transaction.
	watch = 0x1
	// StopWatch represents a stop watch transaction.
	stopWatch = 0x2
)

type upgrade struct {
	NoRequest  byte
	NoResponse byte
	Compress   byte
	Heartbeat  byte
	Watch      byte
	Reserve    byte
}

func (u *upgrade) Reset() {
	*u = upgrade{}
}

func (u *upgrade) Marshal(buf []byte) ([]byte, error) {
	var size uint64 = 1
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	buf[0] = u.NoRequest<<7 + u.NoResponse<<6 + u.Compress<<4 + u.Heartbeat<<3 + u.Watch<<1 + u.Reserve
	offset++
	return buf[:offset], nil
}

func (u *upgrade) Unmarshal(data []byte) (uint64, error) {
	var offset uint64
	if uint64(len(data)) < offset+1 {
		return 0, errors.New("data is too short")
	}
	u.NoRequest = data[0] >> 7 & 0x1
	u.NoResponse = data[0] >> 6 & 0x1
	u.Compress = data[0] >> 4 & 0x3
	u.Heartbeat = data[0] >> 3 & 0x1
	u.Watch = data[0] >> 1 & 0x3
	u.Reserve = data[0] & 0x1
	offset++
	return offset, nil
}
