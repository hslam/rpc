package rpc

import "fmt"

type BatchHandle struct {
	server *Server
	msg *Msg
}

func newBatchHandle(server *Server,msg *Msg) *BatchHandle {
	return &BatchHandle{msg:msg,server:server,}
}

func (b *BatchHandle)Do(requestBody []byte)([]byte,error) {
	fmt.Println(len(requestBody))
	b.msg.data=requestBody
	res_bytes,nr:=b.server.Handler(b.msg)
	if nr==true{
		return []byte(""),nil
	}else {
		return res_bytes,nil
	}
	return nil,nil
}

