// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/socket"
	"sync"
	"testing"
)

func TestNewServerCodec(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	wg := sync.WaitGroup{}
	wg.Add(1)
	sock, _ := socket.NewSocket(network, nil)
	lis, _ := sock.Listen(addr)
	go func() {
		defer wg.Done()
		conn, _ := lis.Accept()
		message := conn.Messages()
		headerEncoder := NewHeaderEncoder("json")
		if NewServerCodec(nil, nil, nil, false, bufferSize) != nil {
			t.Error("should be nil")
		}
		if NewServerCodec(nil, nil, message, false, bufferSize) != nil {
			t.Error("should be nil")
		}
		codec := NewServerCodec(nil, headerEncoder(), message, false, bufferSize)
		if codec == nil {
			t.Error("should not be nil")
		}
		if codec.ReadRequestBody(nil, nil) == nil {
			t.Error("The err should not be nil")
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ServeCodec(codec)
		}()
		message.Close()
		lis.Close()
		codec.ReadRequestHeader(&Context{})
		codec.Close()
		if codec.ReadRequestHeader(nil) == nil {
			t.Error("The err should not be nil")
		}
		if codec.WriteResponse(nil, nil) == nil {
			t.Error("The err should not be nil")
		}
	}()
	sock.Dial(addr)
	wg.Wait()
}

func TestNewServerCodecNoBatch(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	wg := sync.WaitGroup{}
	wg.Add(1)
	sock, _ := socket.NewSocket(network, nil)
	lis, _ := sock.Listen(addr)
	go func() {
		defer wg.Done()
		conn, _ := lis.Accept()
		message := conn.Messages()
		headerEncoder := NewHeaderEncoder("json")
		if NewServerCodec(nil, nil, nil, true, bufferSize) != nil {
			t.Error("should be nil")
		}
		if NewServerCodec(nil, nil, message, true, bufferSize) != nil {
			t.Error("should be nil")
		}
		codec := NewServerCodec(nil, headerEncoder(), message, true, bufferSize)
		if codec == nil {
			t.Error("should not be nil")
		}
		if codec.ReadRequestBody(nil, nil) == nil {
			t.Error("The err should not be nil")
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ServeCodec(codec)
		}()
		message.Close()
		lis.Close()
	}()
	sock.Dial(addr)
	wg.Wait()
}
