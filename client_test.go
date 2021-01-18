// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	var k = "foo"
	var str = "bar"
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
	server.PushFunc(func(key string) (value []byte, ok bool) {
		if key == k {
			return []byte(str), true
		}
		return nil, false
	})
	client := NewClient(opts, addr)
	err = client.Ping()
	if err != nil {
		t.Error(err)
	}
	watch, err := client.Watch(k)
	if err != nil {
		t.Error(err)
	}
	v, err := watch.Wait()
	if err != nil {
		t.Error(err)
	} else if string(v) != str {
		t.Error(string(v))
	}
	watch.Stop()
	if _, err := watch.Wait(); err == nil {
		t.Error()
	}
	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	if err := client.Call("Arith.Multiply", req, &res); err != nil {
		t.Error(err)
	}
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	res = service.ArithResponse{}
	if err := client.CallWithContext(ctx, "Arith.Multiply", req, &res); err != nil {
		t.Error(err)
	}
	if res.Pro != A*B {
		t.Error(res.Pro)
	}
	cancel()

	res = service.ArithResponse{}
	call := client.Go("Arith.Multiply", req, &res, make(chan *Call, 1))
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	call = new(Call)
	call.ServiceMethod = "Arith.Multiply"
	call.Args = req
	call.Reply = &service.ArithResponse{}
	client.RoundTrip(call)
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}
	client.DialTimeout = time.Minute
	cwg := sync.WaitGroup{}
	cwg.Add(1)
	go func() {
		defer cwg.Done()
		client.Ping()
	}()
	time.Sleep(time.Millisecond * 100)
	client.Close()
	cwg.Wait()
	server.Close()
	wg.Wait()
	{
		client := NewClient(opts, addr)
		client.DialTimeout = time.Millisecond * 100
		err = client.Ping()
		if err == nil {
			t.Error()
		}
		_, err := client.Watch(k)
		if err == nil {
			t.Error()
		}
		A := int32(4)
		B := int32(8)
		req := &service.ArithRequest{A: A, B: B}
		var res service.ArithResponse
		if err := client.Call("Arith.Multiply", req, &res); err == nil {
			t.Error()
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		res = service.ArithResponse{}
		if err := client.CallWithContext(ctx, "Arith.Multiply", req, &res); err == nil {
			t.Error()
		}
		cancel()

		res = service.ArithResponse{}
		call := client.Go("Arith.Multiply", req, &res, make(chan *Call, 1))
		<-call.Done
		if call.Error == nil {
			t.Error()
		}

		call = new(Call)
		call.ServiceMethod = "Arith.Multiply"
		call.Args = req
		call.Reply = &service.ArithResponse{}
		client.RoundTrip(call)
		<-call.Done
		if call.Error == nil {
			t.Error()
		}
	}
}

func TestClientFallback(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	opts := DefaultOptions()
	opts.Network = network
	opts.Codec = codec
	{
		client := NewClient(opts, addr)
		client.Fallback(time.Millisecond * 10)
		time.Sleep(time.Millisecond * 100)
		client.Close()
	}
	{
		client := NewClient(opts, addr)
		client.Fallback(time.Millisecond * 10)
		client.Close()
		time.Sleep(time.Millisecond * 100)
	}
}

func TestClientClose(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	opts := DefaultOptions()
	opts.Network = network
	opts.Codec = codec
	client := NewClient(opts, addr)
	cwg := sync.WaitGroup{}
	cwg.Add(1)
	go func() {
		defer cwg.Done()
		client.Ping()
	}()
	time.Sleep(time.Millisecond * 10)
	client.Close()
	cwg.Wait()
	cwg.Add(1)
	go func() {
		defer cwg.Done()
		client.Ping()
	}()
	time.Sleep(time.Millisecond * 10)
	cwg.Wait()
	{
		c := NewClient(opts, addr)
		c.Close()
		done := c.donePool.Get().(chan *waiter)
		w := c.waiterPool.Get().(*waiter)
		w.Done = done
		c.wait(w)
	}
}

