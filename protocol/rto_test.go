package protocol

import (
	"testing"
)
//go test -v -run="none" -bench=. -benchtime=3s  -benchmem
func BenchmarkRTO(t *testing.B) {
	var rto = &RTO{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		rtt:=int64(i%100)
		rto.updateRTT(rtt)
	}
}
