package xml

type ArithRequest struct {
	A                    int32    `xml:"a,omitempty"`
	B                    int32    `xml:"b,omitempty"`
}

type ArithResponse struct {
	Pro                  int32    `xml:"pro,omitempty"`
	Quo                  int32    `xml:"quo,omitempty"`
	Rem                  int32    `xml:"rem,omitempty"`
}