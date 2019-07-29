package protocol

import (
	"testing"
)
//go test -v -run="none" -bench=. -benchtime=3s  -benchmem
func BenchmarkUint16ToBytes(t *testing.B) {
	var n uint16 = 1 << 15
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		uint16ToBytes(n)
	}
}
func BenchmarkBytesToUint16(t *testing.B) {
	var n uint16 = 1 << 15
	var b=uint16ToBytes(n)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		bytesToUint16(b)
	}
}
func BenchmarkUint32ToBytes(t *testing.B) {
	var n uint32 = 1 << 31
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		uint32ToBytes(n)
	}
}
func BenchmarkBytesToUint32(t *testing.B) {
	var n uint32 = 1 << 31
	var b=uint32ToBytes(n)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		bytesToUint32(b)
	}
}
func BenchmarkPacket64(t *testing.B) {
	var msg=make([]byte,64)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		PacketStream(msg)
	}
}
func BenchmarkUnpack64(t *testing.B) {
	var msg=make([]byte,64)
	var pmsg=PacketStream(msg)
	var readerChannel=make(chan []byte,1)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		UnpackStream(pmsg,readerChannel)
		<-readerChannel
	}
}

func BenchmarkPacket512(t *testing.B) {
	var msg=make([]byte,512)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		PacketStream(msg)
	}
}
func BenchmarkUnpack512(t *testing.B) {
	var msg=make([]byte,512)
	var pmsg=PacketStream(msg)
	var readerChannel=make(chan []byte,1)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		UnpackStream(pmsg,readerChannel)
		<-readerChannel
	}
}

func BenchmarkPacket1500(t *testing.B) {
	var msg=make([]byte,1500)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		PacketStream(msg)
	}
}
func BenchmarkUnpack1500(t *testing.B) {
	var msg=make([]byte,1500)
	var pmsg=PacketStream(msg)
	var readerChannel=make(chan []byte,1)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		UnpackStream(pmsg,readerChannel)
		<-readerChannel
	}
}

func BenchmarkPacket65536(t *testing.B) {
	var msg=make([]byte,65536)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		PacketStream(msg)
	}
}
func BenchmarkUnpack65536(t *testing.B) {
	var msg=make([]byte,65536)
	var pmsg=PacketStream(msg)
	var readerChannel=make(chan []byte,1)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		UnpackStream(pmsg,readerChannel)
		<-readerChannel
	}
}
