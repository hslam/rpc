package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"github.com/hslam/socket"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
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
	res = service.ArithResponse{}
	call := conn.Go("Arith.Multiply", req, &res, make(chan *Call, 1))
	<-call.Done
	if res.Pro != A*B {
		t.Error(res.Pro)
	}

	call = new(Call)
	call.ServiceMethod = "Arith.Multiply"
	call.Args = req
	call.Reply = &service.ArithResponse{}
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

func TestNewClientWithCodec(t *testing.T) {
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
	conn, err := sock.Dial(addr)
	if err != nil {
		t.Error(err)
	}
	clientCodec := NewClientCodec(NewJSONCodec(), NewCODEEncoder(), conn.Messages())
	client := NewClientWithCodec(clientCodec)
	if client == nil {
		t.Error("should not be nil")
	}
	err = client.Ping()
	if err != nil {
		t.Error(err)
	}
	client.Close()
	err = client.Ping()
	if err != ErrShutdown {
		t.Error(err)
	}
	server.Close()
	wg.Wait()
}
