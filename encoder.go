// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/codec"
	"github.com/hslam/rpc/encoder"
	"github.com/hslam/rpc/encoder/code"
	"github.com/hslam/rpc/encoder/codepb"
	"github.com/hslam/rpc/encoder/json"
	"github.com/hslam/rpc/encoder/pb"
	"sync"
)

var encoders = sync.Map{}

func init() {
	RegisterEncoder("json", json.NewEncoder)
	RegisterEncoder("code", code.NewEncoder)
	RegisterEncoder("codepb", codepb.NewEncoder)
	RegisterEncoder("pb", pb.NewEncoder)
}

func RegisterEncoder(name string, New func() *encoder.Encoder) {
	encoders.Store(name, New)
}

func NewEncoder(name string) func() *encoder.Encoder {
	if c, ok := encoders.Load(name); ok {
		return c.(func() *encoder.Encoder)
	}
	return nil
}

func DefaultEncoder() *encoder.Encoder {
	return encoder.NewEncoder(code.NewRequest(), code.NewResponse(), &codec.CODECodec{})
}
