// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"testing"
)

func TestUpgrade(t *testing.T) {
	testUpgrade(nil, t)
	testUpgrade(make([]byte, 1), t)
}

func testUpgrade(buf []byte, t *testing.T) {
	u := &upgrade{
		NoRequest:  noRequest,
		NoResponse: noResponse,
		Compress:   flateCompress,
		Heartbeat:  heartbeat,
		Watch:      watch,
		Reserve:    1,
	}
	data, err := u.Marshal(buf)
	if err != nil {
		t.Error(err)
	}
	u.Reset()
	n, err := u.Unmarshal(data)
	if err != nil {
		t.Error(err)
	}
	if n != 1 {
		t.Errorf("%d", n)
	}
	if u.NoRequest != noRequest {
		t.Error("NoRequest Unmarshal error")
	}
	if u.NoResponse != noResponse {
		t.Error("NoResponse Unmarshal error")
	}
	if u.Compress != flateCompress {
		t.Errorf("Compress Unmarshal error %d", u.Compress)
	}
	if u.Heartbeat != heartbeat {
		t.Error("Heartbeat Unmarshal error")
	}
	if u.Watch != watch {
		t.Error("Watch Unmarshal error")
	}
	if u.Reserve != 1 {
		t.Errorf("Reserve Unmarshal error %d", u.Reserve)
	}
}

func TestUpgradeUnmarshal(t *testing.T) {
	u := &upgrade{}
	offset, err := u.Unmarshal(nil)
	if offset != 0 {
		t.Error(offset)
	}
	if err == nil {
		t.Error("should error")
	}
}
