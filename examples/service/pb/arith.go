package pb

import (
	"errors"
	"time"
)

//Arith defines the struct of arith.
type Arith struct{}

//Multiply operation
func (a *Arith) Multiply(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

//Divide operation
func (a *Arith) Divide(req *ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}

//Multiply0ms sleeps 0 ms
func (a *Arith) Multiply0ms(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

//Multiply1ms sleeps 1 ms
func (a *Arith) Multiply1ms(req *ArithRequest, res *ArithResponse) error {
	time.Sleep(time.Millisecond * 1)
	res.Pro = req.A * req.B
	return nil
}

//Multiply10ms sleeps 10 ms
func (a *Arith) Multiply10ms(req *ArithRequest, res *ArithResponse) error {
	time.Sleep(time.Millisecond * 10)
	res.Pro = req.A * req.B
	return nil
}

//Multiply50ms sleeps 50 ms
func (a *Arith) Multiply50ms(req *ArithRequest, res *ArithResponse) error {
	time.Sleep(time.Millisecond * 50)
	res.Pro = req.A * req.B
	return nil
}
