package bytepool

import (
	"testing"
)
//go test -v -run="none" -bench=. -benchtime=3s  -benchmem
func BenchmarkBytePool(t *testing.B) {
	var bytePool =NewBytePool(1024,65536)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		buffer:=bytePool.Get()
		bytePool.Put(buffer)
	}
}

func BenchmarkByteMake(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		GetMake(65536)
	}
}
