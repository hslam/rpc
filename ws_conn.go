package rpc

import (
	"github.com/gorilla/websocket"
)

type wsConn struct {
	*websocket.Conn
}

func (s *wsConn) ReadMessage() (p []byte, err error) {
	_, data, err := s.Conn.ReadMessage()
	return data, err
}

func (s *wsConn) WriteMessage(b []byte) (err error) {
	return s.Conn.WriteMessage(websocket.BinaryMessage, b)
}

func (s *wsConn) Close() error {
	return s.Conn.Close()
}
