// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	network := "tcp"
	addr := ":8880"
	codec := "json"
	SetBufferSize(0)
	err := Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	services := Services()
	if len(services) != 1 && services[0] != "Arith.Multiply" {
		t.Error()
	}
	SetLogLevel(GetLogLevel())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		Listen(network, addr, codec)
	}()
	time.Sleep(time.Millisecond * 10)
	conn, err := Dial(network, addr, codec)
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

	err = Listen(network, addr, codec)
	if err == nil {
		t.Error("The err should be address already in use")
	}
	DefaultServer.Close()
	wg.Wait()
}

func TestListenTLS(t *testing.T) {
	network := "tcp"
	addr := ":8880"
	codec := "json"
	err := Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ListenTLS(network, addr, codec, DefalutServerTLSConfig())
	}()
	time.Sleep(time.Millisecond * 10)
	conn, err := DialTLS(network, addr, codec, DefalutClientTLSConfig())
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
	DefaultServer.Close()
	wg.Wait()
}

func TestListenWithOptions(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	opts := DefaultOptions()
	opts.Network = network
	opts.Codec = codec
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
	conn, err := DialWithOptions(addr, opts)
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
	DefaultServer.Close()
	wg.Wait()
}

func TestServerListen(t *testing.T) {
	addr := ":8880"
	network := ""
	codec := ""

	if err := DefaultServer.Listen(network, addr, codec); err == nil {
		t.Error("The err should not be nil")
	}
	network = "tcp"
	if err := DefaultServer.Listen(network, addr, codec); err == nil {
		t.Error("The err should not be nil")
	}
}

func TestServerListenTLS(t *testing.T) {
	addr := ":8880"
	network := ""
	codec := ""
	if err := DefaultServer.ListenTLS(network, addr, codec, SkipVerifyTLSConfig()); err == nil {
		t.Error("The err should not be nil")
	}
	network = "tcp"
	if err := DefaultServer.ListenTLS(network, addr, codec, SkipVerifyTLSConfig()); err == nil {
		t.Error("The err should not be nil")
	}
}

func TestServerListenWithOptions(t *testing.T) {
	addr := ":9999"
	opts := &Options{}
	opts.Network = "tcp"
	opts.Codec = "json"
	opts.Network = ""
	opts.Codec = ""
	{
		server := NewServer()
		if err := server.ListenWithOptions(addr, opts); err == nil {
			t.Error("The err should not be nil")
		}
		opts.Codec = "json"
		if err := server.ListenWithOptions(addr, opts); err == nil {
			t.Error("The err should not be nil")
		}
	}
	wg := sync.WaitGroup{}
	{
		server := NewServer()
		opts.NewSocket = NewSocket("tcp")
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.ListenWithOptions(addr, opts)
		}()
		time.Sleep(time.Millisecond * 10)
		server.Close()
		wg.Wait()
	}
	{
		server := NewServer()
		opts.Codec = ""
		opts.HeaderEncoder = ""
		opts.NewCodec = nil
		opts.NewHeaderEncoder = NewHeaderEncoder("json")
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.ListenWithOptions(addr, opts)
		}()
		time.Sleep(time.Millisecond * 10)
		server.Close()
		wg.Wait()
	}
}
