// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder/code"
	"github.com/hslam/rpc/encoder/codepb"
	"github.com/hslam/rpc/encoder/json"
	"github.com/hslam/rpc/encoder/pb"
	"testing"
)

func BenchmarkRequestMarshalCODE(t *testing.B) {
	var req = code.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &req)
	}
}

func BenchmarkRequestMarshalCODEPB(t *testing.B) {
	var req = codepb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &req)
	}
}

func BenchmarkRequestMarshalGOGOPB(t *testing.B) {
	var req = pb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &req)
	}
}

func BenchmarkRequestMarshalJSON(t *testing.B) {
	var req = json.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &req)
	}
}

func BenchmarkRequestUnmarshalCODE(t *testing.B) {
	var req = code.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &req)
	var reqCopy code.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkRequestUnmarshalCODEPB(t *testing.B) {
	var req = codepb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &req)
	var reqCopy codepb.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkRequestUnmarshalGOGOPB(t *testing.B) {
	var req = pb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &req)
	var reqCopy pb.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkRequestUnmarshalJSON(t *testing.B) {
	var req = json.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &req)
	var reqCopy json.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkRequestRoundtripCODE(t *testing.B) {
	var req = code.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	var reqCopy code.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &req)
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkRequestRoundtripCODEPB(t *testing.B) {
	var req = codepb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	var reqCopy codepb.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &req)
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkRequestRoundtripGOGOPB(t *testing.B) {
	var req = pb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var buf = make([]byte, 1024)
	var reqCopy pb.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &req)
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkRequestRoundtripJSON(t *testing.B) {
	var req = json.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var buf = make([]byte, 1024)
	var reqCopy json.Request
	d, _ := c.Marshal(buf, &req)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &req)
		c.Unmarshal(data, &reqCopy)
	}
}

func BenchmarkResponseMarshalCODE(t *testing.B) {
	var res = code.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &res)
	}
}

func BenchmarkResponseMarshalCODEPB(t *testing.B) {
	var res = codepb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &res)
	}
}

func BenchmarkResponseMarshalGOGOPB(t *testing.B) {
	var res = pb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &res)
	}
}

func BenchmarkResponseMarshalJSON(t *testing.B) {
	var res = json.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var buf = make([]byte, 1024)
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Marshal(buf, &res)
	}
}

func BenchmarkResponseUnmarshalCODE(t *testing.B) {
	var res = code.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &res)
	var resCopy code.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkResponseUnmarshalCODEPB(t *testing.B) {
	var res = codepb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &res)
	var resCopy codepb.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkResponseUnmarshalGOGOPB(t *testing.B) {
	var res = pb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &res)
	var resCopy pb.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkResponseUnmarshalJSON(t *testing.B) {
	var res = json.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var buf = make([]byte, 1024)
	data, _ := c.Marshal(buf, &res)
	var resCopy json.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkResponseRoundtripCODE(t *testing.B) {
	var res = code.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	var resCopy code.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkResponseRoundtripCODEPB(t *testing.B) {
	var res = codepb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var buf = make([]byte, 1024)
	var resCopy codepb.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkResponseRoundtripGOGOPB(t *testing.B) {
	var res = pb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var buf = make([]byte, 1024)
	var resCopy pb.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkResponseRoundtripJSON(t *testing.B) {
	var res = json.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var buf = make([]byte, 1024)
	var resCopy json.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ := c.Marshal(buf, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkRoundtripCODE(t *testing.B) {
	var req = code.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = code.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var bufq = make([]byte, 1024)
	var bufs = make([]byte, 1024)
	var reqCopy code.Request
	var resCopy code.Response
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

func BenchmarkRoundtripCODEPB(t *testing.B) {
	var req = codepb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = codepb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.CODECodec{}
	var bufq = make([]byte, 1024)
	var bufs = make([]byte, 1024)
	var reqCopy codepb.Request
	var resCopy codepb.Response
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

func BenchmarkRoundtripGOGOPB(t *testing.B) {
	var req = pb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = pb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var bufq = make([]byte, 1024)
	var bufs = make([]byte, 1024)
	var reqCopy pb.Request
	var resCopy pb.Response
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

func BenchmarkRoundtripJSON(t *testing.B) {
	var req = json.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = json.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var bufq = make([]byte, 1024)
	var bufs = make([]byte, 1024)
	var reqCopy json.Request
	var resCopy json.Response
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
