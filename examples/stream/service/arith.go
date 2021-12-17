package service

import (
	"github.com/hslam/rpc"
)

//Arith defines the struct of arith.
type Arith struct{}

//StreamMultiply operation
func (a *Arith) StreamMultiply(stream *Stream) error {
	for {
		var req = &ArithRequest{}
		if err := stream.Read(req); err != nil {
			return err
		}
		res := ArithResponse{}
		res.Pro = req.A * req.B
		if err := stream.Write(&res); err != nil {
			println(err.Error())
			return err
		}
	}
}

// Stream is used to connect rpc Stream.
type Stream struct {
	stream rpc.Stream
}

// Connect connects rpc Stream.
func (s *Stream) Connect(stream rpc.Stream) error {
	s.stream = stream
	return nil
}

//Read reads a message from the rpc stream.
func (s *Stream) Read(req *ArithRequest) error {
	return s.stream.ReadMessage(req)
}

//Write writes a message to the rpc stream.
func (s *Stream) Write(res *ArithResponse) error {
	return s.stream.WriteMessage(res)
}
