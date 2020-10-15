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
		MaxConnsPerHost:     1,
		MaxIdleConnsPerHost: 1,
		KeepAlive:           time.Millisecond * 100,
		IdleConnTimeout:     time.Millisecond * 100,
		Options:             opts,
	}

	err = trans.Ping(addr)
	if err != nil {
		t.Error(err)
	}
	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	cwg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
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
	time.Sleep(time.Second * 3)
	trans.Close()
	server.Close()
	res := service.ArithResponse{}
	call := trans.Go(addr, "Arith.Multiply", req, &res, make(chan *Call, 1))
	<-call.Done
	if call.Error == nil {
		t.Error("should be error")
	}
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
