package rpc

import (
	"github.com/hslam/socket"
	"sync"
	"testing"
)

func TestNewClientCodec(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	wg := sync.WaitGroup{}
	wg.Add(1)
	sock, _ := socket.NewSocket(network, nil)
	lis, _ := sock.Listen(addr)
	go func() {
		defer wg.Done()
		for {
			_, err := lis.Accept()
			if err != nil {
				break
			}
		}
	}()
	conn, _ := sock.Dial(addr)
	message := conn.Messages()
	headerEncoder := NewHeaderEncoder("json")
	if NewClientCodec(nil, nil, nil) != nil {
		t.Error("should be nil")
	}
	if NewClientCodec(nil, nil, message) != nil {
		t.Error("should be nil")
	}
	codec := NewClientCodec(nil, headerEncoder(), message)
	if codec == nil {
		t.Error("should not be nil")
	}
	if codec.ReadResponseBody(nil) == nil {
		t.Error("The err should not be nil")
	}
	message.Close()
	lis.Close()
	wg.Wait()
}
