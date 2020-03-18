package pb

import "errors"

//Seq defines a seq.
type Seq struct {
	seq int32
}

//Reset resets the seq.
func (s *Seq) Reset(res *SeqResponse) error {
	s.seq = -1
	res.B = 1
	return nil
}

//Check checks the seq.
func (s *Seq) Check(req *SeqRequest, res *SeqResponse) error {
	if req.A == s.seq+1 {
		s.seq = req.A
		res.B = req.A
		return nil
	}
	return errors.New("pipelining error")
}
