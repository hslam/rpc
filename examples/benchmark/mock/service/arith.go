package service
import (
	"errors"
	"time"
)

type Arith struct {}

func (this *Arith) Multiply(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

func (this *Arith) Divide(req *ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}

func (this *Arith) Multiply0ms(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

func (this *Arith) Multiply1ms(req *ArithRequest, res *ArithResponse) error {
	time.Sleep(time.Millisecond*1)
	res.Pro = req.A * req.B
	return nil
}

func (this *Arith) Multiply10ms(req *ArithRequest, res *ArithResponse) error {
	time.Sleep(time.Millisecond*10)
	res.Pro = req.A * req.B
	return nil
}
func (this *Arith) Multiply50ms(req *ArithRequest, res *ArithResponse) error {
	time.Sleep(time.Millisecond*50)
	res.Pro = req.A * req.B
	return nil
}