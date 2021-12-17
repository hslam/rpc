// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"fmt"
	"github.com/hslam/rpc/examples/codec/json/service"
	"github.com/hslam/socket"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
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
}

func TestCheckBuffer(t *testing.T) {
	n := 1024
	buf := make([]byte, n)
	if len(checkBuffer(buf, n)) != n {
		t.Error()
	}
	buf = make([]byte, n-1)
	if len(checkBuffer(buf, n)) != n {
		t.Error()
	}
}

type mockClientCodec struct {
}

func (c *mockClientCodec) Messages() socket.Messages {
	return nil
}

func (c *mockClientCodec) WriteRequest(ctx *Context, param interface{}) error {
	return io.EOF
}

func (c *mockClientCodec) ReadResponseHeader(ctx *Context) error {
	return io.EOF
}

func (c *mockClientCodec) ReadResponseBody(reply []byte, x interface{}) error {
	return io.EOF
}

func (c *mockClientCodec) Close() error {
	return io.EOF
}

type mockServerCodec struct {
}

func (c *mockServerCodec) Messages() socket.Messages {
	return nil
}

func (c *mockServerCodec) Concurrency() int {
	return 1
}

func (c *mockServerCodec) ReadRequestHeader(ctx *Context) error {
	return io.EOF
}

func (c *mockServerCodec) ReadRequestBody(args []byte, x interface{}) error {
	return io.EOF
}

func (c *mockServerCodec) WriteResponse(ctx *Context, x interface{}) error {
	return io.EOF
}

func (c *mockServerCodec) Close() error {
	return io.EOF
}

func TestServerCodecAndClientCodec(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	codec := "json"
	opts := DefaultOptions()
	opts.Network = network
	opts.Codec = codec
	opts.HeaderEncoder = "json"
	server := NewServer()
	err := server.Register(new(service.Arith))
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.ListenWithOptions(addr, opts)
	}()
	time.Sleep(time.Millisecond * 10)
	conn, err := DialWithOptions(addr, opts)
	if err != nil {
		t.Error(err)
	}
	A := int32(4)
	B := int32(8)
	req := &service.ArithRequest{A: A, B: B}
	var res service.ArithResponse
	if err := conn.Call("Arith.Multiply", req, &res); err != nil {
		t.Error(err)
	}
	if res.Pro != A*B {
		t.Error(res.Pro)
	}
	conn.Close()
	conn.read(nil, true)
	conn.read(nil, false)
	server.Close()
	var ctx *Context
	ctx = &Context{Error: "error", codec: &mockServerCodec{}, upgrade: getUpgrade(), buffer: server.bufferPool.GetBuffer(0)}
	if err := server.ServeRequest(ctx, nil, nil, nil, nil); err == nil {
		t.Error()
	}
	ctx = &Context{Error: "error", codec: &mockServerCodec{}, upgrade: getUpgrade(), buffer: server.bufferPool.GetBuffer(0)}
	server.sendResponse(ctx)
	wg.Wait()
}

func TestContextKey(t *testing.T) {
	if BufferContextKey.String() != fmt.Sprint(BufferContextKey) {
		t.Error()
	}
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

func TestCodecCODEMore(t *testing.T) {
	var req = request{Seq: 1024, Upgrade: make([]byte, 64), ServiceMethod: strings.Repeat("a", 64), Args: make([]byte, 64)}
	var res = response{Seq: 1024, Error: strings.Repeat("e", 64), Reply: make([]byte, 64)}
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

func TestCodecPBMore(t *testing.T) {
	var c = GOGOPBCodec{}
	{
		var obj pbRequest
		n := byte(reflect.TypeOf(&obj).Elem().NumField())
		for i := byte(1); i <= n; i++ {
			if err := c.Unmarshal([]byte{i<<3 | 7}, &obj); err == nil {
				t.Error()
			}
		}
	}
	{
		var obj pbResponse
		n := byte(reflect.TypeOf(&obj).Elem().NumField())
		for i := byte(1); i <= n; i++ {
			if err := c.Unmarshal([]byte{i<<3 | 7}, &obj); err == nil {
				t.Error()
			}
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
