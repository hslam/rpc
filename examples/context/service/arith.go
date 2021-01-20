package service

import (
	"context"
	"github.com/hslam/rpc"
)

//Arith defines a arith struct.
type Arith struct{}

//Multiply operation
func (a *Arith) Multiply(ctx context.Context, req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	rpc.FreeContextBuffer(ctx)
	return nil
}
