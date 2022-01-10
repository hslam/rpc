// Copyright (c) 2022 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"crypto/tls"
	"github.com/hslam/socket"
)

// LoadServerTLSConfig returns a server TLS config by loading the certificate file and the key file.
func LoadServerTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	return socket.LoadServerTLSConfig(certFile, keyFile)
}

// LoadClientTLSConfig returns a client TLS config by loading the certificate file.
func LoadClientTLSConfig(certFile, serverName string) (*tls.Config, error) {
	return socket.LoadClientTLSConfig(certFile, serverName)
}

// ServerTLSConfig returns a server TLS config by the certificate data and the key data.
func ServerTLSConfig(certPEM []byte, keyPEM []byte) *tls.Config {
	return socket.ServerTLSConfig(certPEM, keyPEM)
}

// ClientTLSConfig returns a client TLS config by the certificate data.
func ClientTLSConfig(certPEM []byte, serverName string) *tls.Config {
	return socket.ClientTLSConfig(certPEM, serverName)
}

// DefalutServerTLSConfig returns a default server TLS config.
func DefalutServerTLSConfig() *tls.Config {
	return socket.DefalutServerTLSConfig()
}

// DefalutClientTLSConfig returns a default client TLS config.
func DefalutClientTLSConfig() *tls.Config {
	return socket.DefalutClientTLSConfig()
}

// SkipVerifyTLSConfig returns a client TLS config which skips security verification.
func SkipVerifyTLSConfig() *tls.Config {
	return socket.SkipVerifyTLSConfig()
}

// DefalutServerName returns a server name with subdomain
func DefalutServerName(sub string) string {
	return socket.DefalutServerName(sub)
}

// DefaultKeyPEM represents the default private key data.
var DefaultKeyPEM = socket.DefaultKeyPEM

// DefaultCertPEM represents the default certificate data.
var DefaultCertPEM = socket.DefaultCertPEM
