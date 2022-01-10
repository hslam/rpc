// Copyright (c) 2022 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"crypto/tls"
	"github.com/hslam/socket"
)

// LoadTLSConfig returns a TLS config by loading the certificate file and the key file.
func LoadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	return socket.LoadTLSConfig(certFile, keyFile)
}

// TLSConfig returns a TLS config by the certificate data and the key data.
func TLSConfig(certPEM []byte, keyPEM []byte) *tls.Config {
	return socket.TLSConfig(certPEM, keyPEM)
}

// DefalutTLSConfig returns a default TLS config.
func DefalutTLSConfig() *tls.Config {
	return socket.DefalutTLSConfig()
}

// SkipVerifyTLSConfig returns a insecure skip verify TLS config.
func SkipVerifyTLSConfig() *tls.Config {
	return socket.SkipVerifyTLSConfig()
}

// DefaultKeyPEM represents the default private key data.
var DefaultKeyPEM = socket.DefaultKeyPEM

// DefaultCertPEM represents the default certificate data.
var DefaultCertPEM = socket.DefaultCertPEM
