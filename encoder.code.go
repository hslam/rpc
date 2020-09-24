// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/codec"
)

// NewCODEEncoder returns a header Encoder.
func NewCODEEncoder() *Encoder {
	return NewEncoder(NewCODERequest(), NewCODEResponse(), NewCODECodec())
}

// NewCODECodec returns the instance of Codec.
func NewCODECodec() codec.Codec {
	return &codec.CODECodec{}
}

// NewCODERequest returns the instance of codeRequest.
func NewCODERequest() Request {
	return &request{}
}

// NewCODEResponse returns the instance of codeResponse.
func NewCODEResponse() Response {
	return &response{}
}
