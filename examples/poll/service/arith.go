package service

//Arith defines a arith struct.
type Arith struct{}

//Multiply operation
func (a *Arith) Multiply(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}
