package service

type ArithRequest struct {
	A                    int32    `json:"a,omitempty"`
	B                    int32    `json:"b,omitempty"`
}

type ArithResponse struct {
	Pro                  int32    `json:"pro,omitempty"`
	Quo                  int32    `json:"quo,omitempty"`
	Rem                  int32    `json:"rem,omitempty"`
}