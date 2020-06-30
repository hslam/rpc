package codec

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/codec/code"
	"github.com/hslam/rpc/codec/json"
	"github.com/hslam/rpc/codec/pb"
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

func BenchmarkRequestUnmarshalPB(t *testing.B) {
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

func BenchmarkRequestRoundtripPB(t *testing.B) {
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

func BenchmarkResponseUnmarshalPB(t *testing.B) {
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

func BenchmarkResponseRoundtripPB(t *testing.B) {
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
	var buf = make([]byte, 1024)
	var reqCopy code.Request
	var resCopy code.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	var data []byte
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ = c.Marshal(buf, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(buf, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkRoundtripPB(t *testing.B) {
	var req = pb.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = pb.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.GOGOPBCodec{}
	var buf = make([]byte, 1024)
	var reqCopy pb.Request
	var resCopy pb.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	var data []byte
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ = c.Marshal(buf, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(buf, &res)
		c.Unmarshal(data, &resCopy)
	}
}

func BenchmarkRoundtripJSON(t *testing.B) {
	var req = json.Request{Seq: 1024, ServiceMethod: "Arith.Multiply", Args: make([]byte, 512)}
	var res = json.Response{Seq: 1024, Error: "", Reply: make([]byte, 512)}
	var c = codec.JSONCodec{}
	var buf = make([]byte, 1024)
	var reqCopy json.Request
	var resCopy json.Response
	d, _ := c.Marshal(buf, &res)
	t.SetBytes(int64(len(d)))
	var data []byte
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		data, _ = c.Marshal(buf, &req)
		c.Unmarshal(data, &reqCopy)
		data, _ = c.Marshal(buf, &res)
		c.Unmarshal(data, &resCopy)
	}
}
