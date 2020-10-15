package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestServerPoll(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	server.SetPoll(true)
	err := server.Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Listen(network, addr, codec)
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
	server.Close()
	wg.Wait()
}

func TestServerPipelining(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	server.SetPipelining(true)
	err := server.Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Listen(network, addr, codec)
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
	server.Close()
	wg.Wait()
}

func TestServerPollPipelining(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	SetPoll(true)
	SetPipelining(true)
	err := RegisterName("Arith", new(service.Arith))
	if err != nil {
		t.Error(err)
	}
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
	DefaultServer.Close()
	wg.Wait()
}

type arith struct{}

func (a *arith) Add() error {
	return nil
}

func (a *arith) Divide(req *service.ArithRequest) error {
	return nil
}

func TestFunc(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	err := server.RegisterName("Arith", new(arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Listen(network, addr, codec)
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
	conn.Close()
	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	{
		conn, err := Dial(network, addr, codec)
		if err != nil {
			t.Error(err)
		}
		if err := conn.Call("Arith.Multiply", req, &res); err == nil {
			t.Error("error")
		}
		conn.Close()
	}
	{
		conn, err := Dial(network, addr, codec)
		if err != nil {
			t.Error(err)
		}
		if err := conn.Call("Arith.Add", req, &res); err == nil {
			t.Error("error")
		}
		conn.Close()
	}
	{
		conn, err := Dial(network, addr, codec)
		if err != nil {
			t.Error(err)
		}
		if err := conn.Call("Arith.Divide", req, &res); err == nil {
			t.Error("error")
		}
		conn.Close()
	}
	server.Close()
	wg.Wait()
}

func TestServerPush(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	var lock sync.Mutex
	var k = "foo"
	var str = "bar"
	var v = []byte(str)
	PushFunc(func(key string) (value []byte, ok bool) {
		if key == k {
			lock.Lock()
			value = v
			lock.Unlock()
			return value, true
		}
		return nil, false
	})
	wg := sync.WaitGroup{}
	done := make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			select {
			case <-ticker.C:
				lock.Lock()
				v = []byte(str)
				value := v
				lock.Unlock()
				Push(k, value)
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
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
	watch := conn.Watch(k)
	for i := 0; i < 3; i++ {
		value, err := watch.Wait()
		if err != nil {
			panic(err)
		}
		if string(value) != str {
			t.Errorf("Watch foo:%s\n", string(value))
		}
	}
	watch.Stop()
	PushFunc(nil)
	watch = conn.Watch(k)
	for i := 0; i < 3; i++ {
		value, err := watch.Wait()
		if err != nil {
			panic(err)
		}
		if string(value) != str {
			t.Errorf("Watch foo:%s\n", string(value))
		}
	}
	conn.Close()
	time.Sleep(time.Millisecond * 100)
	close(done)
	DefaultServer.Close()
	wg.Wait()
}