func TestClientLeastTime(t *testing.T) {
	network := "tcp"
	addrs := []string{":9997", ":9998", ":9999"}
	codec := "json"
	var k = "foo"
	var str = "bar"
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
	server.PushFunc(func(key string) (value []byte, ok bool) {
		if key == k {
			return []byte(str), true
		}
		return nil, false
	})
	client := NewClient(opts, addrs...)
	client.Scheduling = LeastTimeScheduling
	for i := 0; i < 64; i++ {
		err = client.Ping()
		if err != nil {
			t.Error(err)
		}
		watch, err := client.Watch(k)
		if err != nil {
			t.Error(err)
		}
		v, err := watch.Wait()
		if err != nil {
			t.Error(err)
		} else if string(v) != str {
			t.Error(string(v))
		}
		watch.Stop()
		if _, err := watch.Wait(); err == nil {
			t.Error()
		}
		A := int32(4)
		B := int32(8)
		req := &service.ArithRequest{A: A, B: B}
		var res service.ArithResponse
		if err := client.Call("Arith.Multiply", req, &res); err != nil {
			t.Error(err)
		}
		if res.Pro != A*B {
			t.Error(res.Pro)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		res = service.ArithResponse{}
		if err := client.CallWithContext(ctx, "Arith.Multiply", req, &res); err != nil {
			t.Error(err)
		}
		if res.Pro != A*B {
			t.Error(res.Pro)
		}
		cancel()

		res = service.ArithResponse{}
		call := client.Go("Arith.Multiply", req, &res, make(chan *Call, 1))
		<-call.Done
		if res.Pro != A*B {
			t.Error(res.Pro)
		}

		call = new(Call)
		call.ServiceMethod = "Arith.Multiply"
		call.Args = req
		call.Reply = &service.ArithResponse{}
		client.RoundTrip(call)
		<-call.Done
		if res.Pro != A*B {
			t.Error(res.Pro)
		}
	}
	client.Close()
	server.Close()
	wg.Wait()
}

func TestClientOptions(t *testing.T) {
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error()
			}
		}()
		NewClient(&Options{})
	}()
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error()
			}
		}()
		NewClient(&Options{Codec: "json"})
	}()
}

func TestClientEmptyTargets(t *testing.T) {
	client := NewClient(nil)
	client.Director = func() (target string) {
		return ":9999"
	}
	if address, target, err := client.director(); len(address) == 0 && target == nil && err == nil {
		t.Error()
	}
}

func TestClientTransport(t *testing.T) {
	addr := ":9999"
	client := NewClient(nil, addr)
	defer func() {
		if e := recover(); e == nil {
			t.Error("should panic")
		}
	}()
	client.transport()
}

func TestClientUpdate(t *testing.T) {
	addr := ":9999"
	client := NewClient(nil, addr)
	tr := &target{address: addr}
	tr.Update(client.Alpha, 0, ErrDial)
}

func TestClientTarget(t *testing.T) {
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
	client := NewClient(opts, addrs...)
	for i := 0; i < 64; i++ {
		client.Ping()
	}
	time.Sleep(time.Millisecond * 100)
	client.Scheduling = RoundRobinScheduling
	if address, target, err := client.director(); len(address) == 0 && target == nil && err == nil {
		t.Error()
	}
	client.Scheduling = RandomScheduling
	if address, target, err := client.director(); len(address) == 0 && target == nil && err == nil {
		t.Error()
	}
	client.Scheduling = LeastTimeScheduling
	if address, target, err := client.director(); len(address) == 0 && target == nil && err == nil {
		t.Error()
	}
	client.Scheduling = 3
	if address, target, err := client.director(); len(address) == 0 && target == nil && err == nil {
		t.Error()
	}
	client.Director = func() (target string) {
		return addrs[0]
	}
	if address, target, err := client.director(); len(address) == 0 && target == nil && err == nil {
		t.Error()
	}
	client = &Client{}
	if _, _, err := client.schedule(); err == nil {
		t.Error()
	}
	server.Close()
	wg.Wait()
}

func TestTopK(t *testing.T) {
	h := list{&target{latency: 10}, &target{latency: 7}, &target{latency: 2}, &target{latency: 5}, &target{latency: 1}, &target{latency: 6}}
	minHeap(h)
	n := h.Len()
	for i := 1; i < n; i++ {
		if h.Less(i, 0) {
			t.Error("heap error")
		}
	}
}

func TestWaiterDone(t *testing.T) {
	w := &waiter{Done: make(chan *waiter, 1)}
	w.done()
	w.done()
}

func TestResetWaiterDone(t *testing.T) {
	done := make(chan *waiter, 10)
	w := &waiter{Done: done}
	w.done()
	w.done()
	resetWaiterDone(done)
	if len(done) != 0 {
		t.Error(len(done))
	}
	onceWaiterDone(done)
}
