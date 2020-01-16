package rpc

import (
	"github.com/gorilla/websocket"
)

type websocketConn struct {
	*websocket.Conn
}

func (s *websocketConn) ReadMessage() (p []byte, err error) {
	_, data, err := s.Conn.ReadMessage()
	return data, err
}

func (s *websocketConn) WriteMessage(b []byte) (err error) {
	return s.Conn.WriteMessage(websocket.BinaryMessage, b)
}

func (s *websocketConn) Close() error {
	return s.Conn.Close()
}
