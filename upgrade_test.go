// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"testing"
)

func TestUpgrade(t *testing.T) {
	u := &upgrade{
		Heartbeat:  Heartbeat,
		NoRequest:  NoRequest,
		NoResponse: NoResponse,
		Compress:   GzipCompress,
		Reserve:    1,
	}
	buf := make([]byte, 64)
	data, err := u.Marshal(buf)
	if err != nil {
		t.Error(err)
	}
	*u = upgrade{}
	n, err := u.Unmarshal(data)
	if err != nil {
		t.Error(err)
	}
	if n != 1 {
		t.Errorf("%d", n)
	}
	if u.Heartbeat != Heartbeat {
		t.Error("Heartbeat Unmarshal error")
	}
	if u.NoRequest != NoRequest {
		t.Error("NoRequest Unmarshal error")
	}
	if u.NoResponse != NoResponse {
		t.Error("NoResponse Unmarshal error")
	}
	if u.Compress != GzipCompress {
		t.Errorf("Compress Unmarshal error %d", u.Compress)
	}
	if u.Reserve != 1 {
		t.Errorf("Reserve Unmarshal error %d", u.Reserve)
	}
}
