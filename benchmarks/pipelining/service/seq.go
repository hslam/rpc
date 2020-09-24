package service

import (
	"errors"
	"sync"
)

//Seq defines a seq.
type Seq struct {
	mut sync.Mutex
	seq int32
}

//Reset resets the seq.
func (s *Seq) Reset(req *Message, res *Message) error {
	s.mut.Lock()
	s.seq = -1
	res.Value = 1
	s.mut.Unlock()
	return nil
}

//Check checks the seq.
func (s *Seq) Check(req *Message, res *Message) error {
	s.mut.Lock()
	if req.Value == s.seq+1 {
		s.seq = req.Value
		res.Value = req.Value
		s.mut.Unlock()
		return nil
	}
	s.mut.Unlock()
	return errors.New("pipelining error")
}
