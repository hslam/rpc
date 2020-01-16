package rpc

type syncConn struct {
	server *Server
}

func newSyncConn(server *Server) *syncConn {
	return &syncConn{server: server}
}

func (s *syncConn) Do(requestBody []byte) ([]byte, error) {
	_, resBytes := s.server.Serve(requestBody)
	if resBytes != nil {
		return resBytes, nil
	}
	return nil, nil
}
