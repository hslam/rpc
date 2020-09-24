package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"github.com/hslam/socket"
	"sync"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
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
		ListenTLS(network, addr, codec, socket.DefalutTLSConfig())
	}()
	time.Sleep(time.Millisecond * 10)
	conn, err := DialTLS(network, addr, codec, socket.SkipVerifyTLSConfig())
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
