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
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCVHlKfRTju4SuA
R4sbjcy9hNmOOTmvJPBa4gwXkmwpbbrHVUmo5aMSGYiuRfkRQMFGCiJlSG0/TcWe
i4Jbh2BK2jp+dtPaeBoIVgTsAtfN4KD6sSvyvCSH2dcTb+fKpBiqHS3IhqBACkch
qogy/MDAyPioA++MW9rSOthnquDVjMU3rkCf+NgyoeHSrf4MllLG+RRzGX3jWler
0ef8uw+MEK/Y/imQFy6C5ICipF6yQyhmeU62K9Pcm6MShYL40c1xPodAFhT7SHlR
tvldw1ziV7TXzPD5k4ZhV/Vyy5oVEt4Hvp6CjPjt1YJC8XQn+s/wvCeqGD2nPNQX
wtgTOVM9AgMBAAECggEAey+9mY2Z5t1lDmgL5wtRZA7nmrJzkNi3Jp0u2BpB+EeJ
0ToHy9tIx58IZs+vXi1cfPvKRll6xpz88GjXm71OMwfs4qRPh19IQjKthjsjBBTZ
Z8ANSk1a8E9pecksdx7wsTfBprJwl/blpE44jcZ3hcuAf2wg7JbFQn8SXzGu4zpQ
q26JiT0Ra/CXzbP7IHPPv8qarBnjHKD1lmGyLIGlf6Rh6OKmgyeo1/m3ra2m2kAS
OELSuHaVrBrQsYvMU7ETSgvHVkLLN9y9bvy/X1t59YL1oIrha/8UyukFKkjiInvt
heWWhPxcDRLktfWK5zIbyj3u9vggDA/Kf07IbGxVAQKBgQDFv0W8zo4dflns2ZuY
dZIcFJNdv/Yo3/fP1sRXF9EGOdzRRx+Iu5N+Sdk89bh/LIy/tjQHmsslxkMgeywL
4VGBTE+BFO08FM+r2YuKfupM3BXNQvrldadlXd31KogeQHALAWKJ2cMMZxAmzouH
j/80KyhYrE35a/x9VZrZm+6LbQKBgQDBC87pC+LvA+/gQFcu0CkTcXRl7KvpzBSJ
ugwbY/W8Eg5898tNEhtXGiU8HN3a+8JqxHwTg1w5o/O2yrx9Ws8nh/PXMhdE7I/3
Tu1PAAJrzx52x/FDON6aDaTIt2k74reXfnx6cSFghM8SEIp2OhIGund9G9Roty0X
YdJQsEa1EQKBgFxSSXe1o6nnZIpsqfUK5vUPMiHxzjYVIng5V58lsmPKvepC31kR
4fFy/uYz/jf5j5itsyrdvPxczNgsSUsendPU0cV9BKkpOi+MOFannDHYCqGzJLne
LRHpOggNHFGrWeP5eIzNSv/OWj8T7RaURtyPTZ3gi+Ln5JCLV+lCoKMdAoGAJwL3
4Wihh6PICg12kONIKcG3wBE//JNdYyfR4ock1cjgXKjG0OBj3gpOlANRYjuWYnUq
jdbyAEP9sGbwCHUdf+Odh1N8GFWmElhE5L4fvyGwClkFjIwlkARJ1LYb8hoy9857
4VKTaCnunrvw/0tk8S8ljobdOfwqhJskIWI+J8ECgYA5P+t9qh1Ux4OY11lHlbM0
T/orcVhaXUfmf5bXKyriZLT/2LncAUAOkQVUn8nxcxm9gnPqfmmHlLnj7S5JEqSQ
bbVsfZwy1ZeIQS55NraQ6kbZcLoKYYj90Lw9uJeg9ifc9Xgjl08cB7vtOdYlq/3l
6hEsOcxSDQe0TFZ++cEojA==
-----END PRIVATE KEY-----
`)
	certPEM := []byte(`-----BEGIN CERTIFICATE-----
MIIDYjCCAkoCCQCsvKhoqAn3AjANBgkqhkiG9w0BAQsFADByMQswCQYDVQQGEwJD
TjELMAkGA1UECAwCQkoxCzAJBgNVBAcMAkJKMQ0wCwYDVQQKDARURVNUMQ0wCwYD
VQQLDARURVNUMQ0wCwYDVQQDDARURVNUMRwwGgYJKoZIhvcNAQkBFg10ZXN0QHRl
c3QuY29tMCAXDTE5MDgyMjIwMDU0MFoYDzIwNjkwODA5MjAwNTQwWjByMQswCQYD
VQQGEwJDTjELMAkGA1UECAwCQkoxCzAJBgNVBAcMAkJKMQ0wCwYDVQQKDARURVNU
MQ0wCwYDVQQLDARURVNUMQ0wCwYDVQQDDARURVNUMRwwGgYJKoZIhvcNAQkBFg10
ZXN0QHRlc3QuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlR5S
n0U47uErgEeLG43MvYTZjjk5ryTwWuIMF5JsKW26x1VJqOWjEhmIrkX5EUDBRgoi
ZUhtP03FnouCW4dgSto6fnbT2ngaCFYE7ALXzeCg+rEr8rwkh9nXE2/nyqQYqh0t
yIagQApHIaqIMvzAwMj4qAPvjFva0jrYZ6rg1YzFN65An/jYMqHh0q3+DJZSxvkU
cxl941pXq9Hn/LsPjBCv2P4pkBcuguSAoqReskMoZnlOtivT3JujEoWC+NHNcT6H
QBYU+0h5Ubb5XcNc4le018zw+ZOGYVf1csuaFRLeB76egoz47dWCQvF0J/rP8Lwn
qhg9pzzUF8LYEzlTPQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBZCb7jJVgbzqk3
GlcAXC82s6aylrdRnfsyMm5NVyFdlTYWUWV9UySMX1vKOKUs+zzgNQd4RVQ1zAJE
DeUymsGJubpIIj5qToS6GmD+CLBLr2YtDnydjjsONa/0Yl43LDsylFPlPvArvScT
GfN2WLx6yLCLki1oX/1ZDKs3thvsaBH3gWky/tBd3WjobGtygEy9u5bgL+LEvkp1
ZpEA0aRYG9eIpyADngfPEtcU/XmmJdqphPJGyNlfPf2aP26s3LdhCh4aTGAwAQZb
GxEoHEOaqbTS/2QjHgVEr+2nfIVkWzyHi992nxwTESFR7cmmyzJTI+uL+73xZQyd
NTEWX480
-----END CERTIFICATE-----
`)
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
