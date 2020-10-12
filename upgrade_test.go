// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"testing"
)

func TestUpgrade(t *testing.T) {
	u := &upgrade{
		NoRequest:  noRequest,
		NoResponse: noResponse,
		Compress:   flateCompress,
		Heartbeat:  heartbeat,
		Wait:       wait,
		Watch:      watch,
		Reserve:    1,
	}
	buf := make([]byte, 64)
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
	if u.Wait != wait {
		t.Error("Heartbeat Unmarshal error")
	}
	if u.Watch != watch {
		t.Error("Heartbeat Unmarshal error")
	}
	if u.Reserve != 1 {
		t.Errorf("Reserve Unmarshal error %d", u.Reserve)
	}
}
