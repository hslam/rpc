package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestRPC(t *testing.T) {
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
