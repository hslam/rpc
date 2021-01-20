// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"strings"
	"testing"
)

func TestAssignPool(t *testing.T) {
	if len(GetBuffer(0)) > 0 {
		t.Error()
	}
	b := GetBuffer(1024)
	if len(b) < 1024 {
		t.Error(len(b))
	}
	PutBuffer(b)
	assignPool(1024)

}

func TestJSONCodec(t *testing.T) {
	type Object struct {
		A bool `json:"A" xml:"A"`
	}
	var obj = Object{A: true}
	var c = JSONCodec{}
	var buf = make([]byte, 512)
	var objCopy *Object
	data, _ := c.Marshal(buf, &obj)
	c.Unmarshal(data, &objCopy)
}

func TestXMLCodec(t *testing.T) {
	type Object struct {
		A bool `json:"A" xml:"A"`
	}
	var obj = Object{A: true}
	var c = XMLCodec{}
	var buf = make([]byte, 512)
	var objCopy *Object
	data, _ := c.Marshal(buf, &obj)
	c.Unmarshal(data, &objCopy)
}

func TestBYTESCodec(t *testing.T) {
	var obj = []byte{128, 8, 128, 8, 195, 245, 72, 64, 74, 216, 18, 77, 251, 33, 9, 64, 10, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100, 1, 1, 255, 2, 1, 128, 1, 255}
	var c = BYTESCodec{}
	var buf = make([]byte, 512)
	var objCopy []byte
	data, _ := c.Marshal(buf, &obj)
	c.Unmarshal(data, &objCopy)
}

//codeObject is a test struct
type codeObject struct {
}

//Marshal marshals the Object into buf and returns the bytes.
func (o *codeObject) Marshal(buf []byte) ([]byte, error) {
	return nil, nil
}

//Unmarshal unmarshals the Object from buf and returns the number of bytes read (> 0).
func (o *codeObject) Unmarshal(data []byte) (uint64, error) {
	return 0, nil
}

func TestCODECodec(t *testing.T) {
	{
		var obj = codeObject{}
		var c = CODECodec{}
		var buf = make([]byte, 512)
		var objCopy codeObject
		data, _ := c.Marshal(buf, &obj)
		c.Unmarshal(data, &objCopy)
	}
	{
		type Object struct {
		}
		var obj = Object{}
		var c = CODECodec{}
		var buf = make([]byte, 512)
		var objCopy Object
		if _, err := c.Marshal(buf, &obj); err != ErrorCODE {
			t.Error(ErrorCODE)
		}
		if err := c.Unmarshal(nil, &objCopy); err != ErrorCODE {
			t.Error(ErrorCODE)
		}
	}
}

type gogopbObject struct {
	size int
}

func (m *gogopbObject) Size() (n int) {
	return m.size
}

func (m *gogopbObject) Marshal() (dAtA []byte, err error) {
	return
}

func (m *gogopbObject) MarshalTo(dAtA []byte) (int, error) {
	return 0, nil
}

func (m *gogopbObject) Unmarshal(dAtA []byte) error {
	return nil
}

func TestGOGOPBCodec(t *testing.T) {
	{
		var obj = gogopbObject{}
		var c = GOGOPBCodec{}
		var buf = make([]byte, 512)
		var objCopy gogopbObject
		data, _ := c.Marshal(buf, &obj)
		c.Unmarshal(data, &objCopy)
	}
	{
		var obj = gogopbObject{size: 128}
		var c = GOGOPBCodec{}
		var buf = make([]byte, 0)
		var objCopy gogopbObject
		data, _ := c.Marshal(buf, &obj)
		c.Unmarshal(data, &objCopy)
	}
	{
		type Object struct {
		}
		var obj = Object{}
		var c = GOGOPBCodec{}
		var buf = make([]byte, 512)
		var objCopy Object
		if _, err := c.Marshal(buf, &obj); err != ErrorGOGOPB {
			t.Error(ErrorGOGOPB)
		}
		if err := c.Unmarshal(nil, &objCopy); err != ErrorGOGOPB {
			t.Error(ErrorGOGOPB)
		}
	}
}

type msgpObject struct {
}

// MarshalMsg implements msgp.Marshaler
func (z *msgpObject) MarshalMsg(b []byte) (o []byte, err error) {
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *msgpObject) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return
}

