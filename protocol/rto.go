package protocol
const (
	ALPHA = 0.8
	BETA  = 1.6
)

const (
	LBOUND = 100000
	UBOUND = 60000000
)

type RTO struct {
	srtt int64
	rto  int64
}
func (t *RTO) updateRTT(rtt int64) int64 {
	var rto int64
	if t.srtt == 0 {
		t.srtt = rtt
		rto = int64(rtt)
	} else {
		t.srtt = int64(ALPHA*float64(t.srtt) + ((1 - ALPHA) * float64(rtt)))
		rto = min(UBOUND, max(LBOUND, int64(BETA*float32(t.srtt))))
	}
	return rto
}
func (t *RTO) getSRTT() int64 {
	return t.srtt
}
func min(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}