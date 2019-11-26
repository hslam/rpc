package rpc

import "time"

type Call struct {
	ServiceMethod string
	Args		interface{}
	Reply		interface{}
	Error		error
	Done		chan *Call
	start		time.Time
}

func (call *Call) done() {
	select {
	case call.Done <- call:
	default:
	}
}