func TestMSGPCodec(t *testing.T) {
	{
		var obj = msgpObject{}
		var c = MSGPCodec{}
		var buf = make([]byte, 512)
		var objCopy msgpObject
		data, _ := c.Marshal(buf, &obj)
		c.Unmarshal(data, &objCopy)
	}
	{
		type Object struct {
		}
		var obj = Object{}
		var c = MSGPCodec{}
		var buf = make([]byte, 512)
		var objCopy Object
		if _, err := c.Marshal(buf, &obj); err != ErrorMSGP {
			t.Error(ErrorMSGP)
		}
		if err := c.Unmarshal(nil, &objCopy); err != ErrorMSGP {
			t.Error(ErrorMSGP)
		}
	}
}

func TestCodecCODE(t *testing.T) {
	var req = request{Seq: 1024, Upgrade: make([]byte, 512), ServiceMethod: strings.Repeat("a", 512), Args: make([]byte, 512)}
	var res = response{Seq: 1024, Error: strings.Repeat("e", 512), Reply: make([]byte, 512)}
	var c = CODECodec{}
	var data []byte
	{
		var bufq = make([]byte, 10240)
		var bufs = make([]byte, 10240)
		var reqCopy request
		var resCopy response
		data, _ = c.Marshal(bufq, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(bufs, &res)
		c.Unmarshal(data, &resCopy)
	}
	{
		var reqCopy request
		var resCopy response
		data, _ = c.Marshal(nil, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(nil, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkRoundtripCODE(t *testing.B) {
	var req = request{Seq: 1024, Upgrade: make([]byte, 1), ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = CODECodec{}
	var bufq = make([]byte, 1024)
	var bufs = make([]byte, 1024)
	var reqCopy request
	var resCopy response
	qd, _ := c.Marshal(bufq, &req)
	sd, _ := c.Marshal(bufs, &res)
	t.SetBytes(int64(len(qd) + len(sd)))
	var data []byte
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ = c.Marshal(bufq, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(bufs, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func TestCodecPB(t *testing.T) {
	var req = pbRequest{Seq: 1024, Upgrade: make([]byte, 512), ServiceMethod: strings.Repeat("a", 512), Args: make([]byte, 512)}
	var res = pbResponse{Seq: 1024, Error: strings.Repeat("e", 512), Reply: make([]byte, 512)}
	var c = GOGOPBCodec{}
	var data []byte
	{
		var bufq = make([]byte, 10240)
		var bufs = make([]byte, 10240)
		var reqCopy pbRequest
		var resCopy pbResponse
		data, _ = c.Marshal(bufq, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(bufs, &res)
		c.Unmarshal(data, &resCopy)
	}
	{
		var reqCopy pbRequest
		var resCopy pbResponse
		data, _ = c.Marshal(nil, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(nil, &res)
		c.Unmarshal(data, &resCopy)
	}
	{
		if n, err := req.MarshalTo(nil); n > 0 || err == nil {
			t.Error("The err should not be nil")
		}
		if n, err := res.MarshalTo(nil); n > 0 || err == nil {
			t.Error("The err should not be nil")
		}
	}
}

func BenchmarkRoundtripPB(t *testing.B) {
	var req = pbRequest{Seq: 1024, Upgrade: make([]byte, 1), ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = pbResponse{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = GOGOPBCodec{}
	var bufq = make([]byte, 1024)
	var bufs = make([]byte, 1024)
	var reqCopy pbRequest
	var resCopy pbResponse
	qd, _ := c.Marshal(bufq, &req)
	sd, _ := c.Marshal(bufs, &res)
	t.SetBytes(int64(len(qd) + len(sd)))
	var data []byte
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ = c.Marshal(bufq, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(bufs, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func TestCodecJSON(t *testing.T) {
	var req = jsonRequest{Seq: 1024, Upgrade: make([]byte, 512), ServiceMethod: strings.Repeat("a", 512), Args: make([]byte, 512)}
	var res = jsonResponse{Seq: 1024, Error: strings.Repeat("e", 512), Reply: make([]byte, 512)}
	var c = JSONCodec{}
	var bufq = make([]byte, 10240)
	var bufs = make([]byte, 10240)
	var reqCopy jsonRequest
	var resCopy jsonResponse
	var data []byte
	data, _ = c.Marshal(bufq, &req)
	c.Unmarshal(data, &reqCopy)
	data, _ = c.Marshal(bufs, &res)
	c.Unmarshal(data, &resCopy)
}

func BenchmarkRoundtripJSON(t *testing.B) {
	var req = jsonRequest{Seq: 1024, Upgrade: make([]byte, 1), ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = jsonResponse{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = JSONCodec{}
	var bufq = make([]byte, 1024)
	var bufs = make([]byte, 1024)
	var reqCopy jsonRequest
	var resCopy jsonResponse
	qd, _ := c.Marshal(bufq, &req)
	sd, _ := c.Marshal(bufs, &res)
	t.SetBytes(int64(len(qd) + len(sd)))
	var data []byte
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ = c.Marshal(bufq, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(bufs, &res)
		c.Unmarshal(data, &resCopy)
	}
}

//goos: darwin
//goarch: amd64
//pkg: github.com/hslam/rpc
//BenchmarkRequestMarshalCODE-4        	33482854	        31.1 ns/op	17055.25 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRequestMarshalCODEPB-4      	34093851	        34.6 ns/op	15421.15 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRequestMarshalGOGOPB-4      	17655142	        66.6 ns/op	8022.43 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRequestMarshalJSON-4        	  894202	      1247 ns/op	 579.00 MB/s	    1472 B/op	       2 allocs/op
//BenchmarkRequestUnmarshalCODE-4      	46397149	        25.0 ns/op	21232.64 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRequestUnmarshalCODEPB-4    	38636697	        30.6 ns/op	17476.32 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRequestUnmarshalGOGOPB-4    	16259637	        71.6 ns/op	7453.41 MB/s	      16 B/op	       1 allocs/op
//BenchmarkRequestUnmarshalJSON-4      	  183936	      6344 ns/op	 113.81 MB/s	     800 B/op	       6 allocs/op
//BenchmarkRequestRoundtripCODE-4      	19764566	        58.9 ns/op	9017.95 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRequestRoundtripCODEPB-4    	17977524	        64.9 ns/op	8226.33 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRequestRoundtripGOGOPB-4    	 8483186	       140 ns/op	3807.32 MB/s	      16 B/op	       1 allocs/op
//BenchmarkRequestRoundtripJSON-4      	  149979	      7850 ns/op	  91.97 MB/s	    2272 B/op	       8 allocs/op
//BenchmarkResponseMarshalCODE-4       	43076028	        27.8 ns/op	18626.74 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseMarshalCODEPB-4     	41994508	        27.7 ns/op	18668.33 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseMarshalGOGOPB-4     	20898134	        56.2 ns/op	9221.58 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseMarshalJSON-4       	  972465	      1225 ns/op	 577.77 MB/s	    1472 B/op	       2 allocs/op
//BenchmarkResponseUnmarshalCODE-4     	57120253	        20.3 ns/op	25487.98 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseUnmarshalCODEPB-4   	52179763	        22.8 ns/op	22700.75 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseUnmarshalGOGOPB-4   	31025642	        39.5 ns/op	13116.48 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseUnmarshalJSON-4     	  166291	      6240 ns/op	 113.46 MB/s	     784 B/op	       5 allocs/op
//BenchmarkResponseRoundtripCODE-4     	23979372	        49.2 ns/op	10497.67 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseRoundtripCODEPB-4   	22384860	        52.3 ns/op	9901.09 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseRoundtripGOGOPB-4   	12373190	        95.2 ns/op	5439.52 MB/s	       0 B/op	       0 allocs/op
//BenchmarkResponseRoundtripJSON-4     	  149359	      7613 ns/op	  93.00 MB/s	    2256 B/op	       7 allocs/op
//BenchmarkRoundtripCODE-4             	11107009	       105 ns/op	9985.33 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRoundtripCODEPB-4           	10183730	       116 ns/op	9085.41 MB/s	       0 B/op	       0 allocs/op
//BenchmarkRoundtripGOGOPB-4           	 4641254	       259 ns/op	4063.59 MB/s	      16 B/op	       1 allocs/op
//BenchmarkRoundtripJSON-4             	   76912	     15405 ns/op	  92.83 MB/s	    4529 B/op	      15 allocs/op
//PASS
//ok  	github.com/hslam/rpc	34.985s
