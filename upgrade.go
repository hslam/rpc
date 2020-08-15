// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

const (
	Heartbeat     = 0x1
	NoRequest     = 0x1
	NoResponse    = 0x1
	FlateCompress = 0x1
	ZlibCompress  = 0x2
	GzipCompress  = 0x3
)

type upgrade struct {
	Heartbeat  byte
	NoRequest  byte
	NoResponse byte
	Compress   byte
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
	buf[0] = u.Heartbeat<<7 + u.NoRequest<<6 + u.NoResponse<<5 + u.Compress<<3 + u.Reserve
	offset += 1
	return buf[:offset], nil
}

func (u *upgrade) Unmarshal(data []byte) (uint64, error) {
	var offset uint64
	if uint64(len(data)) < offset+1 {
		return 0, nil
	}
	u.Heartbeat = data[0] >> 7
	u.NoRequest = data[0] >> 6 & 0x1
	u.NoResponse = data[0] >> 5 & 0x1
	u.Compress = data[0] >> 3 & 0x3
	u.Reserve = data[0] & 0x7
	offset++
	return offset, nil
}
