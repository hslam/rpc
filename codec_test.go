package rpc

//import (
//	"testing"
//	// "hslam.com/git/x/rpc/pb"
//	//"github.com/golang/protobuf/proto"
//)

//func BenchmarkMe(t *testing.B) {
//	msg = Msg{}
//	name:="Arith.Multiply"
//	req_bytes:=make([]byte,20)
//	t.ResetTimer()
//	for i := 0; i < t.N; i++ {
//		//Serialize
//		Bytes:=msg.Serialize(Version,name,req_bytes)
//		//Deserialize
//		msg.Deserialize(Bytes)
//	}
//}

//func BenchmarkProto(t *testing.B) {
//	name:="Arith.Multiply"
//	req_bytes:=make([]byte,20)
//	t.ResetTimer()
//	for i := 0; i < t.N; i++ {
//		//Serialize
//		rpc_req:=pb.MrpcReq{Version:Version,Method:name,Data:req_bytes}
//		Bytes,_:=proto.Marshal(&rpc_req)
//		//Deserialize
//		var rpc_req_decode pb.MrpcReq
//		proto.Unmarshal(Bytes, &rpc_req_decode)
//	}
//}
//go test -v -bench=. -benchmem
//go test -v -run="none" -bench=. -benchtime=3s
//goos: darwin
//goarch: amd64
//pkg: hslam.com/git/x/mrpc
//BenchmarkMe-4      	20000000	       175 ns/op
//BenchmarkProto-4   	10000000	       405 ns/op
//PASS
//ok  	hslam.com/git/x/mrpc	8.205s