// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestReverseProxy(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	opts := DefaultOptions()
	opts.Network = network
	opts.Codec = codec
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
	trans := &Transport{
		MaxConnsPerHost:     1,
		MaxIdleConnsPerHost: 1,
		KeepAlive:           time.Second * 60,
		IdleConnTimeout:     time.Second * 60,
		Options:             opts,
	}
	proxy := NewReverseProxy(addr)
	proxy.Transport = trans
	err = proxy.Ping()
	if err != nil {
		t.Error(err)
	}

	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	if err := proxy.Call("Arith.Multiply", req, &res); err != nil {
		t.Error(err)
	}
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	res = service.ArithResponse{}
	call := proxy.Go("Arith.Multiply", req, &res, make(chan *Call, 1))
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	call = new(Call)
	call.ServiceMethod = "Arith.Multiply"
	call.Args = req
	call.Reply = &service.ArithResponse{}
	proxy.RoundTrip(call)
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	trans.Close()
	server.Close()
	wg.Wait()
}

func TestReverseProxyLeastTime(t *testing.T) {
	network := "tcp"
	addrs := []string{":9997", ":9998", ":9999"}
	codec := "json"
	opts := DefaultOptions()
	opts.Network = network
	opts.Codec = codec
	server := NewServer()
	err := server.Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < len(addrs); i++ {
		wg.Add(1)
		addr := addrs[i]
		go func() {
			defer wg.Done()
			server.ListenWithOptions(addr, opts)
		}()
	}
	time.Sleep(time.Millisecond * 10)
	trans := &Transport{
		MaxConnsPerHost:     1,
		MaxIdleConnsPerHost: 1,
		KeepAlive:           time.Second * 60,
		IdleConnTimeout:     time.Second * 60,
		Options:             opts,
	}
	proxy := NewReverseProxy(addrs...)
	proxy.Scheduling = LeastTimeScheduling
	proxy.Transport = trans
	err = proxy.Ping()
	if err != nil {
		t.Error(err)
	}
	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	if err := proxy.Call("Arith.Multiply", req, &res); err != nil {
		t.Error(err)
	}
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	res = service.ArithResponse{}
	call := proxy.Go("Arith.Multiply", req, &res, make(chan *Call, 1))
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	call = new(Call)
	call.ServiceMethod = "Arith.Multiply"
	call.Args = req
	call.Reply = &service.ArithResponse{}
	proxy.RoundTrip(call)
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	trans.Close()
	server.Close()
	wg.Wait()
}

func TestReverseProxyEmptyTargets(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("should panic")
		}
	}()
	NewReverseProxy()
}

func TestReverseProxyTransport(t *testing.T) {
	addr := ":9999"
	proxy := NewReverseProxy(addr)
	defer func() {
		if e := recover(); e == nil {
			t.Error("should panic")
		}
	}()
	proxy.transport()
}

func TestReverseProxyTarget(t *testing.T) {
	proxy := NewReverseProxy(":9999", ":9998", ":9997")
	proxy.Scheduling = RoundRobinScheduling
	if address, target := proxy.target(); len(address) == 0 && target == nil {
		t.Error()
	}
	proxy.Scheduling = RandomScheduling
	if address, target := proxy.target(); len(address) == 0 && target == nil {
		t.Error()
	}
	proxy.Scheduling = LeastTimeScheduling
	if address, target := proxy.target(); len(address) == 0 && target == nil {
		t.Error()
	}
	proxy.Scheduling = 3
	if address, target := proxy.target(); len(address) == 0 && target == nil {
		t.Error()
	}
	proxy.Director = func() (target string) {
		return ":9999"
	}
	if address, target := proxy.target(); len(address) == 0 && target == nil {
		t.Error()
	}
	proxy = &ReverseProxy{}
	defer func() {
		if e := recover(); e == nil {
			t.Error("should panic")
		}
	}()
	proxy.target()
}

func TestTopK(t *testing.T) {
	l := list{&target{latency: 10}, &target{latency: 7}, &target{latency: 2}, &target{latency: 5}, &target{latency: 1}, &target{latency: 6}}
	minHeap(l)
	n := l.Len()
	for i := 1; i < n; i++ {
		if l.Less(i, 0) {
			t.Error("heap error")
		}
	}
}
