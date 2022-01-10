// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	addr := ":8880"
	network := ""
	codec := ""
	if _, err := Dial(network, addr, codec); err == nil {
		t.Error("The err should not be nil")
	}
	network = "tcp"
	if _, err := Dial(network, addr, codec); err == nil {
		t.Error("The err should not be nil")
	}
}

func TestDialTLS(t *testing.T) {
	addr := ":8880"
	network := ""
	codec := ""
	if _, err := DialTLS(network, addr, codec, SkipVerifyTLSConfig()); err == nil {
		t.Error("The err should not be nil")
	}
	network = "tcp"
	if _, err := DialTLS(network, addr, codec, SkipVerifyTLSConfig()); err == nil {
		t.Error("The err should not be nil")
	}
}

func TestDialWithOptions(t *testing.T) {
	addr := ":9999"
	opts := &Options{}
	opts.Network = "tcp"
	opts.Codec = "json"
	err := Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ListenWithOptions(addr, opts)
	}()
	time.Sleep(time.Millisecond * 10)
	opts.Network = ""
	opts.Codec = ""
	if _, err := DialWithOptions(addr, opts); err == nil {
		t.Error("The err should not be nil")
	}
	opts.Codec = "json"
	if _, err := DialWithOptions(addr, opts); err == nil {
		t.Error("The err should not be nil")
	}
	opts.Codec = ""
	opts.NewCodec = NewCodec("json")
	opts.NewSocket = NewSocket("tcp")
	if conn, err := DialWithOptions(addr, opts); err != nil {
		t.Error(err)
	} else {
		conn.Close()
	}
	if conn, err := DialWithOptions(addr, opts); err != nil {
		t.Error(err)
	} else {
		conn.Close()
	}
	opts.NewHeaderEncoder = NewHeaderEncoder("json")
	if conn, err := DialWithOptions(addr, opts); err != nil {
		t.Error(err)
	} else {
		conn.Close()
	}
	DefaultServer.Close()
	wg.Wait()
}
