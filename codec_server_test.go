package rpc

import (
	"github.com/hslam/socket"
	"sync"
	"testing"
)

func TestNewServerCodec(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	wg := sync.WaitGroup{}
	wg.Add(1)
	sock, _ := socket.NewSocket(network, nil)
	lis, _ := sock.Listen(addr)
	go func() {
		defer wg.Done()
		conn, _ := lis.Accept()
		message := conn.Messages()
		headerEncoder := NewHeaderEncoder("json")
		if NewServerCodec(nil, nil, nil) != nil {
			t.Error("should be nil")
		}
		if NewServerCodec(nil, nil, message) != nil {
			t.Error("should be nil")
		}
		codec := NewServerCodec(nil, headerEncoder(), message)
		if codec == nil {
			t.Error("should not be nil")
		}
		if codec.ReadRequestBody(nil) == nil {
			t.Error("The err should not be nil")
		}
		message.Close()
		lis.Close()
	}()
	sock.Dial(addr)
	wg.Wait()
}
