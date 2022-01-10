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
	_, err = LoadTLSConfig("", "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	certFile, _ := os.Create(certFileName)
	certFile.Write(DefaultCertPEM)
	certFile.Close()
	defer os.Remove(certFileName)
	_, err = LoadTLSConfig(certFileName, "")
	if err == nil {
		t.Error("should be no such file or directory")
	}
	keyFile, _ := os.Create(keyFileName)
	keyFile.Write(DefaultKeyPEM)
	keyFile.Close()
	defer os.Remove(keyFileName)
	_, err = LoadTLSConfig(certFileName, keyFileName)
	if err != nil {
		t.Error(err)
	}
}

func TestTLSConfig(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("should panic")
		}
	}()
	TLSConfig(DefaultCertPEM, []byte{})
}
