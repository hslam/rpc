// Copyright (c) 2022 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"os"
	"testing"
)

func TestLoadTLSConfig(t *testing.T) {
	var certFileName = "tmpTestCertFile"
	var keyFileName = "tmpTestKeyFile"
	var err error
	_, err = LoadServerTLSConfig("", "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	_, err = LoadClientTLSConfig("", "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	certFile, _ := os.Create(certFileName)
	certFile.Write(DefaultCertPEM)
	certFile.Close()
	defer os.Remove(certFileName)
	_, err = LoadServerTLSConfig(certFileName, "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	_, err = LoadClientTLSConfig(certFileName, "")
	if err != nil {
		t.Error(err)
	}
	keyFile, _ := os.Create(keyFileName)
	keyFile.Write(DefaultKeyPEM)
	keyFile.Close()
	defer os.Remove(keyFileName)
	_, err = LoadServerTLSConfig(certFileName, keyFileName)
	if err != nil {
		t.Error(err)
	}
}

func TestServerTLSConfig(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("should panic")
		}
	}()
	ServerTLSConfig(DefaultCertPEM, []byte{})
}

func TestClientTLSConfig(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("should panic")
		}
	}()
	ClientTLSConfig([]byte{}, "")
}

func TestDefalutServerName(t *testing.T) {
	sub := "Hello"
	if DefalutServerName("Hello")[:len(sub)] != sub {
		t.Error()
	}
}
