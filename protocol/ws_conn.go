package protocol

import (
	"github.com/gorilla/websocket"
)

type WSConn struct{ *websocket.Conn}
func (c WSConn) Write(b []byte) (int, error) {
	err := c.WriteMessage(websocket.BinaryMessage,b)
	return len(b),err
}
func (c WSConn) Read(p []byte) (n int, err error){
	_, data, err := c.ReadMessage()
	copy(p,data)
	return len(data),err
}
func (c WSConn) Close()(error){
	return c.Conn.Close()
}