// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	codeservice "github.com/hslam/rpc/examples/codec/code/service"
	"github.com/hslam/rpc/examples/codec/json/service"
	"github.com/hslam/socket"
	"sync"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
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

	valueCtx := context.WithValue(context.Background(), BufferContextKey, make([]byte, 64))
	ctx, cancel := context.WithTimeout(valueCtx, time.Minute)
	res = service.ArithResponse{}
	if err := conn.CallWithContext(ctx, "Arith.Multiply", req, &res); err != nil {
		t.Error(err)
	}
	if res.Pro != A*B {
		t.Error(res.Pro)
	}
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 0)
	res = service.ArithResponse{}
	if err := conn.CallWithContext(ctx, "Arith.Multiply", req, &res); err == nil {
		t.Error(err)
	}
	cancel()

	res = service.ArithResponse{}
	call := conn.Go("Arith.Multiply", req, &res, make(chan *Call, 1))
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	call = new(Call)
	call.ServiceMethod = "Arith.Multiply"
	call.Args = req
	call.Reply = &res
	conn.RoundTrip(call)
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}
	GetLogLevel()
	time.Sleep(time.Second)
	conn.Close()
	server.Close()
	wg.Wait()
}

func TestConnConcurrent(t *testing.T) {
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
	cwg := sync.WaitGroup{}
	for i := 0; i < 512; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			for j := 0; j < 64; j++ {
				A := int32(j + 4)
				B := int32(j + 8)
				req := &service.ArithRequest{A: A, B: B}
				var res service.ArithResponse
				if err := conn.Call("Arith.Multiply", req, &res); err != nil {
					t.Error(err)
				}
				if res.Pro != A*B {
					t.Error(res.Pro)
				}
			}
		}()
	}
	cwg.Wait()
	time.Sleep(time.Millisecond * 100)
	conn.Close()
	server.Close()
	wg.Wait()
}

func TestNewConnWithCodec(t *testing.T) {
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
	sock := socket.NewTCPSocket(nil)
	c, err := sock.Dial(addr)
	if err != nil {
		t.Error(err)
	}
	if NewConnWithCodec(nil) != nil {
		t.Error("should be nil")
	}
	clientCodec := NewClientCodec(NewJSONCodec(), DefaultEncoder(), c.Messages(), bufferSize)
	conn := NewConnWithCodec(clientCodec)
	if conn == nil {
		t.Error("should not be nil")
	}
	err = conn.Ping()
	if err != nil {
		t.Error(err)
	}
	conn.Close()
	err = conn.Ping()
	if err != ErrShutdown {
		t.Error(err)
	}
	pc := &persistConn{Conn: conn}
	checkPersistConnErr(ErrShutdown, pc)
	if pc.alive == true {
		t.Error("should not be alive")
	}
	server.Close()
	wg.Wait()
}

func TestConnWrite(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "code"
	server := NewServer()
	err := server.Register(new(codeservice.Arith))
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
	req := &codeservice.ArithRequest{A: A, B: B}
	var res codeservice.ArithResponse
	call := new(Call)
	call.ServiceMethod = "Arith.Multiply"
	{
		call.Args = nil
		call.Reply = &res
		conn.RoundTrip(call)
		<-call.Done
		if call.Error == nil {
			t.Error("The err should not be nil")
		}
	}
	{
		call.Args = req
		call.Reply = nil
		conn.RoundTrip(call)
		<-call.Done
		if call.Error == nil {
			t.Error("The err should not be nil")
		}
	}
	conn.Close()
	server.Close()
	wg.Wait()
}

func TestServerCodecReadRequestBody(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "code"
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
	req := &codeservice.ArithRequest{A: A, B: B}
	var res codeservice.ArithResponse
	call := new(Call)
	call.ServiceMethod = "Arith.Multiply"
	call.Args = req
	call.Reply = &res
	conn.RoundTrip(call)
	<-call.Done
	if call.Error == nil {
		t.Error("The err should not be nil")
	}
	conn.Close()
	server.Close()
	wg.Wait()
}

func TestConnWriteWatch(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "code"
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
	conn.Close()
	conn.shutdown = false
	conn.closing = false
	server.Close()
	wg.Wait()
}

type arithMultiply struct{}

func (a *arithMultiply) Multiply(req *codeservice.ArithRequest, res *service.ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

func (a *arithMultiply) MultiplyReturnOut(req *codeservice.ArithRequest) (*codeservice.ArithResponse, error) {
	var res codeservice.ArithResponse
	res.Pro = req.A * req.B
	return &res, nil
}

func (a *arithMultiply) MultiplyReturnOutWithContext(ctx context.Context, req *codeservice.ArithRequest) (*codeservice.ArithResponse, error) {
	var res codeservice.ArithResponse
	res.Pro = req.A * req.B
	return &res, nil
}

func TestServerCodecWriteResponse(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "code"
	server := NewServer()
	err := server.RegisterName("B", new(arithMultiply))
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
		call := new(Call)
		call.ServiceMethod = "B.Multiply"
		call.Args = &codeservice.ArithRequest{A: A, B: B}
		res := codeservice.ArithResponse{}
		call.Reply = &res
		conn.RoundTrip(call)
		<-call.Done
		if call.Error == nil {
			t.Error("The err should not be nil")
		}
	}
	{
		A := int32(4)
		B := int32(8)
		call := new(Call)
		call.ServiceMethod = "B.MultiplyReturnOut"
		call.Args = &codeservice.ArithRequest{A: A, B: B}
		res := codeservice.ArithResponse{}
		call.Reply = &res
		conn.RoundTrip(call)
		<-call.Done
		if call.Error != nil {
			t.Error(call.Error.Error())
		} else if res.Pro != A*B {
			t.Errorf("%d * %d, pro is %d\n", A, B, call.Reply)
		}
	}
	{
		A := int32(4)
		B := int32(8)
		call := new(Call)
		call.ServiceMethod = "B.MultiplyReturnOutWithContext"
		call.Args = &codeservice.ArithRequest{A: A, B: B}
		res := codeservice.ArithResponse{}
		call.Reply = &res
		conn.RoundTrip(call)
		<-call.Done
		if call.Error != nil {
			t.Error(call.Error.Error())
		} else if res.Pro != A*B {
			t.Errorf("%d * %d, pro is %d\n", A, B, call.Reply)
		}
	}
	conn.Close()
	server.Close()
	wg.Wait()
}

func TestConnRead(t *testing.T) {
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
	req.A = 300
	conn, err := Dial(network, addr, codec)
	if err != nil {
		t.Error(err)
	}
	conn.Go("Arith.Sleep", req, &res, nil)
	time.Sleep(time.Millisecond * 10)
	for seq := range conn.pending {
		delete(conn.pending, seq)
	}
	time.Sleep(time.Millisecond * 500)
	conn.Close()
	server.Close()
	wg.Wait()
}

func TestCallDone(t *testing.T) {
	call := &Call{Done: make(chan *Call, 1)}
	call.done()
	call.done()
}

func TestResetDone(t *testing.T) {
	done := make(chan *Call, 10)
	call := &Call{Done: done}
	call.done()
	call.done()
	ResetDone(done)
	if len(done) != 0 {
		t.Error(len(done))
	}
	onceDone(done)
}

func TestCheckDone(t *testing.T) {
	done := checkDone(nil)
	if done == nil {
		t.Error("should not be nil")
	}
	done = make(chan *Call)
	defer func() {
		e := recover()
		if e == nil {
			t.Error("should panix")
		}
	}()
	checkDone(done)
}
