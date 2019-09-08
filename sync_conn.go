package rpc

type SyncConn struct {
	server *Server
}

func newSyncConn(server *Server) *SyncConn {
	return &SyncConn{server:server}
}

func (s *SyncConn)Do(requestBody []byte)([]byte,error) {
	_,res_bytes,err:= s.server.ServeRPC(requestBody)
	if res_bytes!=nil{
		return res_bytes,err
	}
	return nil,nil
}