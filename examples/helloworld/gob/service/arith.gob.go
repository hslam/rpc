package service

type ArithRequest struct {
	A                    int32
	B                    int32
}

type ArithResponse struct {
	Pro                  int32
	Quo                  int32
	Rem                  int32
}