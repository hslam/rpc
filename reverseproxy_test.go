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
	proxy := NewSingleHostReverseProxy(addr)
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
