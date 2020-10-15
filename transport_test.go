package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"sync"
	"testing"
	"time"
)

func TestTransport(t *testing.T) {
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
		MaxConnsPerHost:     64,
		MaxIdleConnsPerHost: 32,
		KeepAlive:           time.Millisecond * 100,
		IdleConnTimeout:     time.Second,
		Options:             opts,
		ticker:              time.Millisecond * 10,
	}
	err = trans.Ping(addr)
	if err != nil {
		t.Error(err)
	}
	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	cwg := sync.WaitGroup{}
	for i := 0; i < 1024; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			var res service.ArithResponse
			if err := trans.Call(addr, "Arith.Multiply", req, &res); err != nil {
				t.Error(err)
			}
			if res.Pro != A*B {
				t.Error(res.Pro)
			}

			res = service.ArithResponse{}
			call := trans.Go(addr, "Arith.Multiply", req, &res, make(chan *Call, 1))
			<-call.Done
			if res.Pro != A*B {
				t.Error(res.Pro)
			}

			call = new(Call)
			call.ServiceMethod = "Arith.Multiply"
			call.Args = req
			call.Reply = &service.ArithResponse{}
			trans.RoundTrip(addr, call)
			<-call.Done
			if res.Pro != A*B {
				t.Error(res.Pro)
			}
		}()
	}
	cwg.Wait()
	time.Sleep(time.Millisecond * 100)
	trans.CloseIdleConnections()
	time.Sleep(time.Millisecond * 100)
	trans.Close()
	server.Close()
	err = trans.Ping(addr)
	if err == nil {
		t.Error("should be error")
	}
	res := service.ArithResponse{}
	call := trans.Go(addr, "Arith.Multiply", req, &res, make(chan *Call, 1))
	<-call.Done
	if call.Error == nil {
		t.Error("should be error")
	}
	if err := trans.Call(addr, "Arith.Multiply", req, &res); err == nil {
		t.Error("should be error")
	}
	call = new(Call)
	call.ServiceMethod = "Arith.Multiply"
	call.Args = req
	call.Reply = &service.ArithResponse{}
	trans.RoundTrip(addr, call)
	<-call.Done
	if call.Error == nil {
		t.Error("should be error")
	}
	wg.Wait()
}

func TestTransportOnceDo(t *testing.T) {
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
	{
		trans := &Transport{
			Options: opts,
		}
		err = trans.Ping(addr)
		if err != nil {
			t.Error(err)
		}
		trans.Close()
	}
	{
		trans := &Transport{
			MaxConnsPerHost:     1,
			MaxIdleConnsPerHost: 2,
			Options:             opts,
		}
		err = trans.Ping(addr)
		if err != nil {
			t.Error(err)
		}
		trans.Close()
	}
	server.Close()
	wg.Wait()
}

func TestNewPersistConn(t *testing.T) {
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
	{
		trans := &Transport{
			MaxConnsPerHost:     1,
			MaxIdleConnsPerHost: 1,
			KeepAlive:           time.Millisecond * 100,
			IdleConnTimeout:     time.Millisecond * 100,
			Network:             network,
			Codec:               codec,
		}
		err = trans.Ping(addr)
		if err != nil {
			t.Error(err)
		}
		trans.Close()
	}
	server.Close()
	wg.Wait()
}

func TestGetConn(t *testing.T) {
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
	trans := &Transport{
		MaxConnsPerHost:     4,
		MaxIdleConnsPerHost: 2,
		KeepAlive:           time.Millisecond * 100,
		IdleConnTimeout:     time.Second * 60,
		Network:             network,
		Codec:               codec,
		ticker:              time.Millisecond * 10,
	}
	err = trans.Ping(addr)
	if err != nil {
		t.Error(err)
	}
	cwg := sync.WaitGroup{}
	for i := 0; i < 64; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			trans.Ping(addr)
		}()
	}
	time.Sleep(time.Millisecond * 300)
	for i := 0; i < 64; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			trans.Ping(addr)
		}()
	}
	trans.CloseIdleConnections()
	time.Sleep(time.Millisecond * 300)
	trans.CloseIdleConnections()
	for i := 0; i < 64; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			trans.Ping(addr)
		}()
	}
	cwg.Wait()
	server.Close()
	wg.Wait()
}

func TestTransportCloseIdleConnections(t *testing.T) {
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
	trans := &Transport{
		MaxConnsPerHost:     4,
		MaxIdleConnsPerHost: 2,
		KeepAlive:           time.Millisecond * 100,
		IdleConnTimeout:     time.Second,
		Network:             network,
		Codec:               codec,
		ticker:              time.Millisecond * 10,
	}
	err = trans.Ping(addr)
	if err != nil {
		t.Error(err)
	}
	cwg := sync.WaitGroup{}
	for i := 0; i < 64; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			trans.Ping(addr)
		}()
	}
	cwg.Wait()
	trans.CloseIdleConnections()
	trans.Close()
	server.Close()
	wg.Wait()
}

func TestTransportClose(t *testing.T) {
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
	trans := &Transport{
		MaxConnsPerHost:     4,
		MaxIdleConnsPerHost: 2,
		KeepAlive:           time.Millisecond * 100,
		IdleConnTimeout:     time.Second,
		Network:             network,
		Codec:               codec,
		ticker:              time.Millisecond * 10,
	}
	err = trans.Ping(addr)
	if err != nil {
		t.Error(err)
	}
	cwg := sync.WaitGroup{}
	for i := 0; i < 64; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			trans.Ping(addr)
		}()
	}
	cwg.Wait()
	time.Sleep(time.Millisecond * 300)
	trans.Close()
	trans.Close()
	server.Close()
	wg.Wait()
}

func TestTransportRun(t *testing.T) {
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
	trans := &Transport{
		MaxConnsPerHost:     4,
		MaxIdleConnsPerHost: 2,
		KeepAlive:           time.Millisecond * 100,
		IdleConnTimeout:     time.Millisecond * 200,
		Network:             network,
		Codec:               codec,
		ticker:              time.Millisecond * 10,
	}
	err = trans.Ping(addr)
	if err != nil {
		t.Error(err)
	}
	cwg := sync.WaitGroup{}
	for i := 0; i < 64; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			trans.Ping(addr)
		}()
	}
	time.Sleep(time.Second)
	trans.Close()
	cwg.Wait()
	server.Close()
	cwg.Wait()
	wg.Wait()
}

func TestConnQueue(t *testing.T) {
	capacity := 1
	cq := newConnQueue(capacity, ":9999")
	if cq.Capacity() != capacity {
		t.Error(cq.Capacity())
	}
	if cq.Front() != nil {
		t.Error(cq.Front())
	}
	if cq.Rear() != nil {
		t.Error(cq.Rear())
	}
	if !cq.Enqueue(&persistConn{}) {
		t.Error("should be ok")
	}
	if cq.Enqueue(&persistConn{}) {
		t.Error("should not be ok")
	}
	if cq.Length() != capacity {
		t.Error(cq.Length())
	}
	if cq.Front() == nil {
		t.Error("should not be nil")
	}
	if cq.Rear() == nil {
		t.Error("should not be nil")
	}
	p := cq.Dequeue()
	if p == nil {
		t.Error("should not be nil")
	}
	p = cq.Dequeue()
	if p != nil {
		t.Error("should be nil")
	}
}
