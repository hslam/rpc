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
		Heartbeat:  heartbeat,
		Stream:     streaming,
	}
	if u.IsZero() {
		t.Error()
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
	if u.Heartbeat != heartbeat {
		t.Error("Heartbeat Unmarshal error")
	}
	if u.Stream != streaming {
		t.Error("Watch Unmarshal error")
	}
}

func TestUpgradeUnmarshal(t *testing.T) {
	u := &upgrade{}
	offset, err := u.Unmarshal(nil)
	if offset != 0 {
		t.Error(offset)
	}
	if err == nil {
		t.Error("The err should not be nil")
	}
}
