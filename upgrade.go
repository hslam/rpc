// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import "errors"

const (
	upgradeSize = 1

	noRequest   = 0x1
	noResponse  = 0x1
	heartbeat   = 0x1
	openStream  = 0x1
	streaming   = 0x2
	closeStream = 0x3
)

type upgrade struct {
	NoRequest  byte //1bit
	NoResponse byte //1bit
	Heartbeat  byte //1bit
	Stream     byte //2bit
}

func (u *upgrade) Reset() {
	*u = upgrade{}
}

func (u *upgrade) IsZero() bool {
	return u.NoRequest+u.NoResponse+u.Heartbeat+u.Stream == 0
}

func (u *upgrade) Marshal(buf []byte) ([]byte, error) {
	var size uint64 = upgradeSize
	if uint64(cap(buf)) >= size {
		buf = buf[:size]
	} else {
		buf = make([]byte, size)
	}
	var offset uint64
	buf[0] = u.NoRequest<<7 + u.NoResponse<<6 + u.Heartbeat<<5 + u.Stream<<3
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
	u.Heartbeat = data[0] >> 5 & 0x1
	u.Stream = data[0] >> 3 & 0x3
	offset++
	return offset, nil
}
