package pb

import "errors"

type Seq struct {
	seq int32
}

func (this *Seq) Reset(res *SeqResponse) error {
	this.seq = -1
	res.B = 1
	return nil
}
func (this *Seq) Check(req *SeqRequest, res *SeqResponse) error {
	if req.A == this.seq+1 {
		this.seq = req.A
		res.B = req.A
		return nil
	} else {
		return errors.New("pipelining error")
	}
}
