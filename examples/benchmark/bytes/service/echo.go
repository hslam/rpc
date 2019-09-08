package service

import "strings"

type Echo struct {
	value []byte
}

func (this *Echo) ToLower(req *[]byte, res *[]byte) error {
	*res=[]byte(strings.ToLower(string(*req)))
	return nil
}

func (this *Echo) Set(req *[]byte) error {
	this.value=*req
	return nil
}

func (this *Echo) Get(res *[]byte) error {
	*res=this.value
	return nil
}
func (this *Echo) Clear() error {
	this.value=nil
	return nil
}