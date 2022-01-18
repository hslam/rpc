// Copyright (c) 2022 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"os"
	"testing"
)

func TestLoadTLSConfig(t *testing.T) {
	var caCertFileName = "tmpTestCaCertFile"
	var serverCertFileName = "tmpTestServerCertFile"
	var serverKeyFileName = "tmpTestServerKeyFile"

	var err error
	_, err = LoadServerTLSConfig("", "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	certFile, _ := os.Create(serverCertFileName)
	certFile.Write(DefaultServerCertPEM)
	certFile.Close()
	defer os.Remove(serverCertFileName)
	_, err = LoadServerTLSConfig(serverCertFileName, "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	keyFile, _ := os.Create(serverKeyFileName)
	keyFile.Write(DefaultServerKeyPEM)
	keyFile.Close()
	defer os.Remove(serverKeyFileName)
	_, err = LoadServerTLSConfig(serverCertFileName, serverKeyFileName)
	if err != nil {
		t.Error(err)
	}

	_, err = LoadClientTLSConfig("", "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	caCertFile, _ := os.Create(caCertFileName)
	caCertFile.Write(DefaultRootCertPEM)
	caCertFile.Close()
	defer os.Remove(caCertFileName)
	_, err = LoadClientTLSConfig(caCertFileName, "")
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
	ServerTLSConfig(DefaultServerCertPEM, []byte{})
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
