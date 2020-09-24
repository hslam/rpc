// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestHeaderEncoder(t *testing.T) {
	encoder := DefaultEncoder()
	if encoder == nil {
		t.Error("encoder is nil")
	}
	testHeaderEncoder("json", t)
	testHeaderEncoder("code", t)
	testHeaderEncoder("pb", t)
}

func testHeaderEncoder(headerEncoder string, t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	opts := DefaultOptions()
	opts.Network = network
	opts.Codec = codec
	opts.HeaderEncoder = headerEncoder
	server := NewServer()
	err := server.Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.ListenWithOptions(addr, opts)
	}()
	time.Sleep(time.Millisecond * 10)
	conn, err := DialWithOptions(addr, opts)
	if err != nil {
		t.Error(err)
	}
	err = conn.Ping()
	if err != nil {
		t.Error(err)
	}
	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	if err := conn.Call("Arith.Multiply", req, &res); err != nil {
		t.Error(err)
	}
	if res.Pro != A*B {
		t.Error(res.Pro)
	}
	conn.Close()
	server.Close()
	wg.Wait()
}
