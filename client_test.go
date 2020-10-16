package rpc

import (
	codeservice "github.com/hslam/rpc/examples/codec/code/service"
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
	if NewClientWithCodec(nil) != nil {
		t.Error("should be nil")
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
	pc := &persistConn{Client: client}
	checkPersistConnErr(ErrShutdown, pc)
	if pc.alive == true {
		t.Error("should not be alive")
	}
	server.Close()
	wg.Wait()
}

func TestClientWrite(t *testing.T) {
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

type arithMultiply struct{}

func (a *arithMultiply) Multiply(req *codeservice.ArithRequest, res *service.ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
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
	A := int32(4)
	B := int32(8)
	call := new(Call)
	call.ServiceMethod = "B.Multiply"
	call.Args = &codeservice.ArithRequest{A: A, B: B}
	call.Reply = nil
	conn.RoundTrip(call)
	<-call.Done
	if call.Error == nil {
		t.Error("The err should not be nil")
	}
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
