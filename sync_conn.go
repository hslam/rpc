package rpc

type SyncConn struct {
	server *Server
}

func newSyncConn(server *Server) *SyncConn {
	return &SyncConn{server:server}
}

func (s *SyncConn)Do(requestBody []byte)([]byte,error) {
	_,res_bytes:= s.server.Serve(requestBody)
	if res_bytes!=nil{
		return res_bytes,nil
	}
	return nil,nil
}
