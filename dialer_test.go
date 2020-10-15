package rpc

import (
	"github.com/hslam/socket"
	"testing"
)

func TestDial(t *testing.T) {
	addr := ":8880"
	network := ""
	codec := ""
	if _, err := Dial(network, addr, codec); err == nil {
		t.Error("should error")
	}
	network = "tcp"
	if _, err := Dial(network, addr, codec); err == nil {
		t.Error("should error")
	}
}

func TestDialTLS(t *testing.T) {
	addr := ":8880"
	network := ""
	codec := ""
	if _, err := DialTLS(network, addr, codec, socket.SkipVerifyTLSConfig()); err == nil {
		t.Error("should error")
	}
	network = "tcp"
	if _, err := DialTLS(network, addr, codec, socket.SkipVerifyTLSConfig()); err == nil {
		t.Error("should error")
	}
}

func TestDialWithOptions(t *testing.T) {
	addr := ":9999"
	opts := &Options{}
	opts.Codec = ""
	opts.Network = ""
	if _, err := DialWithOptions(addr, opts); err == nil {
		t.Error("should error")
	}
	opts.Network = "tcp"
	if _, err := DialWithOptions(addr, opts); err == nil {
		t.Error("should error")
	}
}
