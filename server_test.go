// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"errors"
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

func TestServerWorkers(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
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

func TestServerSetBufferSize(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	server.SetBufferSize(bufferSize)
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

func TestServerDirectIO(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	server.SetDirectIO(true)
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
	conn.SetDirectIO(true)
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
	conn.SetPipelining(true)
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
	SetDirectIO(false)
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
	connwg := sync.WaitGroup{}
	for i := 0; i < 32; i++ {
		connwg.Add(1)
		go func() {
			defer connwg.Done()
			conn, err := Dial(network, addr, codec)
			if err != nil {
				t.Error(err)
			}
			err = conn.Ping()
			if err != nil {
				t.Error(err)
			}
			cwg := sync.WaitGroup{}
			for i := 0; i < 32; i++ {
				cwg.Add(1)
				go func() {
					defer cwg.Done()
					for i := 0; i < 128; i++ {
						A := int32(4)
						B := int32(8)
						req := &service.ArithRequest{A: A, B: B}
						var res service.ArithResponse
						if err := conn.Call("Arith.Multiply", req, &res); err != nil {
							t.Error(err)
						}
					}
				}()
			}
			cwg.Wait()
			conn.Close()
		}()
	}
	connwg.Wait()
	DefaultServer.Close()
	wg.Wait()
	SetPoll(false)
	SetPipelining(false)
	SetDirectIO(false)
}

func TestServerPollPipeliningNoBatch(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	SetPoll(true)
	SetDirectIO(true)
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
	connwg := sync.WaitGroup{}
	for i := 0; i < 32; i++ {
		connwg.Add(1)
		go func() {
			defer connwg.Done()
			conn, err := Dial(network, addr, codec)
			if err != nil {
				t.Error(err)
			}
			err = conn.Ping()
			if err != nil {
				t.Error(err)
			}
			cwg := sync.WaitGroup{}
			for i := 0; i < 32; i++ {
				cwg.Add(1)
				go func() {
					defer cwg.Done()
					for i := 0; i < 128; i++ {
						A := int32(4)
						B := int32(8)
						req := &service.ArithRequest{A: A, B: B}
						var res service.ArithResponse
						if err := conn.Call("Arith.Multiply", req, &res); err != nil {
							t.Error(err)
						}
					}
				}()
			}
			cwg.Wait()
			conn.Close()
		}()
	}
	connwg.Wait()
	DefaultServer.Close()
	wg.Wait()
	SetPoll(false)
	SetPipelining(false)
	SetDirectIO(false)
}

type arithRequest struct {
	A int32
	B int32
}

type arithResponse struct {
	C int32
}

type arith struct{}

func (a *arith) Load() error {
	return nil
}

func (a *arith) Add(req *arithRequest) error {
	return nil
}

func (a *arith) Divide(req *arithRequest, res *arithResponse) error {
	if req.B == 0 {
		return errors.New("B can not be 0")
	}
	res.C = req.A / req.B
	return nil
}

func (a *arith) DivideWithContext(ctx context.Context, req *arithRequest, res *arithResponse) error {
	if req.B == 0 {
		return errors.New("B can not be 0")
	}
	res.C = req.A / req.B
	FreeContextBuffer(ctx)
	return nil
}

func (a *arith) Sleep(req *arithRequest, res *arithResponse) error {
	time.Sleep(time.Millisecond * time.Duration(req.A))
	return nil
}

func (a *arith) Stop(req *arithRequest, res *arithResponse) error {
	return ErrShutdown
}

func TestSetContextBuffer(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	SetContextBuffer(true)
	server.SetContextBuffer(true)
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
	{
		A := int32(4)
		B := int32(8)
		req := &arithRequest{A: A, B: B}
		var res arithResponse
		conn, err := Dial(network, addr, codec)
		if err != nil {
			t.Error(err)
		}
		ctx := context.WithValue(context.Background(), BufferContextKey, make([]byte, 64))
		if err := conn.CallWithContext(ctx, "Arith.DivideWithContext", req, &res); err != nil {
			t.Error(err)
		}
	}
	conn.Close()
	server.Close()
	wg.Wait()
}

func TestSetNoCopy(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	SetNoCopy(true)
	server.SetNoCopy(true)
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
	{
		A := int32(4)
		B := int32(8)
		req := &arithRequest{A: A, B: B}
		var res arithResponse
		conn, err := Dial(network, addr, codec)
		if err != nil {
			t.Error(err)
		}
		ctx := context.WithValue(context.Background(), BufferContextKey, make([]byte, 64))
		if err := conn.CallWithContext(ctx, "Arith.DivideWithContext", req, &res); err != nil {
			t.Error(err)
		}
	}
	conn.Close()
	server.Close()
	wg.Wait()
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
	{
		A := int32(4)
		B := int32(8)
		req := &arithRequest{A: A, B: B}
		var res arithResponse
		conn, err := Dial(network, addr, codec)
		if err != nil {
			t.Error(err)
		}
		if err := conn.CallWithContext(context.Background(), "Arith.DivideWithContext", req, &res); err != nil {
			t.Error(err)
		}
	}
	conn.Close()
	A := int32(4)
	B := int32(8)
	req := &arithRequest{A: A, B: B}
	var res arithResponse
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
		if err := conn.Call("Arith.Load", req, &res); err == nil {
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
		req.B = 0
		if err := conn.Call("Arith.Divide", req, &res); err == nil {
			t.Error("error")
		}
		conn.Close()
	}
	server.Close()
	wg.Wait()
}

func TestServerClose(t *testing.T) {
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
	A := int32(4)
	B := int32(8)
	req := &arithRequest{A: A, B: B}
	var res arithResponse
	req.A = 100
	conn, err := Dial(network, addr, codec)
	if err != nil {
		t.Error(err)
	}
	call := conn.Go("Arith.Sleep", req, &res, nil)
	time.Sleep(time.Millisecond * 10)
	conn.Close()
	<-call.Done
	time.Sleep(time.Millisecond * 500)
	server.Close()
	wg.Wait()
}

func TestServerStopClient(t *testing.T) {
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
	A := int32(4)
	B := int32(8)
	req := &arithRequest{A: A, B: B}
	var res arithResponse
	req.A = 100
	conn, err := Dial(network, addr, codec)
	if err != nil {
		t.Error(err)
	}
	call := conn.Go("Arith.Stop", req, &res, nil)
	time.Sleep(time.Millisecond * 10)
	conn.Close()
	<-call.Done
	time.Sleep(time.Millisecond * 500)
	server.Close()
	wg.Wait()
}
