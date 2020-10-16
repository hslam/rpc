// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

type jsonRequest struct {
	Seq           uint64 `json:"i"`
	Upgrade       []byte `json:"u"`
	ServiceMethod string `json:"m"`
	Args          []byte `json:"p"`
}

type jsonResponse struct {
	Seq       uint64 `json:"i"`
	CallError bool   `json:"c"`
	Error     string `json:"e"`
	Reply     []byte `json:"r"`
}
