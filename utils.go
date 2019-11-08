package rpc

import (
	"crypto/tls"
	"crypto/rsa"
	"math/big"
	"encoding/pem"
	"crypto/rand"
	"crypto/x509"
	"encoding/binary"
	"math"
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
func float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func byteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func byteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
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

func DefalutTLSConfig() *tls.Config {
	keyPEM := []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQx2lxR73msUrD
EthvjZCG8nmxY/m6ZXBRdaygc++Ie4baSHgS5DsLUcMEPhUOzzkUjjvc9dFy/ROD
yw/hL1jmgC2a8bA+czWlK1V6tP/bkEzy9MZpUIM7aY48W0/t4O9ZaRML8zavfozo
Vq5UXoS6meIzS5afg6vVNFB9bLV/d/AmT86BdZ0PYUzYg4KJcGT2nFJnEecjBFjj
khIzhQs02IxaHH7xiddcIjZpzXEcaGFxJGecDFe+LFHIN9pgCG17kppK2Ebt93ck
EpIM1LuiemX2CBckBlXCACrIVkeQ6pJTKcEMFbEcU9zYz8G9Apu3CgZApQUxuSlc
VbdStupHAgMBAAECggEBAMfqm0I441Pa4u8Gxa+UDBfcacD+LaxZ5AJsvt9qBK1Z
M5kjma7lUtCU+zu9wuZKcasIQ7RiwqvsQFqMAGmtr+AQTxs2YYB7S5wccZ6tYO67
L6PZ1YAU84TZn7SV72SmZirknbNssim78WutTQNG+qvAHMTnlZSrPchLbuObS/SA
hufs9fqxzBUeVXVE18vyhBq4i2SRyKcBxKEPRdzhbYem8SPi7hr/NnPomNUOr8Cw
JGY0XsD6MVs4myzdpPRTNL8farzvpmLdJ9C1Nc8DKIwZZsoGkoMXXxtqJzYUYp7C
deqUCMJY381R3mJUpUd6RU2ocL521CyA2YinGKljspkCgYEA60CC0B0I2cQ+W1Bf
i/E5rItT91cseCJTnpH9m2ok3Pe4CiO/BLK9Ew9rjlHb7wNofch/EwNb4veIxRf8
2MMjBuecTvqn/akirtPtjDPuQoKW7yBEKjqW8+yt97INHhMVkhMkoZ4JTmq2gWJg
oT02Eq9YN52dpJjqYM7OvSywHPUCgYEA4zEyKo1/KWNNCCyBMlmPzRNyhxLk0h+t
84UlwTKRogsQLgkZrseRMKakiHjF4Ez6qczC5MrHqH4I3e5RrAQP9uNQJIA0g+Ma
mhpI5xq1ZnkCq/Z69K808trzFS4rAS7NZfwn8mV0ZXFHzmTDH1tDEvg+NXUXykFE
zI15yYyXpMsCgYEA1QSaLvZLgFyplifGDMLGVY3n3yzJcJKsowZQ3PyVGp0Ywd2y
Zv+uI2cwHjPTca7lXBhDsKS2/GLmLonVAzZXLjZlHELuAMu5QxNVo0GWuhTjtO3D
q3VYINGsiYBpTlU7kATTg6DFjoMkdS3uj7IMl4i82cdX6qYofLZnD3c6lU0CgYBZ
t1uwIiBNH8GTsL90Opnmyf84B+YEdC4lNDcsi+Omsee5xi42LujO5X+jxM2fPcbe
ttVftBQUHXEy8qGd5BzJygoj39zdGBmxMSAI4ysvRCoh7juv1GB8ZqoHeyvQU8MY
uvKrbhUA2jMY9gF3qHpcS1uFkK/MVunsPRIS3Uok8QKBgGq5Auz3ByjA9gKPk8qO
p8S35+ldnAnIsMeojW+rXmUBrHaSp9ea9+obFHp1JKjI4dMoT92EKRAU5nnM3iqR
q2/umYjD1XcbHSRdxz8H7hkZfQ4GSAwG/KNeGWL5FRSq3p6WeobLL4WdTcugS5cz
cXoVckT5dS+Mj/SxSJ/oyvyu
-----END PRIVATE KEY-----
`)
	certPEM := []byte(`-----BEGIN CERTIFICATE-----
MIIDbjCCAlYCCQCZymVmCboLuDANBgkqhkiG9w0BAQsFADB4MQswCQYDVQQGEwJD
TjELMAkGA1UECAwCQkoxCzAJBgNVBAcMAkJKMQ4wDAYDVQQKDAVobHNhbTEOMAwG
A1UECwwFaHNsYW0xDjAMBgNVBAMMBWhzbGFtMR8wHQYJKoZIhvcNAQkBFhA3OTE4
NzQxNThAcXEuY29tMCAXDTE5MTEwODEzNDYzNloYDzIxMTkxMDE1MTM0NjM2WjB4
MQswCQYDVQQGEwJDTjELMAkGA1UECAwCQkoxCzAJBgNVBAcMAkJKMQ4wDAYDVQQK
DAVobHNhbTEOMAwGA1UECwwFaHNsYW0xDjAMBgNVBAMMBWhzbGFtMR8wHQYJKoZI
hvcNAQkBFhA3OTE4NzQxNThAcXEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEA0MdpcUe95rFKwxLYb42QhvJ5sWP5umVwUXWsoHPviHuG2kh4EuQ7
C1HDBD4VDs85FI473PXRcv0Tg8sP4S9Y5oAtmvGwPnM1pStVerT/25BM8vTGaVCD
O2mOPFtP7eDvWWkTC/M2r36M6FauVF6EupniM0uWn4Or1TRQfWy1f3fwJk/OgXWd
D2FM2IOCiXBk9pxSZxHnIwRY45ISM4ULNNiMWhx+8YnXXCI2ac1xHGhhcSRnnAxX
vixRyDfaYAhte5KaSthG7fd3JBKSDNS7onpl9ggXJAZVwgAqyFZHkOqSUynBDBWx
HFPc2M/BvQKbtwoGQKUFMbkpXFW3UrbqRwIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
AQBMfGE+zL6XMc/CJK59rRJQdypFG2gPxzpHi4XXTUfTAAYhZZRMYfiefxII8s5V
MG+n+c1wU/nubE3xj9dgq7aIC1L3EPyVkWu/s8lPNWKMOO1FchZghBHYsImD5uM6
sD1euV5nOmPnirK/vrfBuemGLtOFEgDrCEk39bd8AWLgrdpeqfVpW4K6QkDh1V4u
Qe2ZVXa6qRwJ8dAvo79JmW8txruJ6/5s4Af7Gogr/F2BHYlbLMgdjAYTOT0X7mOs
ZtmM8OaZrcg7EDEFHsV3k56S6i4EUH2VDOCz+v2BAJmWqDwVmjI4kEKKbgjI9fwv
st022MDOYYU7dKz1lAiLJr1G
-----END CERTIFICATE-----
`)
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
