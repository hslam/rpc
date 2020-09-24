package rpc

import (
	"github.com/hslam/rpc/examples/codec/json/service"
	"github.com/hslam/socket"
	"sync"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
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
	network = ""
	if _, err := Dial(network, addr, codec); err == nil {
		t.Error(err)
	}
	codec = ""
	if _, err = Dial(network, addr, codec); err == nil {
		t.Error(err)
	}
	DefaultServer.Close()
	wg.Wait()
}

func TestDialTLS(t *testing.T) {
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
	network = ""
	if _, err := DialTLS(network, addr, codec, socket.SkipVerifyTLSConfig()); err == nil {
		t.Error(err)
	}
	codec = ""
	if _, err = DialTLS(network, addr, codec, socket.SkipVerifyTLSConfig()); err == nil {
		t.Error(err)
	}
	DefaultServer.Close()
	wg.Wait()
}

func TestDialWithOptions(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	opts := &Options{}
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
	opts.Network = ""
	if _, err := DialWithOptions(addr, opts); err == nil {
		t.Error(err)
	}
	opts.Codec = ""
	if _, err = DialWithOptions(addr, opts); err == nil {
		t.Error(err)
	}
	DefaultServer.Close()
	wg.Wait()
}
