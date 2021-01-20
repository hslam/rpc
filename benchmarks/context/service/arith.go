package service

import (
	"context"
)

//Arith defines a arith struct.
type Arith struct{}

//Multiply operation
func (a *Arith) Multiply(ctx context.Context, req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}
