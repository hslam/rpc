package service

import "strings"

type Echo struct {}

func (this *Echo) ToLower(req *[]byte, res *[]byte) error {
	*res=[]byte(strings.ToLower(string(*req)))
	return nil
}
