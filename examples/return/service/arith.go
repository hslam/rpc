package service

import (
	"context"
)

//Arith defines a arith struct.
type Arith struct{}

//Multiply operation
func (a *Arith) Multiply(ctx context.Context, req *ArithRequest) (*ArithResponse, error) {
	var res ArithResponse
	res.Pro = req.A * req.B
	return &res, nil
}
