package rpc

import (
	"crypto/tls"
	"crypto/rsa"
	"math/big"
	"encoding/pem"
	"crypto/rand"
	"crypto/x509"
)

func uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}
func bytesToUint16(array []byte) uint16 {
	var data uint16 =0
	for i:=0;i< len(array);i++  {
		data = data+uint16(uint(array[i])<<uint(8*i))
	}
	return data
}
func uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}
func bytesToUint32(array []byte) uint32 {
	var data uint32 =0
	for i:=0;i< len(array);i++  {
		data = data+uint32(uint(array[i])<<uint(8*i))
	}
	return data
}

func checkSum(b []byte) uint16 {
	sum := 0
	for n := 1; n < len(b)-1; n += 2 {
		sum += int(b[n])<<8 + int(b[n+1])
	}
	lowbit:=sum & 0xffff
	highbit:=sum >> 16
	checksum := lowbit +highbit
	checksum += (checksum >> 16)
	var ans = uint16(^checksum)
	return ans
}
// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
