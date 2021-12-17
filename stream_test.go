// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"sync"
	"testing"
	"time"
)

func TestStream(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	server := NewServer()
	server.Register(new(mockStreamArith))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Listen(network, addr, codec)
	}()
	time.Sleep(time.Millisecond * 10)
	conn, err := Dial(network, addr, codec)
	if err != nil {
		panic(err)
	}
	stream, err := conn.NewStream("mockStreamArith.StreamMultiply")
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 3; i++ {
		req := &mockStreamArithRequest{A: int32(i), B: int32(i + 1)}
		var res mockStreamArithResponse
		if err := stream.WriteMessage(req); err != nil {
			t.Error(err)
		}
		if err := stream.ReadMessage(&res); err != nil {
			t.Error(err)
		}
	}
	stream.Close()
	conn.Close()
	server.Close()
	wg.Wait()
}

type mockStreamArithRequest struct {
	A int32
	B int32
}

type mockStreamArithResponse struct {
	Pro int32 //product
}

type mockStreamArith struct{}

func (a *mockStreamArith) StreamMultiply(stream *mockStream) error {
	for {
		var req = &mockStreamArithRequest{}
		if err := stream.Read(req); err != nil {
			return err
		}
		res := mockStreamArithResponse{}
		res.Pro = req.A * req.B
		if err := stream.Write(&res); err != nil {
			println(err.Error())
			return err
		}
	}
}

type mockStream struct {
	stream Stream
}

func (s *mockStream) Connect(stream Stream) error {
	s.stream = stream
	return nil
}

func (s *mockStream) Read(req *mockStreamArithRequest) error {
	return s.stream.ReadMessage(req)
}

func (s *mockStream) Write(res *mockStreamArithResponse) error {
	return s.stream.WriteMessage(res)
}
